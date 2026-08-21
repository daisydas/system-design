package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"net/http"
	"strconv"
	"sync"
	"system-design/loadbalancer/config"
	"system-design/loadbalancer/config/models"
	"system-design/loadbalancer/health"
	"system-design/loadbalancer/proxy"
	"system-design/loadbalancer/redisclient"
	"system-design/loadbalancer/servers"
	"time"
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
	}, &mutex, getClientMap(healthyServers))

	hc := health.NewHealthChecker(cp, configChanged)

	go hc.HealthCheck(ctx)

	reqChan := make(chan *proxy.Request, requestChanBuffer)
	lp := proxy.New(cp, reqChan, configChanged)

	go lp.Send(ctx)

	sh := serveHandler{reqChan: reqChan}
	http.HandleFunc("/api/fibonacci", sh.requestHandler)

	ps := redisclient.NewPubSubClient(
		redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}), cp, configChanged, "api_read_lb_write", "api_write_lb_read")
	go ps.ReadFromPubSub(context.Background())
	fmt.Println("Server starting on :8080...")
	_ = http.ListenAndServe(":8080", nil)

}

func getClientMap(healthyServers []string) map[string]*http.Client {
	clientMap := make(map[string]*http.Client)
	for i := 0; i < len(healthyServers); i++ {
		client := http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				MaxIdleConns:    5,
				MaxConnsPerHost: 10,
				IdleConnTimeout: 5 * time.Second,
			},
		}
		clientMap[healthyServers[i]] = &client
	}
	return clientMap
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
	reqChan  chan *proxy.Request
	respChan chan http.Response
}

func (sh *serveHandler) requestHandler(w http.ResponseWriter, r *http.Request) {
	completed := make(chan bool)
	sh.reqChan <- &proxy.Request{
		Completed:   completed,
		HttpRequest: r,
		HttpWriter:  w,
	}
	select {
	case <-completed:
	}

}
