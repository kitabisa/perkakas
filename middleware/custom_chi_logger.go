package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	cmiddleware "github.com/go-chi/chi/middleware"
)

var customChiLogger = customRequestLogger(&cmiddleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false})

// CustomChiLogger create custom chi logger-like that print request log only for 4xx and 5xx error
func CustomChiLogger(next http.Handler) http.Handler {
	return customChiLogger(next)
}

func customRequestLogger(f cmiddleware.LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := f.NewLogEntry(r)
			ww := cmiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			next.ServeHTTP(ww, cmiddleware.WithLogEntry(r, entry))

			if ww.Status() >= http.StatusBadRequest {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}
		}
		return http.HandlerFunc(fn)
	}
}
