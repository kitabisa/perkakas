package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cmiddleware "github.com/go-chi/chi/middleware"
	"github.com/kitabisa/perkakas/v2/httputil"
	"github.com/kitabisa/perkakas/v2/log"
	zlog "github.com/rs/zerolog/log"
)

type HttpRequestLoggerMiddleware struct {
	logger *log.Logger
}

// TODO: to be deprecated
func NewHttpRequestLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.SetRequest(r)
			next.ServeHTTP(w, r)
			logger.Print()
		})
	}
}

// RequestLogger middleware for request logging using zerolog
func RequestLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := cmiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		var body string
		if !strings.Contains(r.Header.Get("Content-type"), "multipart/form-data") {
			body = httputil.ReadRequestBody(r)
			if body != "" {
				bodyClean := new(bytes.Buffer)
				if err := json.Compact(bodyClean, []byte(body)); err != nil {
					zlog.Err(err).Send()
				}

				body = bodyClean.String()
				httputil.ExcludeSensitiveRequestBody(&body)

			}
		}

		next.ServeHTTP(ww, r)

		if ww.Status() < http.StatusBadRequest {
			return
		}

		subLog := zlog.With().
			Str(log.FieldEndpoint, r.URL.String()).
			Str(log.FieldMethod, r.Method).
			Int(log.FieldHTTPStatus, ww.Status()).
			Logger()

		if body != "" {
			subLog = subLog.With().Str(log.FieldRequestBody, body).Logger()
		}

		h := r.Header.Clone()
		h.Del("Authorization")

		var hStr []string
		for k, v := range h {
			hStr = append(hStr, fmt.Sprintf("%s: %s", k, v))
		}
		subLog = subLog.With().Str(log.FieldRequestHeaders, strings.Join(hStr, "|")).Logger()

		subLog.Info().Send()
	}
	return http.HandlerFunc(fn)
}
