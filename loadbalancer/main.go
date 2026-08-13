package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"system-design/loadbalancer/config"
	"system-design/loadbalancer/config/models"
	"system-design/loadbalancer/health"
	"system-design/loadbalancer/proxy"
	"system-design/loadbalancer/pubsub"
	"system-design/loadbalancer/servers"
)

const maxServers = 5
const requestChanBuffer = 50

func main() {

	configChanged := make(chan struct{})
	mutex := sync.RWMutex{}
	ctx, _ := context.WithCancel(context.Background())

	healthyServers := startServers(maxServers)

	cp := config.NewProcessor(models.Configuration{
		CurrentAPIServers: healthyServers,
		Algorithm:         "RR",
	}, &mutex)
	hc := health.NewHealthChecker(cp, configChanged)

	go hc.HealthCheck(ctx)

	reqChan := make(chan *proxy.Request, requestChanBuffer)
	lp := proxy.New(cp, reqChan, configChanged)

	go lp.Send(ctx)

	pb := pubsub.NewPubSub(configChanged)
	pb.Publish(ctx, "pub-sub-channel", pubsub.RedisPubSubMessage{})
	go pb.Subscribe(ctx, "pub-sub-channel", func(c chan struct{}, message pubsub.RedisPubSubMessage) {
		c <- struct{}{}
	})

	sh := serveHandler{reqChan: reqChan}
	http.HandleFunc("/load-balance", sh.requestHandler)

	fmt.Println("Server starting on :8080...")
	_ = http.ListenAndServe(":8080", nil)

}

func startServers(maxServers int) []string {
	srvs := make([]string, 0, maxServers)
	port := 8085
	for i := 0; i < maxServers; i++ {
		srvPort := strconv.Itoa(port + i)
		si := servers.NewServerN(srvPort, 1500)
		srvs = append(srvs, fmt.Sprintf("localhost:%s", srvPort))
		go si.Start()
	}
	return srvs
}

type serveHandler struct {
	reqChan chan *proxy.Request
}

func (sh *serveHandler) requestHandler(w http.ResponseWriter, r *http.Request) {
	sh.reqChan <- &proxy.Request{
		HttpRequest: r,
		HttpWriter:  w,
	}

}
