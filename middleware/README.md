# Package Middleware

This package contains http middleware.
The middleware are actually a pure `http.Handler`. We keep the handler in this way since it's the form that compatible with
golang standard.

## JWT Middleware
JWT middleware is middleware for checking whether token is valid or not.

## Log Middleware
Log middleware is middleware that will help logging the application. The logging prints out log from [Kitabisa log specification](https://app.gitbook.com/@kitabisa-engineering/s/backend/standardization-1/log-format).

## How To Use The Middleware
```go
func main() {
	meta := structs.Meta{
		Version: "v1.2.3",
		Status:  "stable",
		APIEnv:  "prod-test",
	}

	handlerCtx := phttp.NewContextHandler(meta)
	
	newHandler := phttp.NewHttpHandler(handlerCtx)
	helloHandler := newHandler(HelloHandler)

    // middleware initialization. Just initialized it and use.
	midd := middleware.NewJWTMiddleware(handlerCtx, []byte("abcde"))
	logMiddleware := middleware.NewHttpRequestLoggerMiddleware(logger)

	router := chi.NewRouter()
    
    // use the middleware
    router.Use(midd)
    router.Use(logMiddleware)
    router.Use(RequestLogger) // will print request log only for 4xx and 5xx HTTP status code
    
	router.Get("/", helloHandler.ServeHTTP)

	http.ListenAndServe(":5678", router)
}
```