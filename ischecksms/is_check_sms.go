package ischecksms

// Prevent Sms For Android/IOS
func isCheckSms(donationSource string) bool {
	return !(donationSource == "android" || donationSource == "ios")
}
