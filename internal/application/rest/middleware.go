package rest

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/contextx"
	"github.com/rs/zerolog"
)

type tokenFinder interface {
	FindRawToken(*http.Request) (string, error)
}

type accessTokenParser interface {
	Parse(token value.AccessToken) (value.AccessTokenClaims, error)
}

type BearerAuth struct {
	tokenFinder       tokenFinder
	accessTokenParser accessTokenParser
}

func NewBearerAuth(
	tokenFinder tokenFinder,
	accessTokenParser accessTokenParser,
) BearerAuth {
	return BearerAuth{
		tokenFinder:       tokenFinder,
		accessTokenParser: accessTokenParser,
	}
}

func (j BearerAuth) JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := zerolog.Ctx(ctx)

			rawAccessToken, err := j.tokenFinder.FindRawToken(r)
			if err != nil {
				logger.Error().Msg(fmt.Sprintf("failed to retrieve access token, error: %s", err))

				return
			}

			accessTokenClaims, err := j.accessTokenParser.Parse(value.AccessToken(rawAccessToken))
			if err != nil {
				logger.Error().Msg(fmt.Sprintf("failed to parse access token, error: %s", err))

				ctx = contextx.WithPersonID(ctx, accessTokenClaims.PersonID)

				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
