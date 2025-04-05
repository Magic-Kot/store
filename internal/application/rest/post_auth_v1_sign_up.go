package rest

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"git.appkode.ru/pub/go/failure"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

func (s Server) PostAuthV1SignUp(w http.ResponseWriter, r *http.Request) error {
	var request entity.PostAuthSignUpRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}

	if request.Login == "" || request.Password == "" {
		return failure.NewInvalidArgumentError("login or password is empty")
	}

	if utf8.RuneCountInString(request.Login) < 4 || utf8.RuneCountInString(request.Login) > 20 {
		return failure.NewInvalidArgumentError("login is invalid")
	}

	if utf8.RuneCountInString(request.Password) < 6 || utf8.RuneCountInString(request.Password) > 20 {
		return failure.NewInvalidArgumentError("password is invalid")
	}

	if err := s.auth.Registration(r.Context(), value.Login(request.Login), value.Password(request.Password)); err != nil {
		return fmt.Errorf("auth.Registration: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
