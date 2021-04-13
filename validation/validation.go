package validation

import "strings"

// Prevent Sms For Android/IOS
func IsSourceNotMobileApp(source string) bool {
	source = strings.ToLower(source)
	return !(source == "android" || source == "ios")
}
