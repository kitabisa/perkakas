package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kitabisa/perkakas/v2/httputil"
	"github.com/kitabisa/perkakas/v2/log"
	zlog "github.com/rs/zerolog/log"
)

type HttpRequestLoggerMiddleware struct {
	logger *log.Logger
}

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
		subLog := zlog.With().
			Str(log.FieldEndpoint, r.URL.String()).
			Str(log.FieldMethod, r.Method).
			Logger()

		h := r.Header.Clone()
		h.Del("Authorization")

		var hStr []string
		for k, v := range h {
			hStr = append(hStr, fmt.Sprintf("%s: %s", k, v))
		}
		subLog = subLog.With().Str(log.FieldRequestHeaders, strings.Join(hStr, "|")).Logger()

		if !strings.Contains(r.Header.Get("Content-type"), "multipart/form-data") {
			body := httputil.ReadRequestBody(r)

			bodyClean := new(bytes.Buffer)
			if err := json.Compact(bodyClean, []byte(body)); err != nil {
				subLog.Err(err).Send()
			}

			body = bodyClean.String()
			httputil.ExcludeSensitiveRequestBody(&body)

			if body != "" {
				subLog = subLog.With().Str(log.FieldRequestBody, body).Logger()
			}
		}

		subLog.Info().Send()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
