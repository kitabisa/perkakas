package text

import (
	"testing"
	"time"
)

func TestFormatCurrency(t *testing.T) {
	FormatIDR(1000000)

	// Output: Rp1.000.000
}

func TestTranslateDateToIndonesian(t *testing.T) {
	date := time.Date(2020, time.October, 6, 10, 46, 03, 0, time.UTC)
	TranslateDateToIndonesian(date)

	// Output: Selasa, 6 Oktober 2020, 17:46 WIB
}
