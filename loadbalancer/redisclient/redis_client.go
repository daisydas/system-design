package redisclient

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"system-design/loadbalancer/config"
	cfgModels "system-design/loadbalancer/config/models"
	"system-design/loadbalancer/models"
)

type PubSub interface {
	ReadFromPubSub(ctx context.Context) error
	WriteToPubSub(ctx context.Context, redisPubSubMessage models.RedisPubSubMessage)
}

func (ps *pubSubClient) ReadFromPubSub(ctx context.Context) error {
	pubSub := ps.redisClient.Subscribe(ctx, ps.readChannelName)
	_, err := pubSub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", ps.readChannelName, err)
	}

	msg := pubSub.Channel()
	for {
		select {
		case redisMsg, ok := <-msg:
			if !ok {
				fmt.Println("redis channel closed")
			}
			fmt.Println("received message in redis client")
			x := models.RedisPubSubMessage{}
			_ = jsoniter.Unmarshal([]byte(redisMsg.Payload), &x)
			cfgMessage := cfgModels.Configuration{CurrentAPIServers: x.CurrentAPIServers}
			ps.cfgProcessor.Write(cfgMessage, config.Update)
			ps.configChanged <- struct{}{}

			fmt.Println(cfgMessage)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ps *pubSubClient) WriteToPubSub(ctx context.Context, redisPubSubMessage models.RedisPubSubMessage) {
	jsonData, _ := jsoniter.Marshal(redisPubSubMessage)
	ps.redisClient.Publish(ctx, ps.writeChannelName, string(jsonData))
}

type pubSubClient struct {
	redisClient      *redis.Client
	cfgProcessor     config.Processor
	writeChannelName string
	readChannelName  string
	configChanged    chan struct{}
}

func NewPubSubClient(redisClient *redis.Client, cfgProcessor config.Processor, configChanged chan struct{}, writeChannelName, readChannelName string) PubSub {
	return &pubSubClient{
		redisClient:      redisClient,
		cfgProcessor:     cfgProcessor,
		configChanged:    configChanged,
		writeChannelName: writeChannelName,
		readChannelName:  readChannelName,
	}
}
