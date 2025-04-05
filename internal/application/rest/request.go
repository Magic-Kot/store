package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.appkode.ru/pub/go/failure"
)

func readRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return failure.NewInvalidArgumentError(
			fmt.Errorf("json.Decode: %w", err).Error(),
			failure.WithCode(failure.ErrorCode(VALIDATION_ERROR)),
		)
	}

	return nil
}
