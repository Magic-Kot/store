package component

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/benbjohnson/clock"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

type JWTGeneratorES256[ClaimsT jwt.Claims] struct {
	clock         clock.Clock
	signingMethod jwt.SigningMethod
	privateKey    *ecdsa.PrivateKey
}

func NewJWTGeneratorES256[ClaimsT jwt.Claims](
	rawPrivateKey string,
	clock clock.Clock, //nolint:gocritic
) JWTGeneratorES256[ClaimsT] {
	privateKey := lo.Must(jwt.ParseECPrivateKeyFromPEM([]byte(rawPrivateKey)))

	return JWTGeneratorES256[ClaimsT]{
		privateKey:    privateKey,
		clock:         clock,
		signingMethod: jwt.SigningMethodES256,
	}
}

func (j JWTGeneratorES256[ClaimsT]) Generate(
	claims ClaimsT,
) (string, error) {
	token := jwt.NewWithClaims(j.signingMethod, claims)

	signedString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return signedString, nil
}
