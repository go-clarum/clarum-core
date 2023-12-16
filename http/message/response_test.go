package message

import (
	"github.com/goclarum/clarum/http/constants"
	"maps"
	"testing"
)

func TestBuilder(t *testing.T) {
	actual := Response(200).
		ContentType("text/plain").
		ETag("5555").
		Payload("batman!")

	expected := ResponseMessage{
		StatusCode: 200,
		Message: Message{
			MessagePayload: "batman!",
			Headers: map[string]string{
				constants.ContentTypeHeaderName: "text/plain",
				constants.ETagHeaderName:        "5555"},
		},
	}

	if !responseEqual(actual, &expected) {
		t.Errorf("Message is not as expected.")
	}
}

func TestClone(t *testing.T) {
	message := Response(500).
		ContentType("text/plain").
		ETag("5555").
		Payload("my payload")

	clonedMessage := message.Clone()

	if clonedMessage == message {
		t.Errorf("Message has not been cloned.")
	}

	if !responseEqual(clonedMessage, message) {
		t.Errorf("Messages are not equal.")
	}
}

func responseEqual(m1 *ResponseMessage, m2 *ResponseMessage) bool {
	if m1.StatusCode != m2.StatusCode {
		return false
	} else if !maps.Equal(m1.Headers, m2.Headers) {
		return false
	} else if m1.MessagePayload != m2.MessagePayload {
		return false
	}
	return true
}
