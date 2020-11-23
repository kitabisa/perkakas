package http

import (
	"fmt"
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	zlog "github.com/kitabisa/perkakas/v2/log"
)

type HanlderOption func(*HttpHandler)

type HttpHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)
	CustomWriter
	Metric *statsd.Client
}

func NewHttpHandler(c HttpHandlerContext, opts ...HanlderOption) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
		h := HttpHandler{H: handler, CustomWriter: CustomWriter{C: c}}

		// Option paremeters values:
		for _, opt := range opts {
			opt(&h)
		}

		return h
	}
}

// OptionMetric wire statsd client to perkakas handler
func OptionMetric(m *statsd.Client) HanlderOption {
	return func(h *HttpHandler) {
		h.Metric = m
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add time.now
	data, pageToken, err := h.H(w, r)
	// add time.now, calculate to get response time, kirim ke metric
	if err != nil {
		var statusCode int
		if erResp, ok := h.C.E[err]; ok {
			statusCode = erResp.HttpStatus
		}

		if h.Metric != nil {
			h.Metric.Incr("ERROR", []string{fmt.Sprintf("http_status:%d", statusCode)}, 1)
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
