package main

import (
	"context"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"net/http"
	"system-design/loadbalancer/models"
	"time"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	server := NewApiServer(redisClient, "api_write_lb_read", "api_read_lb_write")
	go func() {
		_ = server.subscribe(context.Background())
	}()
	http.HandleFunc("/read", server.ReadConfigHandler)
	http.HandleFunc("/write", server.WriteConfigHandler)
	fmt.Println("server started")
	_ = http.ListenAndServe(":8054", nil)
}

type ApiServer interface {
	ReadConfigHandler(w http.ResponseWriter, r *http.Request)
	WriteConfigHandler(w http.ResponseWriter, r *http.Request)
	subscribe(ctx context.Context) error
}

func (as *apiServer) ReadConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	bytes, _ := jsoniter.Marshal(as.message)
	_, _ = w.Write(bytes)
}

func (as *apiServer) subscribe(ctx context.Context) error {
	msg := as.redisClient.Subscribe(ctx, as.readChannelName).Channel()
	for {
		select {
		case redisMsg := <-msg:
			x := models.RedisPubSubMessage{}
			_ = jsoniter.Unmarshal([]byte(redisMsg.Payload), &x)
			as.message = x
			fmt.Println("received message in api server")
			fmt.Println(as.message)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (as *apiServer) WriteConfigHandler(w http.ResponseWriter, r *http.Request) {
	msg := &models.RedisPubSubMessage{}
	_ = json.NewDecoder(r.Body).Decode(msg)

	fmt.Println(as.redisClient.Ping(r.Context()).Err())
	jsonData, _ := jsoniter.Marshal(msg)
	x := as.redisClient.Publish(r.Context(), as.writeChannelName, string(jsonData)).Err()
	fmt.Println(x)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("published successfully!!"))
}

type apiServer struct {
	writeChannelName string
	readChannelName  string
	redisClient      *redis.Client
	message          models.RedisPubSubMessage
}

func NewApiServer(redisClient *redis.Client, writeChannelName string, readChannelName string) ApiServer {
	return &apiServer{
		redisClient:      redisClient,
		writeChannelName: writeChannelName,
		readChannelName:  readChannelName,
	}
}
