package influx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func emptyTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("X-Influxdb-Version", "x.x")
	}))
}

func TestWritePoint(t *testing.T) {
	mockServer := emptyTestServer()

	config := ClientConfig{
		Addr:     mockServer.URL,
		Database: "myDB",
		Timeout:  5 * time.Second,
	}

	tags := Tags{
		"tag1": "tagsValue1",
		"tag2": "tagsValue2",
	}

	fields := Fields{
		"fields1": 1,
		"fields2": "value2",
	}

	c, err := NewClient(config)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = c.WritePoints("dummy", tags, fields, "ns")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestWriteBatchPoint(t *testing.T) {
	mockServer := emptyTestServer()

	config := ClientConfig{
		Addr:     mockServer.URL,
		Database: "myDB",
		Timeout:  5 * time.Second,
	}

	tags := Tags{
		"tag1": "tagsValue1",
		"tag2": "tagsValue2",
	}

	fields := Fields{
		"fields1": 1,
		"fields2": "value2",
	}

	c, err := NewClient(config)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	b, err := c.NewBatchPointsWriter("s")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	b.AddPoints("dummy", tags, fields) // you can add more points later

	b.Write() // finally write the points
}
