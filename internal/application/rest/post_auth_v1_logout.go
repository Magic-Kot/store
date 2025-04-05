package rest

import (
	"fmt"
	"net/http"

	"git.appkode.ru/pub/go/failure"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

func (s Server) PostAuthV1Logout(w http.ResponseWriter, r *http.Request) error {
	var request entity.PostAuthLogoutRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}

	if request.RefreshToken == "" {
		return failure.NewInvalidArgumentError("refreshToken is empty")
	}

	if err := s.auth.Logout(r.Context(), value.RefreshToken(request.RefreshToken)); err != nil {
		return fmt.Errorf("auth.Logout: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
