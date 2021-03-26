package ischecksms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSendingSMSIsSuccess(t *testing.T) {
	source := "pwa"
	assert.True(t, IsSendingSMS(source))
}

func TestIsSendingSMSIsFail(t *testing.T) {
	source := "Android"
	assert.False(t, IsSendingSMS(source))
}
