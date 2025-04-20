package rest

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/store/pkg/contextx"

	"git.appkode.ru/pub/go/failure"
)

func (s Server) GetUserV1Info(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	personID, err := contextx.PersonIDFromContext(ctx)
	if err != nil {
		return failure.NewInvalidArgumentErrorFromError(fmt.Errorf("contextx.PersonIDFromContext: %w", err))
	}

	result, err := s.user.UserInfo(ctx, personID)
	if err != nil {
		return fmt.Errorf("user.UserInfo: %w", err)
	}

	replyJSON(r.Context(), w, http.StatusOK, result)

	return nil
}
