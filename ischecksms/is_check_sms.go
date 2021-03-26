package ischecksms

import "strings"

// Prevent Sms For Android/IOS
func IsSendingSMS(donationSource string) bool {
	donationSource = strings.ToLower(donationSource)
	return !(donationSource == "android" || donationSource == "ios")
}
