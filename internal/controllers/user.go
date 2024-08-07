package controllers

import (
	"net/http"

	"online-store/internal/services/user"

	"github.com/labstack/echo/v4"
)

type ApiController struct {
	UserService user.UserService
}

type request struct {
	name  string //`json:"user_name"`
	login string //`json:"login"`
}

func NewApiController(userService *user.UserService) *ApiController {
	return &ApiController{
		UserService: *userService,
	}
}

// GetUser - получение login по user_name
func (ac *ApiController) GetUser(c echo.Context) error {
	req := request{
		name: c.QueryParam("user_name"),
		// ... необходимо распарсить запрос в структуру
	}

	name, err := ac.UserService.GetUser(c.Request().Context(), req.name)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, "Hello,"+name)
}

// CreateUser - создание нового пользователя
func (ac *ApiController) CreateUser(c echo.Context) error {
	req := request{
		name:  c.QueryParam("user_name"),
		login: c.QueryParam("login"),
	}

	rez, err := ac.UserService.CreateUser(c.Request().Context(), req.name, req.login)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, rez)
}

func (ac *ApiController) UpdateUser(c echo.Context) error {
	req := request{
		name:  c.QueryParam("user_name"),
		login: c.QueryParam("login"),
	}

	rez, err := ac.UserService.UpdateUser(c.Request().Context(), req.name, req.login)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, rez)
}

func (ac *ApiController) DeleteUser(c echo.Context) error {
	req := request{
		name: c.QueryParam("user_name"),
	}

	rez, err := ac.UserService.DeleteUser(c.Request().Context(), req.name)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.String(http.StatusOK, rez)
}
