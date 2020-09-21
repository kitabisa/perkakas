# Kafka
This package is a helper to create kafka sync and/or async producer, using sarama under the hood.

## How To Create Producers
```go
package main

func main() {
    // Create synchronous producer
    producer, err := NewKafkaProducer([]string{"localhost:9092"}, "2.5.0", WithClientID("dhuwit"), WithRetryMax(5))
    if err != nil {
        panic(err)
    }

    // Create asynchronous producer
    asyncProducer, err := NewKafkaAsyncProducer([]string{"localhost:9092"}, "2.5.0", WithClientID("dhuwit"), WithRetryMax(5))
    if err != nil {
        panic(err)
    }
}
```

When create the producer, it will set some default values.

## Options
When creating producer, this helper using options pattern, so you can inject the options as needed with `WithXXX()`
functions. `WithXXX()` functions have `ProducerConfigOption` type. It will act as configuration properties that sarama
producer has.
Right now `ProducerConfigOption` only provides kafka producer configuration options that often used.
For perkakas maintainers, you can create additional `ProducerConfigOption` option in the future by create `WithXXX()`
functions, based on the sarama producer configuration as needed.

### Options Available
`WithClientID(clientID string)` -- Set kafka producer client ID

`WithMaxMessageBytes(maxMessageBytes int)` -- Set a single message size in bytes

`WithoutSASL()` -- Disable SASL support

`WithRequiredAcks(mode int)` -- Set the required acks mode

`WithRetryBackoff(d time.Duration)` -- Set the delay between retries if retry happens

`WithRetryMax(max int)` -- Set the max retry attempt

`WithReturnErrors(t bool)` -- Enable/disable the kafka producer to return error to the error channel

`WithSuccesses(t bool)` -- Enable/disable the kafka producer to return successfully delivered messages to Successes channel.
Always `on` on synchronous producer

`WithSASL(username, password string)` -- Set SASL support to the producer. Takes SASL username and password

`WithTLS(t bool)` -- Enable/disable TLS support when connecting to kafka brokers

`WithVerbose()` -- Set the kafka producer logger, discard by default
