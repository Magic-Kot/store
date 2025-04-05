package middlewarex

import (
	"net/http"

	"github.com/rs/xid"

	"github.com/Magic-Kot/store/pkg/contextx"
)

const headerNameTraceID = "X-Trace-ID"

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(headerNameTraceID)

		if traceID == "" {
			traceID = xid.New().String()
		}

		ctx := contextx.WithTraceID(r.Context(), contextx.TraceID(traceID))

		w.Header().Set(headerNameTraceID, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
