package ischecksms

// Prevent Sms For Android/IOS
func isSendingSMS(donationSource string) bool {
	return !(donationSource == "android" || donationSource == "ios")
}
