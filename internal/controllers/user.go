package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/internal/services/user"

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("invalid id: %d", id))
	}
	req.ID = id

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
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully created user, id: %d", id))
}

// AuthorizationUser - авторизация пользователя
func (ac *ApiController) AuthorizationUser(c echo.Context) error {
	req := new(models.UserAuthorization)
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

	token, err := ac.UserService.AuthorizationUser(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successful authorization, token: %d", token))
}

// UpdateUser - обновление данных пользователя по ID
func (ac *ApiController) UpdateUser(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ac.logger.Debug().Msgf("deleteUser: invalid id: %d, err: %v", id, err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}
	req.ID = id

	err = ac.validator.Struct(req)
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

	return c.JSON(http.StatusOK, fmt.Sprint("successfully updated login user"))
}

// DeleteUser - удаление пользователя по ID
func (ac *ApiController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ac.logger.Debug().Msgf("deleteUser: invalid id: %d, err: %v", id, err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	err = ac.UserService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted user: %d", id))
}
