# Perkakas Slack Library
This library helps you to send message to slack via simple webhook

# How To Send Message
```go
slackMsg := NewWebhook("https://hooks.slack.com/services/your-webhook-url-path")
slackMsg.AddField("Name", "Mau tau aja")
slackMsg.AddField("Phonenumber", "081211113333")
slackMsg.AddField("Email", "qwe@testing.com")
slackMsg.AddText("Hello from Tuman-app staging")

err := slackMsg.Send()
if err != nil {
    fmt.Println(err.Error())
}
```
