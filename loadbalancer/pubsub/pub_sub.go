package pubsub

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

type RedisPubSubMessage struct {
	Message string `json:"message"`
}

type RedisPubsub interface {
	Publish(ctx context.Context, channel string, data RedisPubSubMessage)
	Subscribe(ctx context.Context, channel string, handler func(chan struct{}, RedisPubSubMessage)) error
}

type pubSub struct {
	client        *redis.Client
	configUpdates chan struct{}
}

func (ps *pubSub) Publish(ctx context.Context, channel string, data RedisPubSubMessage) {
	ps.client.Publish(ctx, channel, data.Message)
	return
}

func (ps *pubSub) Subscribe(ctx context.Context, channel string, handler func(configUpdates chan struct{}, msg RedisPubSubMessage)) error {
	receive := ps.client.Subscribe(ctx, channel)
	msgChan := receive.Channel()

	for {
		select {
		case msg := <-msgChan:
			bytes := []byte(msg.Payload)
			x := RedisPubSubMessage{}
			_ = jsoniter.Unmarshal(bytes, &x)
			handler(ps.configUpdates, x)
		case <-ctx.Done():
			return ctx.Err()

		}
	}
}

func NewPubSub(configUpdates chan struct{}) RedisPubsub {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DB:           0,
		MaxRetries:   0,
		DialTimeout:  10,
		ReadTimeout:  10,
		WriteTimeout: 10,
	})
	return &pubSub{
		client:        redisClient,
		configUpdates: configUpdates,
	}
}
