package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/internal/services/user"
	"github.com/Magic-Kot/store/pkg/utils/jwtToken"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiController struct {
	UserService user.UserService
	logger      *zerolog.Logger
	validator   *validator.Validate
}

func NewApiController(userService *user.UserService, logger *zerolog.Logger, validator *validator.Validate) *ApiController {
	return &ApiController{
		UserService: *userService,
		logger:      logger,
		validator:   validator,
	}
}

// GetUser - получение сущности пользователя по ID
func (ac *ApiController) GetUser(c echo.Context) error {
	req := new(models.User)
	id := c.Get("userId")
	userID, ok := id.(int)
	if ok != true {
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userID

	result, err := ac.UserService.GetUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateUser - регистрация нового пользователя
func (ac *ApiController) CreateUser(c echo.Context) error {
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

// SignIn - индетификация, аутентификация пользователя, парсинг Username, Password
func (ac *ApiController) SignIn(c echo.Context) error {
	req := new(models.UserAuthorization)
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

	token, err := ac.UserService.SignIn(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successful authorization, jwtToken: %s", token))
}

// AuthorizationUser - авторизация пользователя, парсинг токена
func (ac *ApiController) AuthorizationUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			ac.logger.Debug().Msgf("header 'Authorization': invalid request: %v", header)

			return c.JSON(http.StatusUnauthorized, fmt.Sprint("empty auth header"))
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ac.logger.Debug().Msgf("header 'Authorization': invalid request: %v", header)

			return c.JSON(http.StatusUnauthorized, fmt.Sprint("invalid auth header"))
		}

		userId, err := jwtToken.ParseToken(headerParts[1])
		if err != nil {
			ac.logger.Debug().Msgf("parseToken: %v", err)

			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("userId", userId)
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
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	id := c.Get("userId")
	userID, ok := id.(int)
	if ok != true {
		ac.logger.Debug().Msgf("UpdateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userID

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" {
				switch err.Tag() {
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 4 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Email" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("Incorrect email"))
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err = ac.UserService.UpdateUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprint("successfully updated"))
}

// DeleteUser - удаление пользователя по ID
func (ac *ApiController) DeleteUser(c echo.Context) error {
	id := c.Get("userId")
	userID, ok := id.(int)
	if ok != true {
		ac.logger.Debug().Msgf("DeleteUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	err := ac.UserService.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted user: %d", userID))
}
