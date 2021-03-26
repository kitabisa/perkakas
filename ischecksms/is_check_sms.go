package ischecksms

import "strings"

// Prevent Sms For Android/IOS
func isSendingSMS(donationSource string) bool {
	donationSource = strings.ToLower(donationSource)
	return !(donationSource == "android" || donationSource == "ios")
}
