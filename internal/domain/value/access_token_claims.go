package value

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims

	PersonID PersonID `json:"personId"`
}
