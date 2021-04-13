package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSourceNotMobileAppSuccess(t *testing.T) {
	source := "pwa"
	assert.True(t, IsSourceNotMobileApp(source))
}

func TestIsSourceNotMobileAppFail(t *testing.T) {
	source := "Android"
	assert.False(t, IsSourceNotMobileApp(source))
}
