package rest

import (
	"fmt"
	"github.com/Magic-Kot/store/internal/domain/entity"
	"net/http"

	"github.com/Magic-Kot/store/pkg/contextx"

	"git.appkode.ru/pub/go/failure"
)

func (s Server) PatchSettingsV1User(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	personID, err := contextx.PersonIDFromContext(ctx)
	if err != nil {
		return failure.NewInvalidArgumentErrorFromError(fmt.Errorf("contextx.PersonIDFromContext: %w", err))
	}

	var request entity.UserData

	if err = readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}

	if request.Age == 0 && request.Username == "" && request.Name == "" && request.Surname == "" {
		return failure.NewInvalidArgumentError("request is empty")
	}

	if err = s.user.UpdateUser(ctx, personID, request); err != nil {
		return fmt.Errorf("user.UpdateUser: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
