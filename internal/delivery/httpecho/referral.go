package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"
	"github.com/Magic-Kot/store/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetReferralRoutes(e *echo.Echo, apiController *controllers.ApiReferralController, middleware *middleware.Middleware) {
	referral := e.Group("/bonuses", middleware.AuthorizationUser)
	{
		referral.POST("/friends", apiController.CreateReferral)
		//referral.GET("/counter", apiController.CounterReferral)
	}
	e.GET("/baf/:url", apiController.GetReferral)
}
