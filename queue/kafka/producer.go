package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

// NewKafkaProducerConfig is a function to initialize default kafka producer configuration
// It accepts version of kafka and options. This options are taken from sarama configuration as
// ProducerConfigOption, but not all covered. ProducerConfigOption only provides kafka producer
// configuration options that often used. For perkakas maintainers, you can add the ProducerConfigOption
// provided in the future by create WithXXX() functions, based on the sarama producer configuration.
func NewKafkaProducerConfig(version string, opts ...ProducerConfigOption) (*sarama.Config, error) {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = kafkaVersion

	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}

// ProducerConfigOption is an function type alias that accepts *sarama.Config. It will be used as producer configuration
// options type.
type ProducerConfigOption func(*sarama.Config)

// WithClientID sets the client id in producer configuration
func WithClientID(clientID string) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.ClientID = clientID
	}
}

// WithVerbose sets the logger to output the log into stdout. If you not set this, the producer will discard the log.
func WithVerbose() ProducerConfigOption {
	return func(_ *sarama.Config) {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}
}

// WithRetryMax sets max retry if producer failed to produce the message. Default 3 if you not set this.
func WithRetryMax(max int) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.Retry.Max = max
	}
}

// WithRetryBackoff sets duration between retry. Default 100 ms.
func WithRetryBackoff(duration time.Duration) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.Retry.Backoff = duration
	}
}

// WithReturnSuccesses tells producers to return success. Default disabled.
func WithReturnSuccesses(isReturnSuccess bool) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.Return.Successes = isReturnSuccess
	}
}

// WithReturnErrors tells the producers to return error. Default enabled.
func WithReturnErrors(isReturnError bool) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.Return.Errors = isReturnError
	}
}

// WithRequiredAcks sets the level of aknowledement reliability. Default to WaitForLocal.
func WithRequiredAcks(reqAcks sarama.RequiredAcks) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.RequiredAcks = reqAcks
	}
}

// WithMaxMessageBytes sets the max permitted size of message. Defaults to 1000000.
func WithMaxMessageBytes(max int) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Producer.MaxMessageBytes = max
	}
}

// WithTLS sets producer to connect with tls or not. Default false/disabled.
func WithTLS(withTLS bool) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Net.TLS.Enable = withTLS
	}
}

// WithSASL sets the expected username and password for SASL and enabling the SASL mode in producer configuration.
func WithSASL(user, pass string) ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Net.SASL.Enable = true
		c.Net.SASL.User = user
		c.Net.SASL.Password = pass
	}
}

// WithoutSASL disable SASL mode in producer and emptying SASL username and password configuration.
func WithoutSASL() ProducerConfigOption {
	return func(c *sarama.Config) {
		c.Net.SASL.Enable = false
		c.Net.SASL.User = ""
		c.Net.SASL.Password = ""
	}
}

// NewKafkaProducer initialize synchronous producer for kafka.
func NewKafkaProducer(brokers []string, version string, opts ...ProducerConfigOption) (producer sarama.SyncProducer, err error) {
	opts = append(opts, WithReturnSuccesses(true))
	config, err := NewKafkaProducerConfig(version, opts...)
	if err != nil {
		return
	}

	producer, err = sarama.NewSyncProducer(brokers, config)
	return
}

// NewKafkaAsyncProducer initialize asynchronous producer for kafka.
func NewKafkaAsyncProducer(brokers []string, version string, opts ...ProducerConfigOption) (producer sarama.AsyncProducer, err error) {
	config, err := NewKafkaProducerConfig(version, opts...)
	if err != nil {
		return
	}

	producer, err = sarama.NewAsyncProducer(brokers, config)
	return
}
