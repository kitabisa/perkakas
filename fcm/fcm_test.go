package fcm

import (
	"context"
	"testing"
)

func TestSendWithoutNotification(t *testing.T) {
	f, err := NewFcm("path_to_key.json")
	if err != nil {
		t.Fail()
	}

	data := map[string]string{
		"msg": "perkakas test send message to fcm",
	}

	_, err = f.SendWithoutNotification(context.Background(), "topic_test", data)
	if err != nil {
		t.Fail()
	}
}

func TestSendWithNotification(t *testing.T) {
	f, err := NewFcm("path_to_key.json")
	if err != nil {
		t.Fail()
	}

	data := map[string]string{
		"msg": "perkakas test send message to fcm",
	}

	_, err = f.Send(context.Background(), "topic_test", "this is title notification", "this is the body", "https://image.com/location/img.jpg", data)
	if err != nil {
		t.Fail()
	}
}
