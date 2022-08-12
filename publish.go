package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

type message struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

func generateMessage() *message {
	minTemp := float32(20)
	maxTemp := float32(35)
	minHumidity := float32(85)
	maxHumidity := float32(99)
	rand.Seed(time.Now().UnixNano())

	msg := message{
		minTemp + rand.Float32()*(maxTemp-minTemp),
		minHumidity + rand.Float32()*(maxHumidity-minHumidity),
	}

	return &msg
}

func Publish(w io.Writer, projectID, topicID string) error {
	ctx := context.Background()
	msg, _ := json.Marshal(generateMessage())
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	defer client.Close()

	t := client.Topic(topicID)
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
		Attributes: map[string]string{
			"origin":   "home-sensor",
			"username": "iotx",
		},
	})

	id, err := result.Get(ctx)

	if err != nil {
		return fmt.Errorf("get: %v", err)
	}

	fmt.Fprintf(os.Stdout, "[%s] Sent message: %s, message ID: %s\n",
		time.Now().Format("01-02-2006 15:04:05"), string(msg), id)

	return nil
}

func main() {
	projectID := "marifw-data-pipelines"
	topicID := "hello-topic"

	buf := new(bytes.Buffer)

	for {
		if err := Publish(buf, projectID, topicID); err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
