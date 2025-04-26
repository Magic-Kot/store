package component

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

type JWTParserES256 struct {
	keyFunc jwt.Keyfunc
}

func NewJWTParserES256(rawPublicKey string) JWTParserES256 {
	publicKey := lo.Must(jwt.ParseECPublicKeyFromPEM([]byte(rawPublicKey)))

	return JWTParserES256{
		keyFunc: func(*jwt.Token) (any, error) { return publicKey, nil },
	}
}

//nolint:nolintlint,ireturn
func (j JWTParserES256) Parse(
	rawToken fmt.Stringer,
	claims jwt.Claims,
) error {
	_, err := jwt.ParseWithClaims(rawToken.String(), claims, j.keyFunc)
	if err != nil {
		return fmt.Errorf("jwt.ParseWithClaims: %w", err)
	}

	return nil
}
