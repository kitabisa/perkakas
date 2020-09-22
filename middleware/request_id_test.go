package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	plog "github.com/kitabisa/perkakas/v2/log"

	uuid "github.com/satori/go.uuid"
)

var testReqIDHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	plog.Zlogger(r.Context()).Info().Msg("test")
})

func TestRequestIDToContextAndLogMiddleware(t *testing.T) {
	handlerToTest := RequestIDToContextAndLogMiddleware(RequestLogger(testReqIDHandler))
	ts := httptest.NewServer(handlerToTest)
	defer ts.Close()

	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.FailNow()
	}

	reqID := uuid.NewV4().String()
	req.Header.Add("X-Ktbs-Request-ID", reqID)

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.FailNow()
	}

	assert.Contains(t, out.String(), reqID)
}

func TestRaceRequestIDToContextAndLogMiddleware(t *testing.T) {
	handlerToTest := RequestIDToContextAndLogMiddleware(RequestLogger(testReqIDHandler))
	ts := httptest.NewServer(handlerToTest)
	defer ts.Close()
	file, err := os.Create("log.txt")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove("log.txt")
	log.Logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	maxConcurrent := 10
	sem := make(chan int, maxConcurrent)
	var wg sync.WaitGroup
	var reqIDs []string

	for i := 0; i < 1000; i++ {
		sem <- 1
		wg.Add(1)

		go func(*httptest.Server) {
			defer wg.Done()
			// defer pipeWriter.Close()
			reqID := callHTTP(ts)
			reqIDs = append(reqIDs, reqID)
			<-sem
		}(ts)
	}

	wg.Wait()
	file.Close()
	out, err := ioutil.ReadFile("log.txt")
	if err != nil {
		t.FailNow()
	}

	// the test is still fail, always lack of 1 line from total looping on the our buffer.
	for _, reqID := range reqIDs {
		assert.Contains(t, string(out), reqID)
		// assert.Contains(t, "123", reqID)
	}
}

func callHTTP(ts *httptest.Server) string {
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		log.Err(err).Send()
		return ""
	}

	reqID := uuid.NewV4().String()
	req.Header.Add("X-Ktbs-Request-ID", reqID)

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Err(err).Send()
		return ""
	}

	return reqID
}
