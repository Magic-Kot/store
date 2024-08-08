package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"online-store/internal/models"
	"online-store/internal/services/user"

	"github.com/labstack/echo/v4"
)

type ApiController struct {
	UserService user.UserService
}

func NewApiController(userService *user.UserService) *ApiController {
	return &ApiController{
		UserService: *userService,
	}
}

// GetUser - получение ID пользователя по login
func (ac *ApiController) GetUser(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}

	id, err := ac.UserService.GetUser(c.Request().Context(), req.Login)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

// CreateUser - создание нового пользователя
func (ac *ApiController) CreateUser(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}

	rez, err := ac.UserService.CreateUser(c.Request().Context(), req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, rez)
}

// UpdateUser - обновление login, password пользователя по ID
func (ac *ApiController) UpdateUser(c echo.Context) error {
	req := new(models.User)
	if err := c.Bind(req); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("invalid id: %d", id))
	}
	req.ID = id

	fmt.Println(req)

	rez, err := ac.UserService.UpdateUser(c.Request().Context(), req.ID, req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, rez)
}

// DeleteUser - удаление пользователя по ID
func (ac *ApiController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("invalid id: %d", id))
	}

	req := models.User{ID: id}

	fmt.Println(req.ID)

	rez, err := ac.UserService.DeleteUser(c.Request().Context(), req.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, rez)
}
