package tracing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunctionName(t *testing.T) {
	functionName := FunctionName()
	assert.Equal(t, "tracing.TestFunctionName", functionName)
}
