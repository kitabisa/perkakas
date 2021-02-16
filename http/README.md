# Package http

This package contains http handler related to help you write handler easier and have a standardized response.

## How to create custom handler
To create custom handler using this package, you need to declare the handler using this function signature:

```go
func(w http.ResponseWriter, r *http.Request) (data interface{}, pageToken *string, err error)
```

This handler requires you to return the `response data`, `page token` and `error` if any.
`response data` is the struct for the response data
`page token` is the optional hashed or encoded page token for pagination. 
We enforce this approach for pagination since it's the most clean way to do a pagination.'
This approach support infinite scroll and traditional pagination.
`error` is the error (if any)

Example:
```go
func HelloHandler(w http.ResponseWriter, r *http.Request) (data interface{}, pageToken *string, err error) {
    // Data we will return
	data = Person{
		FirstName: "Budi",
		LastName:  "Ariyanto",
	}

    // The token
	token := "o934ywjkhk67j78sd9af=="
    pageToken = &token
    
    // error, should be something like this
    // if err != nil {
    //    return err
    //}

	return
}

func main() {
    // Meta is metadata for the response
	meta := structs.Meta{
		Version: "v1.2.3",
		Status:  "stable",
		APIEnv:  "prod-test",
	}

	// When new context handler is created, it will inject all general error map.
	// Then you should add your necessary own error to the handler.
	handlerCtx := phttp.NewContextHandler(meta)

	// add error individualy
	var ErrCustom *structs.ErrorResponse = &structs.ErrorResponse{
		Response: structs.Response{
			ResponseCode: "00011",
			ResponseDesc: structs.ResponseDesc{
				ID: "Custom error",
				EN: "Custom error",
			},
		},
		HttpStatus: http.StatusInternalServerError,
	}
	handlerCtx.AddError(errors.New("custom error"), ErrCustom)

	// add error individualy
	var ErrCustom2 *structs.ErrorResponse = &structs.ErrorResponse{
		Response: structs.Response{
			ResponseCode: "00011",
			ResponseDesc: structs.ResponseDesc{
				ID: "Custom error",
				EN: "Custom error",
			},
		},
		HttpStatus: http.StatusInternalServerError,
	}

	// add error individualy
	var ErrCustom3 *structs.ErrorResponse = &structs.ErrorResponse{
		Response: structs.Response{
			ResponseCode: "00011",
			ResponseDesc: structs.ResponseDesc{
				ID: "Custom error",
				EN: "Custom error",
			},
		},
		HttpStatus: http.StatusInternalServerError,
	}

	// add error by setup error map at first
	customError2 := errors.New("Custom error 2")
	customError3 := errors.New("Custom error 3")
	errMap := map[error]*structs.ErrorResponse{
		customError2: ErrCustom2,
		customError3: ErrCustom3,
	}

	// add error map
	handlerCtx.AddErrorMap(errMap)

    // newHandler is a function that will create function for create new custom handler with injected handler context
    newHandler := phttp.NewHttpHandler(handlerCtx)
    
    // helloHandler is the handler
	helloHandler := newHandler(HelloHandler)

	router := chi.NewRouter()
	router.Get("/hello", helloHandler.ServeHTTP)

	http.ListenAndServe(":5678", router)
}
```

From the example above, you can see you only care about the data, pageToken and error,
then this custom handler will construct the response itself.
This response are refer to [Kitabisa API response standardization](https://app.gitbook.com/@kitabisa-engineering/s/backend/standardization-1/api-response).

## How to use http metrics

This http metrics will send your http metrics (response time, error status, and success status) to telegraf.

Usage:

```go
handlerCtx := phttp.NewContextHandler(pstructs.Meta{
	Version: "v1.2.3",
	Status:  "stable",
	APIEnv:  "prod-test",
})

telegrafHost := "telegraf.localhost"
telegrafPort := 8125
serviceName := "my-service" // serviceName identifies your service in the telegraf storage (e.g., Influxdb, etc)

phandler := phttp.NewHttpHandler(
	handlerCtx,
	phttp.WithMetric(telegrafHost, telegrafPort, serviceName),
)
```
