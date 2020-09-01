package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kitabisa/perkakas/v2/httpclient"
	"github.com/Shopify/sarama"
)

// ISseClient defines interface of SSE client
type ISseClient interface {
	SendEvent(ctx context.Context, eventPath string, payload interface{}) (err error)
}

// Client defines object for SSE instance client
type Client struct {
	Host       string
	Username   string
	Password   string
	HTTPClient *httpclient.HttpClient
	Producer sarama.AsyncProducer
}

// NewSseClient initializes new instance of SSE client
func NewSseClient(host, username, password string, producer sarama.AsyncProducer) ISseClient {
	return &Client{
		Host:       host,
		Username:   username,
		Password:   password,
		HTTPClient: httpclient.NewHttpClient(nil),
		Producer: producer,
	}
}

// SendEvent will send event request to kitabisa SSE server
func (s *Client) SendEvent(ctx context.Context, eventPath string, payload interface{}) (err error) {
	url := fmt.Sprintf("%s%s", s.Host, eventPath)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(s.Username, s.Password)

	resp, err := s.HTTPClient.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[sse-client] %d Reader error when reading response: %s", resp.StatusCode, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[sse-client] %d:%s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *Client) PublishEvent(ctx context.Context, topic string, key string, payload interface{}) (err error) {
	if s.Producer == nil {
		return errors.New("[sse-client]: Want to publish message but producer is uninitialized")
	}

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
	case s.Producer.Input() <- msg:
	case _ = <-s.Producer.Errors():
	}

	return nil
}
