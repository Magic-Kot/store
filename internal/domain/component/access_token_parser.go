//nolint:dupl // not a duplicate
package component

import (
	"errors"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type AccessTokenParser struct {
	parser JWTParserES256
}

func NewAccessTokenParser(
	rawPublicKey string,
) AccessTokenParser {
	parser := NewJWTParserES256(rawPublicKey)

	return AccessTokenParser{parser: parser}
}

func (j AccessTokenParser) Parse(token value.AccessToken) (value.AccessTokenClaims, error) {
	var claims value.AccessTokenClaims

	err := j.parser.Parse(token, &claims)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, fmt.Errorf("parser.Parse: %w", err) //Access token expired
		}

		return claims, fmt.Errorf("parser.Parse: %w", err) // Access token invalid
	}

	return claims, nil
}
