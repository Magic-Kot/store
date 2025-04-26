//nolint:dupl // not a duplicate
package component

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type RefreshTokenParser struct {
	parser JWTParserES256
}

func NewRefreshTokenParser(
	rawPublicKey string,
) RefreshTokenParser {
	parser := NewJWTParserES256(rawPublicKey)

	return RefreshTokenParser{parser: parser}
}

func (j RefreshTokenParser) Parse(token value.RefreshToken) (value.RefreshTokenClaims, error) {
	var claims value.RefreshTokenClaims

	err := j.parser.Parse(token, &claims)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, fmt.Errorf("parser.Parse: %w", err) // Refresh token expired
		}

		return claims, fmt.Errorf("parser.Parse: %w", err) // Refresh token invalid
	}

	return claims, nil
}
