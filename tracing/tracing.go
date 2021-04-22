package tracing

import (
	"path/filepath"
	"runtime"
)

// FunctionName get function name programmatically. It will return function name complete from the package and method
// receiver name.
// This function should called when you start span.
//
// Do this:
//span, ctx := opentracing.StartSpanFromContext(ctx, FunctionName())
//
// Instead of:
// span, ctx := opentracing.StartSpanFromContext(ctx, "Service.Brand.Create")
func FunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return filepath.Base(runtime.FuncForPC(pc).Name())
}
