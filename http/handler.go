package http

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	zlog "github.com/kitabisa/perkakas/v2/log"
)

type HandlerOption func(*HttpHandler)

type HttpHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)
	CustomWriter
	Metric      *statsd.Client
	ServiceName string
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
func WithMetric(m *statsd.Client, svcName string) HandlerOption {
	return func(h *HttpHandler) {
		h.Metric = m
		h.ServiceName = svcName
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startHandleRequest := time.Now()
	data, pageToken, err := h.H(w, r)
	diff := time.Since(startHandleRequest)

	var tag []string
	URL := PathPattern(r.URL.Path)

	if err != nil {
		if h.Metric != nil {
			var table string = "ERROR"
			var statusCode int
			var responseCode string
			if erResp, ok := h.C.E[err]; ok {
				statusCode = erResp.HttpStatus
				responseCode = erResp.Response.ResponseCode
			}

			var status string
			if statusCode >= 400 && statusCode < 500 {
				status = "CLIENT_ERROR"
			} else {
				status = "SERVER_ERROR"
			}

			tag = append(tag, fmt.Sprintf("service_name:%s", h.ServiceName), fmt.Sprintf("endpoint:%s", URL), fmt.Sprintf("http_status:%d", statusCode), fmt.Sprintf("response_code:%s", responseCode), fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID")), fmt.Sprintf("status:%s", status))

			h.Metric.Incr(table, tag, 1)
		}

		zlog.Zlogger(r.Context()).Err(err).Msgf("Response: %+v", data)
		h.WriteError(w, err)
		return
	}

	if h.Metric != nil {
		tag = append(tag, fmt.Sprintf("service_name:%s", h.ServiceName), fmt.Sprintf("endpoint:%s", r.URL.Path), fmt.Sprintf("http_status:%d", 200), fmt.Sprintf("response_code:%s", "000000"), fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID")))

		h.Metric.Incr("SUCCESS", tag, 1)

		// response time
		responseTimeTag := []string{fmt.Sprintf("service_name:%s", h.ServiceName), fmt.Sprintf("endpoint:%s", r.URL.Path), fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID"))}

		h.Metric.Incr("RESPONSE_TIME", responseTimeTag, float64(diff.Milliseconds()))
	}

	h.Write(w, data, pageToken)
}

const paramSign = "PARAM"

// PathPattern modify params on url
func PathPattern(input string) string {
	u, _ := url.Parse(input)
	path := u.Path
	for _, pic := range strings.Split(path, "/") {
		itParam, _ := regexp.Match("[0-9][a-z+A-Z]*[0-9]", []byte(pic))
		if itParam {
			path = strings.Replace(path, pic, paramSign, 1)
		} else if len(pic) > 12 {
			path = strings.Replace(path, pic, paramSign, 1)
		}
	}

	return path
}
