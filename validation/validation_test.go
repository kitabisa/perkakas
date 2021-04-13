package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotMobileAppSuccess(t *testing.T) {
	source := "pwa"
	assert.True(t, IsNotMobileApp(source))
}

func TestIsNotMobileAppFail(t *testing.T) {
	source := "Android"
	assert.False(t, IsNotMobileApp(source))
}
