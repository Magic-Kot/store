package middlewarex

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func() {
			if rec := recover(); rec != nil {
				zerolog.Ctx(ctx).Error().Str("panic in handler", "").Any("error", rec).Bytes("stack", debug.Stack()).Msg("recovered from panic")

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
