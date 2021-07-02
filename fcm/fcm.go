package fcm

import (
	"context"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Fcm struct {
	fapp            *firebase.App
	messagingClient *messaging.Client
}

func NewFcm(credentialsFilePath string) (fcm *Fcm, err error) {
	f, err := newFirebase(credentialsFilePath)
	if err != nil {
		return
	}

	messagingClient, err := firebaseMessagingClient(f)
	if err != nil {
		return
	}

	fcm = &Fcm{
		fapp:            f,
		messagingClient: messagingClient,
	}

	return
}

func newFirebase(creds string) (app *firebase.App, err error) {
	opt := option.WithCredentialsFile(creds)
	return firebase.NewApp(context.Background(), nil, opt)
}

func firebaseMessagingClient(app *firebase.App) (messagingClient *messaging.Client, err error) {
	ctx := context.Background()
	return app.Messaging(ctx)
}

func (f *Fcm) send(ctx context.Context, topic, title, body, imageURL string, message map[string]string) (messageID string, err error) {
	// See documentation on defining a message payload.
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Data:  message,
		Topic: topic,
	}

	messageID, err = f.messagingClient.Send(ctx, msg)
	if err != nil {
		return
	}

	return
}

func (f *Fcm) SendWithoutNotification(ctx context.Context, topic string, message map[string]string) (messageID string, err error) {
	return f.send(ctx, topic, "", "", "", message)
}

func (f *Fcm) Send(ctx context.Context, topic, title, body, imageURL string, message map[string]string) (messageID string, err error) {
	return f.send(ctx, topic, title, body, imageURL, message)
}
