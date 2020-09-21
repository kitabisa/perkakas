package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"
)

var testReqIDHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("test")
})

func TestRequestIDToContextAndLogMiddleware(t *testing.T) {
	handlerToTest := RequestIDToContextAndLogMiddleware(testReqIDHandler)
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
	// log.Output(os.Stdout)

	assert.Contains(t, out.String(), reqID)
}
