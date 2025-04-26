package rest

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/store/pkg/contextx"

	"git.appkode.ru/pub/go/failure"
)

func (s Server) DeleteSettingsV1User(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	personID, err := contextx.PersonIDFromContext(ctx)
	if err != nil {
		return failure.NewInvalidArgumentErrorFromError(fmt.Errorf("contextx.PersonIDFromContext: %w", err))
	}

	if err = s.user.RemoveUser(ctx, personID); err != nil {
		return fmt.Errorf("user.RemoveUser: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
