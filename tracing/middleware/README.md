# Middleware Handler Tracing

## Usage

* [Default middleware](#default-middleware)
* [Enable tracing logger](#enable-tracing-logger)
* [Enable tracing logger by http status](#enable-tracing-logger-by-http-status)

## Default middleware

Default handler tracing will create opentracing span for handler level without logging any traceID to console.

```go
handlerTracing := middleware.NewHandlerTracing()

r := chi.NewRouter()
r.Use(handlerTracing)

// put your routings below
```

## Enable tracing logger

Default handler tracing will create opentracing span for handler level and log its traceID to console.

```go
handlerTracing := middleware.NewHandlerTracing(
    middleware.WithEnabledLog(true),
)

r := chi.NewRouter()
r.Use(handlerTracing)

// put your routings below
```

Example result:

![Screen Shot 2020-12-11 at 17 11 43](https://user-images.githubusercontent.com/9508513/101891410-9000a900-3bd4-11eb-8e38-24ea155495b7.png)

## Enable tracing logger by http status

Default handler tracing will create opentracing span for handler level, but only log traceID if its http status is equal above specified limit.

```go
handlerTracing := middleware.NewHandlerTracing(
    middleware.WithEnabledLog(true),
    middleware.WithLimitLogHTTPStatus(http.StatusBadRequest), // this sample only print traceID log if htt status is >= 400
)

r := chi.NewRouter()
r.Use(handlerTracing)

// put your routings below
```

Example result:

![Screen Shot 2020-12-11 at 17 14 36](https://user-images.githubusercontent.com/9508513/101891486-ad357780-3bd4-11eb-8383-7d34ceed7fc1.png)