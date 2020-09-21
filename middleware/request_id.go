package middleware

import (
	"context"
	"net/http"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog/log"
)

// RequestIDToContextMiddleware set X-Ktbs-Request-ID header value to context
func RequestIDToContextAndLogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := r.Header.Get(ctxkeys.CtxXKtbsRequestID.String())
		r = r.WithContext(context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID))

		log.Logger = log.With().
			Str(ctxkeys.CtxXKtbsRequestID.String(), reqID).
			Logger()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
