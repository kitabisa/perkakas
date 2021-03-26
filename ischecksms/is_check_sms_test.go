package ischecksms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSendingSMSIsSuccess(t *testing.T) {
	source := "pwa"
	assert.True(t, isSendingSMS(source))
}

func TestIsSendingSMSIsFail(t *testing.T) {
	source := "Android"
	assert.False(t, isSendingSMS(source))
}
