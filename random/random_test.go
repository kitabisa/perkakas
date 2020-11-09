package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var length int

func TestString(t *testing.T) {
	length = 10

	generator := NewRandomStringGenerator()

	randomString := generator.String(length)

	assert.Equal(t, length, len(randomString))
	assert.NotEmpty(t, randomString)
}

func TestStringWithZeroLength(t *testing.T) {
	length = 0

	generator := NewRandomStringGenerator()

	randomString := generator.String(length)

	assert.Equal(t, randomString, "")
}
