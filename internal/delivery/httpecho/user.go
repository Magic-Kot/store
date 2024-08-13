package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	r := e.Group("/user")

	r.POST("/sign-up", apiController.CreateUser)
	r.POST("/sign-in", apiController.AuthorizationUser)

	r.GET("/:id", apiController.GetUser)
	r.PUT("/update/:id", apiController.UpdateUser)
	r.DELETE("/delete/:id", apiController.DeleteUser)
}
