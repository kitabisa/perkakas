package text

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

func TranslateDateToIndonesian(date time.Time) string {
	location, _ := time.LoadLocation("Asia/Jakarta")
	localDate := date.In(location)

	weekday := WeekdayId[date.Weekday()]
	day := localDate.Day()
	month := MonthId[localDate.Month()]
	year := localDate.Year()
	hour := localDate.Hour()
	minute := localDate.Minute()
	timeZone, _ := localDate.Zone()

	return fmt.Sprintf("%s, %d %s %d, %02d:%02d %s", weekday, day, month, year, hour, minute, timeZone)
}

func formatCurrency(localeID string, amount int64, useCurrencySymbol bool) (formattedCurrency string) {
	switch localeID {
	case "id_ID":
		p := message.NewPrinter(language.Indonesian)

		if useCurrencySymbol {
			formattedCurrency = p.Sprintf("Rp%d", amount)
		} else {
			formattedCurrency = p.Sprintf("%d", amount)
		}
	}

	return
}

func FormatIDR(amount int64) string {
	return formatCurrency("id_ID", amount, true)
}

var (
	WeekdayId = make(map[time.Weekday]string)
	MonthId   = make(map[time.Month]string)
)

func init() {
	WeekdayId[time.Sunday] = "Minggu"
	WeekdayId[time.Monday] = "Senin"
	WeekdayId[time.Tuesday] = "Selasa"
	WeekdayId[time.Wednesday] = "Rabu"
	WeekdayId[time.Thursday] = "Kamis"
	WeekdayId[time.Friday] = "Jumat"
	WeekdayId[time.Saturday] = "Sabtu"

	MonthId[time.January] = "Januari"
	MonthId[time.February] = "Februari"
	MonthId[time.March] = "Maret"
	MonthId[time.April] = "April"
	MonthId[time.May] = "Mei"
	MonthId[time.June] = "Juni"
	MonthId[time.July] = "Juli"
	MonthId[time.August] = "Agustus"
	MonthId[time.September] = "September"
	MonthId[time.October] = "Oktober"
	MonthId[time.November] = "November"
	MonthId[time.December] = "Desember"
}
