package httpecho

import (
	"github.com/Magic-Kot/store/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetReferralRoutes(e *echo.Echo, middleware *controllers.ApiController, apiController *controllers.ApiReferralController) {
	referral := e.Group("/bonuses", middleware.AuthorizationUser)
	{
		referral.POST("/friends", apiController.CreateReferral) // создание реферальной ссылки
		//referral.GET("/counter", apiController.CounterReferral)
	}
	e.GET("/baf/:url", apiController.GetReferral) // переход по реферальной ссылке
}
