# Online store project

## Инфраструктура

Для запуска сервису необходимы:

- Подключение к базе Postgres для хранения данных приложения
- Подключение к Redis для хранения активных JWT Refresh-токенов.

Для миграции баз Postgres используется библиотека "github.com/pressly/goose"

Для для подключения к Postgres и выполнения SQL запросов используются
библиотеки:

- github.com/jackc/pgx/v5
- github.com/jmoiron/sqlx

Для для подключения к Redis и выполнения запросов используется библиотека
"github.com/redis/go-redis/v9"

## Web-сервер

Web-сервер реализован с использованием стандартной библиотеки "net/http", в
качестве роутера используется "github.com/go-chi/chi"

Каждый HTTP-запрос и ответ логируется в формате JSON в STDOUT. Для логирования
используется библиотека "log/slog". Список полей, которые присутствуют в
каждой строке лога:

- app-name
- app-version
- trace-id
- url
- personID

Список дополнительных полей, которые присутствуют в каждой строке лога
HTTP-запроса:

- request-body

Список дополнительных полей, которые присутствуют в каждой строке лога ответа:

- response-status
- response-headers
- response-body
- duration-ms

Если длина любого поля в строке лога превышает 2000 символов, оно обрезается до
этого значения (конфигурируется через ENV).

## Пользователи и сессии

Для работы с сессиями используется JWT с алгоритмом шифрования ES256. После
успешной регистрации, аутентификации или обновления сессии сервер выдаёт
клиенту новую пару AccessToken + RefreshToken. Используемая библиотека:
"github.com/golang-jwt/jwt/v5".

- Время жизни AccessToken'а: 5 минут (конфигурируется через ENV)
- Время жизни RefreshToken'а: 7 дней (конфигурируется через ENV)

Пользовательские активные RefreshToken'ы сохраняются в базе Redis. Это нужно
для того, чтобы инвалидировать все старые RefreshToken'ы при обновлении сессии
и получении новой пары AccessToken + RefreshToken.
## Интеграционные тесты

Для запуска интеграционных тестов нужно сначала запустить тестовую среду:

```shell
make test-infrastructure
```

После этого, когда всё запустится, можно выполнить интеграционные тесты:

```shell
make test
```

## API

1) User registration. Request Example:

| Path               | Method | Request                                                   | Description  |
|--------------------|--------|-----------------------------------------------------------|--------------|
| `/auth/v1/sign-up` | POST   | Body: `{"login": "username", "password": "userpassword"}` | Registration |

2) Authorization. Request Example:

| Path               | Method | Request                                                   | Description    |
|--------------------|--------|-----------------------------------------------------------|----------------|
| `/auth/v1/sign-in` | POST   | Body: `{"login": "username", "password": "userpassword"}` | Authorization  |
| `/auth/v1/refresh` | POST   | Body: `{"refreshToken": "token"}`                         | Refresh tokens |
| `/auth/v1/logout`  | POST   | Body: `{"refreshToken": "token"}`                         | Logout         |

3) Working with the user. Request Example:

| Path                | Method | Request                                                                                                              | Description                           |
|---------------------|--------|----------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| `/user/v1/info`     | GET    | Header: `Authorization: token`                                                                                       | Getting user data                     |
| `/settings/v1/user` | PATCH  | Header: `Authorization: token`<br/>Body: `{"login": "username", "name": "name", "surname": "surname", "age": "age"}` | Changing the login or other user data |
| `/settings/v1/user` | DELETE | Header: `Authorization: token`                                                                                       | Deleting the user                     |
