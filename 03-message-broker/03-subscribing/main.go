package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewStdLogger(false, false)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		logger.Error("failed to create subscriber", err, nil)
		return
	}

	messages, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		logger.Error("failed to subscribe", err, nil)
		return
	}

	for msg := range messages {
		go func(m *message.Message) {
			orderID := string(m.Payload)
			fmt.Printf("Message ID: %s - %s%%\n", m.UUID, orderID)
			m.Ack()
		}(msg)
	}

}
