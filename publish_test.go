package pubsub

import (
	"bytes"
	"testing"
)

func TestPublish(t *testing.T) {
	projectID := "marifw-data-pipelines"
	topicID := "hello-topic"

	buf := new(bytes.Buffer)

	if err := Publish(buf, projectID, topicID); err != nil {
		t.Errorf("failed to publish message %v", err)
	}
}
