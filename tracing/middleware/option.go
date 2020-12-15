package middleware

type tracingConfig struct {
	enableLogTraceID   bool
	limitHTTPStatusLog int
}

// TracingOption defines functional object for tracing option
type TracingOption func(*tracingConfig)

// WithEnabledLog will set traceID log toggle
func WithEnabledLog(toggle bool) TracingOption {
	return func(cfg *tracingConfig) {
		cfg.enableLogTraceID = toggle
	}
}

// WithLimitLogHTTPStatus will set traceID logger if http status is equal above certain limit
func WithLimitLogHTTPStatus(httpStatus int) TracingOption {
	return func(cfg *tracingConfig) {
		cfg.limitHTTPStatusLog = httpStatus
	}
}
