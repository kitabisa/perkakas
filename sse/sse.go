package sse

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/kitabisa/perkakas/v2/queue/kafka"
)

// ISseClient defines interface of SSE client
type ISseClient interface {
	// PublishEvent has functionality to publish event to kafka brokers
	PublishEvent(ctx context.Context, topic string, key string, payload interface{}) (err error)

	// SetKafkaVersion sets the kafka version
	SetKafkaVersion(ctx context.Context, version string)

	// GetKafkaVersion gets the kafka version
	GetKafkaVersion(ctx context.Context) string
}

// Client defines object for SSE instance client
type Client struct {
	// Host of the kafka brokers. Currently designed for one broker.
	Host       string

	// KafkaVersion denotes the expecting kafka version used by this client.
	KafkaVersion    string

	// opts are array of producer config options use to build the kafka producer used by this client.
	opts []kafka.ProducerConfigOption
}

// NewSseClient initializes new instance of SSE client
func NewSseClient(host string, opts ...kafka.ProducerConfigOption) ISseClient {
	return &Client{
		Host:       host,
		KafkaVersion: "2.5.0",
		opts: opts,
	}
}

// SetKafkaVersion sets the kafka version
func (s *Client) SetKafkaVersion(ctx context.Context, version string) {
	s.KafkaVersion = version
}

// GetKafkaVersion gets the kafka version
func (s *Client) GetKafkaVersion(ctx context.Context) string {
	return s.KafkaVersion
}

// PublishEvent has functionality to publish event to kafka brokers.
// For payload, you can use marshalable types, such as struct or map[string]interface{}. PublishEvent will publish
// message to kafka broker in asynchronous fashion.
func (s *Client) PublishEvent(ctx context.Context, topic string, key string, payload interface{}) (err error) {
	producer, err := kafka.NewKafkaAsyncProducer([]string{s.Host}, s.KafkaVersion, s.opts...)
	if err != nil {
		return
	}

	defer producer.Close()

	val, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(val),
	}

	select {
	case producer.Input() <- msg:
	case <- producer.Errors():
	}

	return nil
}
