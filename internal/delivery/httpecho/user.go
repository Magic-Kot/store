package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"
	"github.com/Magic-Kot/store/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController, middleware *middleware.Middleware) {
	e.POST("/sign-up", apiController.SignUp)

	user := e.Group("/user", middleware.AuthorizationUser)
	{
		user.GET("/get", apiController.GetUser)
		user.PATCH("/update", apiController.UpdateUser)
		user.DELETE("/delete", apiController.DeleteUser)
	}
}
