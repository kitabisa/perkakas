# Perkakas SSE Library
This library helps you to send message to SSE server

## How To Send Message
```go
package main

import (
    "context"
    "fmt"
    "github.com/kitabisa/perkakas/v2/sse"
    "github.com/kitabisa/perkakas/v2/queue/kafka"
)

func main() {
    data := map[string]interface{} {
        "donation_id": 5234577,
        "user_id": 267182,
        "amount": 50000,
    }

    // By default, SSE client will set the kafka version to 2.5.0. You can change the kafka version using
    // `SetKafkaVersion()` and see the kafka version currently applied with `GetKafkaVersion()`
    kafkaHost := []string{"localhost:9092"}
    client := sse.NewSseClient(kafkaHost, kafka.WithClientID("katresnan"), kafka.WithRetryMax(5))

    err := client.PublishEvent(context.Background(), "topic", "key", data)
    if err != nil {
        fmt.Println(err)
    }
}
```
