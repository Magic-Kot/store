package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/labstack/echo/v4"
)

// AuthorizationUser - авторизация пользователя, парсинг токена
func (ac *ApiController) AuthorizationUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			ac.logger.Debug().Msgf("header 'Authorization': invalid request: %v", header)

			return c.JSON(http.StatusUnauthorized, fmt.Sprint("empty auth header"))
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ac.logger.Debug().Msgf("header 'Authorization': invalid request: %v", header)

			return c.JSON(http.StatusUnauthorized, fmt.Sprint("invalid auth header"))
		}

		userId, err := jwt_token.ParseToken(headerParts[1])
		if err != nil {
			ac.logger.Debug().Msgf("parseToken: %v", err)

			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("userId", userId)
		err = next(c)
		if err != nil {
			ac.logger.Warn().Msgf("next HandlerFunc: %v", err)

			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return nil
	}
}
