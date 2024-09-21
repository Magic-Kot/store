# Online store project

## API

1) User registration. Request Example:  
   `host:port`

| Path       | Method | Request                                                   | Description  |
|------------|--------|-----------------------------------------------------------|--------------|
| `/sign-up` | POST   | Body: `{"login": "username", "password": "userpassword"}` | Registration |

When registering, the login must contain from 4 to 20 characters, the password - from 6 to 20.

2) Authorization. Request Example:  
   `host:port/auth`

| Path       | Method | Request                                                                                 | Description    |
|------------|--------|-----------------------------------------------------------------------------------------|----------------|
| `/sign-in` | POST   | Query Params: `GUID=guid`<br/>Body: `{"login": "username", "password": "userpassword"}` | Authorization  |
| `/refresh` | POST   | Cookie: `refreshToken=token; Path=/auth/refresh; HttpOnly;`                             | Refresh tokens |

When logging in, the Login and password must contain from 1 to 20 characters.

3) Working with the user. Request Example:  
   `host:port/user`

| Path      | Method | Request                                                                                                                                | Description                           |
|-----------|--------|----------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| `/get`    | GET    | Header: `Authorization: token`                                                                                                         | Getting user data                     |
| `/update` | PUT    | Header: `Authorization: token`<br/>Body: `{"login": "username", "name": "name", "surname": "surname", "age": "age", "email": "email"}` | Changing the login or other user data |
| `/delete` | DELETE | Header: `Authorization: token`                                                                                                         | Deleting the user                     |

The login must contain from 4 to 20 characters.

4) Referral link. Request Example:  
   `host:port`

| Path               | Method | Request                                                   | Description                 |
|--------------------|--------|-----------------------------------------------------------|-----------------------------|
| `/bonuses/friends` | POST   | Header: `Authorization: token`<br/>Body: `{"url": "url"}` | Getting a referral link     |
| `/baf/:url`        | GET    |                                                           | Clicking on a referral link |
