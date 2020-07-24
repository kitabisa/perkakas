package log

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kitabisa/perkakas/v2/internal"
)

func TestSublogger(t *testing.T) {
	ctx := context.WithValue(context.Background(), internal.CtxXKtbsRequestID, "111111")
	f, err := os.Create("logfile-test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	subLog := GetSublogger(ctx, "test-context-name-1")
	subLog.Output(f)
	subLog.Err(errors.New("test-error")).Msg("test-message")

	subLog2 := GetSublogger(ctx, "test-context-name-2")
	subLog2.Output(f)
	subLog2.Info().Msg("test-message")

	// as log is not printed when test is success, need to uncomment line below :D
	// t.FailNow()
}
