package middlewarex

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/Magic-Kot/store/pkg/contextx"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID, err := contextx.TraceIDFromContext(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("contextx.TraceIDFromContext")
		}

		ctx = contextx.WithLogger(
			ctx,
			zerolog.Ctx(ctx).With().
				Str("trace-id", traceID.String()).
				Str("url", r.URL.Path).
				Str("http-method", r.Method).
				Str("ip", r.RemoteAddr).
				Logger(),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
