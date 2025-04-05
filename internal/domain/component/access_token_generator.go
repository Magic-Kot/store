package component

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type AccessTokenGenerator struct {
	generator JWTGeneratorES256[value.AccessTokenClaims]
	ttl       time.Duration
}

func NewAccessTokenGenerator(
	rawPrivateKey string,
	ttl time.Duration,
	clock clock.Clock, //nolint:gocritic
) AccessTokenGenerator {
	generator := NewJWTGeneratorES256[value.AccessTokenClaims](rawPrivateKey, clock)

	return AccessTokenGenerator{
		generator: generator,
		ttl:       ttl,
	}
}

func (j AccessTokenGenerator) Generate(
	personID value.PersonID,
) (value.AccessToken, error) {
	timeNow := timeNowUTC(j.generator.clock)

	claims := value.AccessTokenClaims{ //nolint:exhaustruct
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
		},
		PersonID: personID,
	}

	token, err := j.generator.Generate(claims)
	if err != nil {
		return "", fmt.Errorf("generator.Generate: %w", err)
	}

	return value.AccessToken(token), nil
}
