package log

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestZlogger(t *testing.T) {
	var out bytes.Buffer
	ctx := context.WithValue(context.Background(), ctxkeys.CtxLogger, zerolog.New(&out).With().Str("X-Ktbs-Request-ID", "any-request-id").Logger())

	Zlogger(ctx).Log().Send()
	assert.Contains(t, out.String(), "X-Ktbs-Request-ID")
	assert.Contains(t, out.String(), "any-request-id")
}

func TestZloggerEmpty(t *testing.T) {
	ctx := context.Background()
	file, err := os.Create("log.txt")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove("log.txt")
	log.Logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	Zlogger(ctx).Log().Msg("anylog")
	file.Close()
	out, err := ioutil.ReadFile("log.txt")
	if err != nil {
		t.FailNow()
	}

	assert.Contains(t, string(out), "anylog")
}

func TestZloggerWithoutInitiate(t *testing.T) {
	ctx := context.Background()

	Zlogger(ctx).Log().Msg("anylog")

	// as long as there's no panic error, it's okay.
	// Code below only make sure that the log printed on os.stdout
	// assert.Equal(t, "ole", "ale")
}
