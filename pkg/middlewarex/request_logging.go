package middlewarex

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/rs/zerolog"
)

func RequestLogging(
	sensitiveDataMasker sensitiveDataMasker,
	logFieldMaxLen int,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			dumpBody := true

			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				dumpBody = false
			}

			dump, _ := httputil.DumpRequest(r, dumpBody)

			if len(dump) > logFieldMaxLen {
				dump = dump[:logFieldMaxLen]
			}

			zerolog.Ctx(ctx).Info().Str("http-request", "").Bytes("request-body", sensitiveDataMasker.Mask(dump)).Msg("http-request")

			next.ServeHTTP(w, r)
		})
	}
}
