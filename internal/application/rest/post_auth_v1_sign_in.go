package rest

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

func (s Server) PostAuthV1SignIn(w http.ResponseWriter, r *http.Request) error {
	var request entity.PostAuthSignInRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}

	tokenPair, err := s.auth.Authenticate(r.Context(), value.Login(request.Login), value.Password(request.Password))
	if err != nil {
		return fmt.Errorf("auth.Authenticate: %w", err)
	}

	replyJSON(r.Context(), w, http.StatusOK, tokenPair)

	return nil
}
