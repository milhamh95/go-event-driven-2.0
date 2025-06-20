package main

import (
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

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		logger.Error("failed to create publisher", err, nil)
		return
	}

	msg1 := message.NewMessage(watermill.NewUUID(), []byte("50"))
	msg2 := message.NewMessage(watermill.NewUUID(), []byte("100"))

	err = publisher.Publish("progress", msg1)
	if err != nil {
		logger.Error("failed to publish message", err, nil)
		return
	}

	err = publisher.Publish("progress", msg2)
	if err != nil {
		logger.Error("failed to publish message", err, nil)
		return
	}
}
