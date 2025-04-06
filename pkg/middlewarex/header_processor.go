package middlewarex

import (
	"fmt"
	"net/http"
	"strings"
)

//nolint:gosec
const (
	headerNameAuthorization = "Authorization"
	bearerPrefix            = "Bearer "
)

type HeaderTokenFinder struct {
	authorizationHeader string
	authorizationPrefix string
}

func NewHeaderTokenFinder(authorizationHeader, authorizationPrefix string) HeaderTokenFinder {
	return HeaderTokenFinder{
		authorizationHeader: authorizationHeader,
		authorizationPrefix: authorizationPrefix,
	}
}

func NewHeaderAuthorizationBearerTokenFinder() HeaderTokenFinder {
	return HeaderTokenFinder{
		authorizationHeader: headerNameAuthorization,
		authorizationPrefix: bearerPrefix,
	}
}

func (h HeaderTokenFinder) FindRawToken(r *http.Request) (string, error) {
	header := r.Header.Get(h.authorizationHeader)

	if !strings.HasPrefix(header, h.authorizationPrefix) {
		return "", fmt.Errorf("bad %s header - access token invalid", h.authorizationHeader)
	}

	return strings.TrimPrefix(header, h.authorizationPrefix), nil
}
