package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthorizationUser - авторизация пользователя, парсинг токена
func (ac *ApiController) AuthorizationUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			ac.logger.Debug().Msgf("empty 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ac.logger.Debug().Msgf("invalid 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		id, err := ac.token.ParseToken(headerParts[1])
		if err != nil {
			ac.logger.Debug().Msgf("invalid authorization token: %v", err)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		ac.logger.Debug().Msgf("id: %s", id)
		c.Set("id", id)

		err = next(c)
		if err != nil {
			ac.logger.Warn().Msgf("next HandlerFunc: %v", err)

			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return nil
	}
}
