package ischecksms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCheckSmsIsSuccess(t *testing.T) {
	source := "pwa"
	assert.True(t, isSendingSMS(source))
}