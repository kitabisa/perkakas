package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	cmiddleware "github.com/go-chi/chi/middleware"
	"github.com/kitabisa/perkakas/v2/log"
	"github.com/opentracing/opentracing-go"
	otLog "github.com/opentracing/opentracing-go/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
)

const paramSign = "{param}"

// NewHandlerTracing initializes opentracing context for handler level, via middleware
func NewHandlerTracing(opts ...TracingOption) func(next http.Handler) http.Handler {
	// set default value
	tracingCfg := &tracingConfig{
		enableLogTraceID:   false,
		limitHTTPStatusLog: http.StatusOK,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			operationName, pathParams := makeOperationName(r)

			span, ctx := opentracing.StartSpanFromContext(r.Context(), operationName)
			defer span.Finish()

			for _, opt := range opts {
				opt(tracingCfg)
			}

			// add ktbs request_id to span tag
			if r.Header.Get("X-Ktbs-Request-ID") != "" {
				span.SetTag("request_id", r.Header.Get("X-Ktbs-Request-ID"))
			}

			// inject path & query params to span log
			var fields []otLog.Field
			if len(pathParams) > 0 {
				fields = append(fields, otLog.Object("path_params", pathParams))
			}

			if len(r.URL.Query()) > 0 {
				fields = append(fields, otLog.Object("query_params", r.URL.Query()))
			}

			if len(fields) > 0 {
				span.LogFields(fields...)
			}

			r = r.WithContext(ctx)
			ww := cmiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			if tracingCfg.enableLogTraceID && ww.Status() >= tracingCfg.limitHTTPStatusLog {
				if sc, ok := span.Context().(jaeger.SpanContext); ok {
					subLog := zlog.With().
						Str(log.FieldEndpoint, r.URL.String()).
						Str(log.FieldMethod, r.Method).
						Int(log.FieldHTTPStatus, ww.Status()).
						Str("trace_id", sc.TraceID().String()).
						Logger()

					if r.Header.Get("X-Ktbs-Request-ID") != "" {
						subLog = subLog.With().Str("request_id", r.Header.Get("X-Ktbs-Request-ID")).Logger()
					}
					subLog.Info().Send()
				}
			}
		})
	}
}

func makeOperationName(r *http.Request) (urlPath string, pathParams []string) {
	urlPath, pathParams = pathPattern(r.URL.Path)
	return fmt.Sprintf("%s %s", r.Method, urlPath), pathParams
}

func pathPattern(input string) (path string, arrParam []string) {
	u, _ := url.Parse(input)
	path = u.Path

	for _, pic := range strings.Split(path, "/") {
		itParam, _ := regexp.Match("[0-9][a-z+A-Z]*[0-9]", []byte(pic))
		if itParam {
			path = strings.Replace(path, pic, paramSign, 1)
			arrParam = append(arrParam, pic)
		} else if len(pic) > 12 {
			path = strings.Replace(path, pic, paramSign, 1)
			arrParam = append(arrParam, pic)
		}
	}

	return path, arrParam
}
