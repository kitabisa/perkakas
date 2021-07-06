package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistSuccess(t *testing.T) {
	platformName := "pwa"
	assert.True(t, IsExist(platformName, []string{"pwa", "android"}))
}

func TestIsExistSuccessUpperCase(t *testing.T) {
	platformName := "iOS"
	assert.True(t, IsExist(platformName, []string{"pwa", "ios"}))
}

func TestIsExistNotFoundFail(t *testing.T) {
	platformName := "KanvasHebat"
	assert.False(t, IsExist(platformName, []string{"pwa", "iOS"}))
}
