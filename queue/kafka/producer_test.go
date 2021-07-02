package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWithClientID(t *testing.T) {
	cfg := sarama.NewConfig()
	WithClientID("katresnan")(cfg)
	assert.Equal(t, "katresnan", cfg.ClientID)
}

func TestWithMaxMessageBytes(t *testing.T) {
	maxBytes := 10000
	cfg := sarama.NewConfig()
	WithMaxMessageBytes(maxBytes)(cfg)
	assert.Equal(t, maxBytes, cfg.Producer.MaxMessageBytes)
}

func TestWithoutSASL(t *testing.T) {
	cfg := sarama.NewConfig()
	WithoutSASL()(cfg)
	assert.Equal(t, false, cfg.Net.SASL.Enable)
	assert.Equal(t, "", cfg.Net.SASL.User)
	assert.Equal(t, "", cfg.Net.SASL.Password)
}

func TestWithRequiredAcks(t *testing.T) {
	cfg := sarama.NewConfig()
	WithRequiredAcks(sarama.WaitForAll)(cfg)
	assert.Equal(t, sarama.WaitForAll, cfg.Producer.RequiredAcks)
}

func TestWithRetryBackoff(t *testing.T) {
	backoffDuration := 3 * time.Second
	cfg := sarama.NewConfig()
	WithRetryBackoff(backoffDuration)(cfg)
	assert.Equal(t, backoffDuration, cfg.Producer.Retry.Backoff)
}

func TestWithRetryMax(t *testing.T) {
	maxRetry := 3
	cfg := sarama.NewConfig()
	WithRetryMax(maxRetry)(cfg)
	assert.Equal(t, maxRetry, cfg.Producer.Retry.Max)
}

func TestWithReturnErrors(t *testing.T) {
	cfg := sarama.NewConfig()
	WithReturnErrors(true)(cfg)
	assert.Equal(t, true, cfg.Producer.Return.Errors)
}

func TestWithReturnSuccesses(t *testing.T) {
	cfg := sarama.NewConfig()
	WithReturnSuccesses(true)(cfg)
	assert.Equal(t, true, cfg.Producer.Return.Successes)
}

func TestWithSASL(t *testing.T) {
	username := "myusername"
	password := "mypassword"

	cfg := sarama.NewConfig()
	WithSASL(username, password)(cfg)
	assert.Equal(t, true, cfg.Net.SASL.Enable)
	assert.Equal(t, username, cfg.Net.SASL.User)
	assert.Equal(t, password, cfg.Net.SASL.Password)
}

func TestWithTLS(t *testing.T) {
	cfg := sarama.NewConfig()
	WithTLS(true)(cfg)
	assert.Equal(t, true, cfg.Net.TLS.Enable)
}

func TestWithVerbose(t *testing.T) {
	cfg := sarama.NewConfig()
	WithVerbose()(cfg)
	sarama.Logger.Println("hello")
	// Output: hello
}

func TestNewKafkaProducerConfig(t *testing.T) {
	cfg, err := NewKafkaProducerConfig("2.5.0")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestNewKafkaProducerConfigVersionError(t *testing.T) {
	cfg, err := NewKafkaProducerConfig("1234")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestNewKafkaProducer(t *testing.T) {
	mockBroker := sarama.NewMockBroker(t, 1)
	mockBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t),
	})

	producer, err := NewKafkaProducer([]string{mockBroker.Addr()}, "2.5.0", WithClientID("dhuwit"), WithRetryMax(5))
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewKafkaAsyncProducer(t *testing.T) {
	mockBroker := sarama.NewMockBroker(t, 2)
	mockBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t),
	})

	producer, err := NewKafkaAsyncProducer([]string{mockBroker.Addr()}, "2.5.0", WithClientID("dhuwit"), WithRetryMax(5))
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}
