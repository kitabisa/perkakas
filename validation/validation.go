package validation

import "strings"

// Prevent Sms For Android/IOS
func IsNotMobileApp(source string) bool {
	source = strings.ToLower(source)
	return !(source == "android" || source == "ios")
}
