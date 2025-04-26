package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.appkode.ru/pub/go/failure"

	"github.com/Magic-Kot/store/pkg/errcodes"
)

func readRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return failure.NewInvalidArgumentError(
			fmt.Errorf("json.Decode: %w", err).Error(),
			failure.WithCode(errcodes.ValidationError),
		)
	}

	return nil
}
