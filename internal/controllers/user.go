package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/internal/services/user"
	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var errAutorizationUser = errors.New("an unauthorized user")

type ApiController struct {
	UserService user.UserService
	logger      *zerolog.Logger
	validator   *validator.Validate
	token       *jwt_token.Manager
}

func NewApiController(userService *user.UserService, logger *zerolog.Logger, validator *validator.Validate, token *jwt_token.Manager) *ApiController {
	return &ApiController{
		UserService: *userService,
		logger:      logger,
		validator:   validator,
		token:       token,
	}
}

// GetUser - получение сущности пользователя по ID
func (ac *ApiController) GetUser(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	req := new(models.User)
	id := c.Get("id")

	userID, ok := id.(string)
	if ok != true {
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userIdInt

	result, err := ac.UserService.GetUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateUser - регистрация нового пользователя
func (ac *ApiController) CreateUser(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your login"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 4 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Password" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your password"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum password length is 6 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum password length is 20 characters"))
				}
			}
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := ac.UserService.CreateUser(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		// нет обработки ошибки на уникальность логина
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully created user, id: %d", id))
}

// SignIn - аутентификация пользователя, парсинг Username, Password
func (ac *ApiController) SignIn(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	req := new(models.UserAuthorization)

	req.GUID = c.QueryParam("GUID")

	if err := c.Bind(req); err != nil {
		ac.logger.Debug().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your login"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 1 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Password" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your password"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum password length is 1 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum password length is 20 characters"))
				}
			}
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	tokens, err := ac.UserService.SignIn(ctx, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = tokens.RefreshToken
	cookie.Path = "/auth"
	cookie.Expires = time.Now().Add(4 * time.Hour) // TODO: применить переменную из созданной сессии // =expiresAt
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}

// RefreshToken - получение новых refresh и access токенов
func (ac *ApiController) RefreshToken(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	cookieRequest, err := c.Cookie("refreshToken")
	if err != nil {
		return err
	}

	// TODO: валидация токена?

	tokens, err := ac.UserService.RefreshToken(ctx, cookieRequest.Value)
	if err != nil {
		//return c.Redirect(http.StatusUnauthorized, "/sign-in")
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = tokens.RefreshToken
	cookie.Path = "/auth"
	cookie.Expires = time.Now().Add(4 * time.Hour) // TODO: применить переменную из созданной сессии // =expiresAt
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}

// TODO: некорректное расположение кода AuthorizationUser

// AuthorizationUser - авторизация пользователя, парсинг токена
func (ac *ApiController) AuthorizationUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			ac.logger.Debug().Msgf("empty 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAutorizationUser)
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ac.logger.Debug().Msgf("invalid 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAutorizationUser)
		}

		id, err := ac.token.ParseToken(headerParts[1])
		if err != nil {
			ac.logger.Debug().Msgf("invalid authorization token: %v", err)

			return c.JSON(http.StatusUnauthorized, errAutorizationUser)
		}

		ac.logger.Debug().Msgf("id: %s", id)
		c.Set("id", id)

		err = next(c)
		if err != nil {
			ac.logger.Warn().Msgf("next HandlerFunc: %v", err)

			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return nil
	}
}

// UpdateUser - обновление данных пользователя по ID
func (ac *ApiController) UpdateUser(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	//req := new(models.User)
	var req models.User
	if err := c.Bind(&req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	id := c.Get("id")
	userId, ok := id.(string)
	if ok != true {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userIdInt

	err = ac.validator.Struct(&req)

	// TODO: оптимизировать код для валидации

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" && err.Value() != "" {
				switch err.Tag() {
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 4 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Email" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("Incorrect email"))
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err = ac.UserService.UpdateUser(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprint("successfully updated"))
}

// DeleteUser - удаление пользователя по ID
func (ac *ApiController) DeleteUser(c echo.Context) error {
	// TODO: прокинуть логер в контекст echo
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	id := c.Get("id")
	userId, ok := id.(string)
	if ok != true {
		ac.logger.Debug().Msgf("DeleteUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	err = ac.UserService.DeleteUser(c.Request().Context(), userIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted user: %d", userIdInt))
}
