# Open Tracing Context

> Wrapping Service Dependencies with opentracing

* [Wraping Handler](#wrapping-handler)
* [Wraping go-redis](#wrapping-go-redis)

## Wraping Handler

This middleware wraps handler with opentracing context along with any identified request_id, path params, and query params.

```go
handlerTracing := middleware.NewHandlerTracing()

r := chi.NewRouter()
r.Use(handlerTracing)

// put your routings below
```

![](https://user-images.githubusercontent.com/9508513/101141474-ed31b300-3646-11eb-937a-561cf12c1c17.png)

Plase read [middleware readme.md](https://github.com/kitabisa/perkakas/blob/master/tracing/middleware/README.md) to know more about handler tracing middleware usage example.

## Wraping go-redis

```go
redisClient := redis.WrapRedisClient(ctx, redisClient)
```
