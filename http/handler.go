package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	zlog "github.com/kitabisa/perkakas/v2/log"
)

type HandlerOption func(*HttpHandler)

type HttpHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)
	CustomWriter
	Metric *statsd.Client
}

func NewHttpHandler(c HttpHandlerContext, opts ...HandlerOption) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
		h := HttpHandler{H: handler, CustomWriter: CustomWriter{C: c}}

		// Option paremeters values:
		for _, opt := range opts {
			opt(&h)
		}

		return h
	}
}

// WithMetric wire statsd client to perkakas handler
func WithMetric(m *statsd.Client) HandlerOption {
	return func(h *HttpHandler) {
		h.Metric = m
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startHandleRequest := time.Now()
	data, pageToken, err := h.H(w, r)
	finishHandleRequest := time.Now()

	diff := finishHandleRequest.Sub(startHandleRequest)

	responseTimeTag := []string{fmt.Sprintf("response_time(s):%0.1f", diff.Seconds())}

	if h.Metric != nil {
		h.Metric.Incr("RESPONSE_TIME", responseTimeTag, 1)
	}

	if err != nil {
		if h.Metric != nil {
			var statusCode int
			var responseCode string
			var tag []string
			if erResp, ok := h.C.E[err]; ok {
				statusCode = erResp.HttpStatus
				responseCode = erResp.Response.ResponseCode
			}

			tag = append(tag, fmt.Sprintf("http_status:%d", statusCode), fmt.Sprintf("response_code:%s", responseCode), fmt.Sprintf("endpoint:%s", r.URL.Path))

			fmt.Println(tag)

			var table string
			if statusCode >= 400 && statusCode < 500 {
				table = "WARN"
			} else {
				table = "ERROR"
			}

			h.Metric.Incr(table, tag, 1)
		}

		zlog.Zlogger(r.Context()).Err(err).Msgf("Response: %+v", data)
		h.WriteError(w, err)
		return
	}

	if h.Metric != nil {
		h.Metric.Incr("SUCCESS", []string{"http_status:200"}, 1)
	}

	h.Write(w, data, pageToken)
}
