package httpecho

import (
	"github.com/labstack/echo/v4"

	"online-store/internal/controllers"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	r := e.Group("/user")

	r.GET("/", apiController.GetUser)
	r.POST("/create/:id", apiController.CreateUser)
	r.PUT("/update", apiController.UpdateUser)
	r.DELETE("/delete", apiController.DeleteUser)
}
