package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestZlogger(t *testing.T) {
	var out bytes.Buffer
	ctx := context.WithValue(context.Background(), ctxkeys.CtxLogger, zerolog.New(&out).With().Str("X-Ktbs-Request-ID", "any-request-id").Logger())

	Zlogger(ctx).Log().Send()
	assert.Contains(t, out.String(), "X-Ktbs-Request-ID")
	assert.Contains(t, out.String(), "any-request-id")
}
