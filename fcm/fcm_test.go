package fcm

import (
	"context"
	"testing"
)

func TestSendMessage(t *testing.T) {
	f, err := NewFcm("path_to_key_file.json")
	if err != nil {
		t.Fail()
	}

	data := map[string]string {
		"msg": "perkakas test send message to fcm",
	}

	_, err = f.Send(context.Background(), "topic_test", data)
	if err != nil {
		t.Fail()
	}
}