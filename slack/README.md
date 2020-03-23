# Perkakas Slack Library
This library helps you to send message to slack via simple webhook

# How To Send Message
```go
slackMsg := NewWebhook("https://hooks.slack.com/services/your-webhook-url-path")
slackMsg.AddField("Name", "Mau tau aja")
slackMsg.AddField("Phonenumber", "081211113333")
slackMsg.AddField("Email", "qwe@testing.com")
slackMsg.AddText("Hello from Tuman-app staging")
slackMsg.SetChannelMention(true)

err := slackMsg.Send()
if err != nil {
    fmt.Println(err.Error())
}
```

The code will produce this kind of message:

![sample message img](https://user-images.githubusercontent.com/9508513/77311868-83c97880-6d33-11ea-9f0a-b868ec4dfaa7.png)
