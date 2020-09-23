package middleware

import (
	"context"
	"net/http"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog/log"
)

// RequestIDToContextAndLogMiddleware set X-Ktbs-Request-ID header value and logger to context
func RequestIDToContextAndLogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := r.Header.Get(ctxkeys.CtxXKtbsRequestID.String())
		r = r.WithContext(context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID))

		logger := log.With().
			Str(ctxkeys.CtxXKtbsRequestID.String(), reqID).
			Logger()
		r = r.WithContext(context.WithValue(ctx, ctxkeys.CtxLogger, logger))

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
