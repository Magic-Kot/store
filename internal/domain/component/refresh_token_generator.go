package component

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type RefreshTokenGenerator struct {
	generator JWTGeneratorES256[value.RefreshTokenClaims]
	ttl       time.Duration
}

func NewRefreshTokenGenerator(
	rawPrivateKey string,
	ttl time.Duration,
	clock clock.Clock, //nolint:gocritic
) RefreshTokenGenerator {
	generator := NewJWTGeneratorES256[value.RefreshTokenClaims](rawPrivateKey, clock)

	return RefreshTokenGenerator{
		generator: generator,
		ttl:       ttl,
	}
}

func (j RefreshTokenGenerator) Generate(
	personID value.PersonID,
	id value.RefreshTokenID,
) (value.RefreshToken, error) {
	timeNow := timeNowUTC(j.generator.clock)

	claims := value.RefreshTokenClaims{
		//nolint:exhaustruct
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
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

	return value.RefreshToken(token), nil
}
