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
func WithMetric(telegrafHost string, telegrafPort int, svcName string) HandlerOption {
	return func(h *HttpHandler) {
		host := fmt.Sprintf("%s:%d", telegrafHost, telegrafPort)
		m, err := statsd.New(host)
		if err != nil {
			panic(err)
		}

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

		// don't calculate metrics if endpoint is health check
		if h.Metric != nil && !isHealthEndpoint(r.URL.Path, r.Method) {
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

			tag = append(tag,
				fmt.Sprintf("service_name:%s", h.ServiceName),
				fmt.Sprintf("method:%s", r.Method),
				fmt.Sprintf("endpoint:%s", URL),
				fmt.Sprintf("http_status:%d", statusCode),
				fmt.Sprintf("response_code:%s", responseCode),
				fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID")),
				fmt.Sprintf("status:%s", status),
			)

			h.Metric.Incr(table, tag, 1)
		}

		zlog.Zlogger(r.Context()).Err(err).Msgf("Response: %+v", data)
		h.WriteError(w, err)
		return
	}

	// don't calculate metrics if endpoint is health check
	if h.Metric != nil && !isHealthEndpoint(r.URL.Path, r.Method) {
		tag = append(tag, fmt.Sprintf("service_name:%s", h.ServiceName), fmt.Sprintf("endpoint:%s", URL), "http_status:200", "response_code:000000", fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID")), fmt.Sprintf("method:%s", r.Method))

		h.Metric.Incr("SUCCESS", tag, 1)

		// response time
		responseTimeTag := []string{
			fmt.Sprintf("service_name:%s", h.ServiceName),
			fmt.Sprintf("method:%s", r.Method),
			fmt.Sprintf("endpoint:%s", URL),
			fmt.Sprintf("request_id:%s", r.Header.Get("X-Ktbs-Request-ID")),
		}

		h.Metric.Count("RESPONSE_TIME", diff.Milliseconds(), responseTimeTag, 1)
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

// isHealthEndpoint determines whether the endpoint is health check or not
func isHealthEndpoint(ep, method string) bool {
	healthWords := []string{"health", "liveness", "readiness", "ready"}
	result := false

	if method != http.MethodGet {
		return result
	}

	for _, subStr := range healthWords {
		isHealth := strings.Contains(ep, subStr)
		if isHealth {
			result = true
			break
		}
	}
	return result
}
