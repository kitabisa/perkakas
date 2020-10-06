package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatCurrency(t *testing.T) {
	amountIDR := FormatIDR(1000000)
	assert.Equal(t, "Rp1.000.000", amountIDR)
}

func TestTranslateDateToIndonesian(t *testing.T) {
	date := time.Date(2020, time.October, 6, 10, 46, 03, 0, time.UTC)
	formattedDate := TranslateDateToIndonesian(date)
	assert.Equal(t, "Selasa, 6 Oktober 2020, 17:46 WIB", formattedDate)
}
