package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", apiController.CreateUser)
		auth.POST("/sign-in", apiController.SignIn)
		auth.POST("/refresh", apiController.RefreshToken)
	}

	r := e.Group("/user", apiController.AuthorizationUser)
	{
		r.GET("/get", apiController.GetUser)
		r.PUT("/update", apiController.UpdateUser)
		r.DELETE("/delete", apiController.DeleteUser)
	}
}
