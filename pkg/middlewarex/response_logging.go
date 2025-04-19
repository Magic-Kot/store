package middlewarex

import (
	"bytes"
	"cmp"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/zenazn/goji/web/mutil"
)

// The trouble with optional interfaces:
// https://blog.merovius.de/posts/2017-07-30-the-trouble-with-optional-interfaces/
// https://medium.com/@cep21/interface-wrapping-method-erasure-c523b3549912
func ResponseLogging(
	sensitiveDataMasker sensitiveDataMasker,
	logFieldMaxLen int,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()
			lw := mutil.WrapWriter(w)

			var buf bytes.Buffer

			lw.Tee(&buf)

			next.ServeHTTP(lw, r)

			responseHeaders, err := responseHeaders(w)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("responseHeaders")
			}

			dump := buf.Bytes()

			if len(dump) > logFieldMaxLen {
				dump = dump[:logFieldMaxLen]
			}

			// Если в хэндлере принудительно не установлен статус, то
			// lw.Status() будет возвращать 0 (упоминание этого есть в
			// документации). Поэтому устанавливаем статус 200 вручную.
			status := cmp.Or(lw.Status(), http.StatusOK)

			zerolog.Ctx(ctx).Info().
				Str("http-response", "").
				Int("response-status", status).
				Bytes("response-headers", sensitiveDataMasker.Mask(responseHeaders)).
				Bytes("response-body", sensitiveDataMasker.Mask(dump)).
				Int64("duration-ms", time.Since(start).Milliseconds()).Msg("http-response")
		})
	}
}

func responseHeaders(w http.ResponseWriter) ([]byte, error) {
	var buf bytes.Buffer

	if err := w.Header().WriteSubset(&buf, nil); err != nil {
		return nil, fmt.Errorf("header.WriteSubset: %w", err)
	}

	return buf.Bytes(), nil
}
