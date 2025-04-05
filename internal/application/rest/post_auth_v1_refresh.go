package rest

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

func (s Server) PostAuthV1Refresh(w http.ResponseWriter, r *http.Request) error {
	var request entity.PostAuthRefreshRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}

	tokenPair, err := s.auth.Refresh(r.Context(), value.RefreshToken(request.RefreshToken))
	if err != nil {
		return fmt.Errorf("auth.Refresh: %w", err)
	}

	replyJSON(r.Context(), w, http.StatusOK, tokenPair)

	return nil
}
