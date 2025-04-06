package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"git.appkode.ru/pub/go/failure"
	"github.com/rs/zerolog"

	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/errcodes"
)

func replyJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(fmt.Sprintf("json.Encode: %v", err))
	}
}

//nolint:funlen,gocognit,gocyclo,cyclop
func replyError(ctx context.Context, w http.ResponseWriter, err error) {
	zerolog.Ctx(ctx).Error().Err(err).Msg(err.Error())

	var statusCode int

	response := value.ErrorModel{
		Code:    failure.Code(err).String(),
		Message: failure.Description(err),
	}

	switch {
	case failure.IsInvalidArgumentError(err):
		statusCode = http.StatusBadRequest

		if failure.Code(err) == "" {
			response.Code = errcodes.ValidationError.String()
		}

		if response.Message == "" {
			response.Message = errcodes.ValidationErrorMessage
		}

	case failure.IsNotFoundError(err):
		statusCode = http.StatusNotFound

		if response.Code == "" {
			response.Code = errcodes.NotFound.String()
		}

		if response.Message == "" {
			response.Message = errcodes.NotFoundErrorMessage
		}

	case failure.IsUnprocessableEntityError(err):
		statusCode = http.StatusUnprocessableEntity

	case failure.IsUnauthorizedError(err):
		statusCode = http.StatusUnauthorized

	case failure.IsForbiddenError(err):
		statusCode = http.StatusForbidden

		if response.Code == "" {
			response.Code = errcodes.Forbidden.String()
		}

		if response.Message == "" {
			response.Message = errcodes.ForbiddenErrorMessage
		}

	case failure.IsConflictError(err):
		statusCode = http.StatusConflict

	default:
		statusCode = http.StatusInternalServerError

		if response.Code == "" {
			response.Code = errcodes.InternalServerError.String()
		}

		if response.Message == "" {
			response.Message = errcodes.InternalServerErrorMessage
		}
	}

	switch {
	case failure.IsInvalidArgumentError(err):
		replyJSON(ctx, w, statusCode, response)
	case failure.IsNotFoundError(err):
		replyJSON(ctx, w, statusCode, response)
	case failure.IsUnauthorizedError(err):
		replyJSON(ctx, w, statusCode, response)
	case failure.IsForbiddenError(err):
		replyJSON(ctx, w, statusCode, response)
	case failure.IsConflictError(err):
		replyJSON(ctx, w, statusCode, response)
	case failure.IsUnprocessableEntityError(err):
		replyJSON(ctx, w, statusCode, response)
	default:
		replyJSON(ctx, w, http.StatusInternalServerError, response)
	}
}
