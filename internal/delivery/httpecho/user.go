package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", apiController.SignUp)
		auth.POST("/sign-in", apiController.SignIn)
		auth.POST("/refresh", apiController.RefreshToken)
	}

	user := e.Group("/user", apiController.AuthorizationUser)
	//r := e.Group("/user", middleware.AuthorizationUser)
	{
		user.GET("/get", apiController.GetUser)
		user.PUT("/update", apiController.UpdateUser)
		user.DELETE("/delete", apiController.DeleteUser)
	}
}
