package httpecho

import (
	"online-store/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	r := e.Group("/user")

	r.GET("/", apiController.GetUser)
	r.POST("/create", apiController.CreateUser)
	r.PUT("/update/:id", apiController.UpdateUser)
	r.DELETE("/delete/:id", apiController.DeleteUser)
}
