# Perkakas SSE Library
This library helps you to send message to SSE server

# How To Send Message
```go
type Donation struct {
    DonationID int64 `json:"donation_id"`
    UserID     int64 `json:"user_id"`
    Amount     int64 `json:"amount"` 
}

donationData := Donation{
    DonationID: 334534,
    UserID:     2463,
    Amount:     50000,
}

sseClient := NewSseClient("your-sse-host", "username", "password")

err := sseClient.SendEvent(context.Background(), "/notification/donation-created", donationData)
if err != nil {
    panic(err)
}
```
