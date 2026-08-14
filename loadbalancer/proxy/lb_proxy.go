package proxy

import (
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/http/httputil"
	"system-design/loadbalancer/config"
)

type Sender interface {
	Send(ctx context.Context)
}

type sender struct {
	configProcessor config.Processor
	requestChan     chan *Request
	forceReadConfig chan struct{}
}

type Request struct {
	HttpRequest *http.Request
	HttpWriter  http.ResponseWriter
}

func (s *sender) Send(ctx context.Context) {
	cfg := s.configProcessor.Read()
	currentServerIndex := 0
	for {
		select {
		case <-s.forceReadConfig:
			cfg = s.configProcessor.Read()
			continue
		case req := <-s.requestChan:

			if len(cfg.CurrentAPIServers) == 0 {
				http.Error(req.HttpWriter, "No healthy upstream servers", http.StatusServiceUnavailable)
				continue
			}

			if len(cfg.CurrentAPIServers) <= currentServerIndex {
				currentServerIndex = 0
			}

			host := cfg.CurrentAPIServers[currentServerIndex]
			client := s.configProcessor.GetClients()[host]
			req.HttpRequest.Host = host
			req.HttpRequest.URL.Host = host
			go s.forwardRequest(client.Transport, req)
			currentServerIndex++
		case <-ctx.Done():
			log.Print(ctx.Err())
			break
		}
	}
}

func (s *sender) forwardRequest(transport http.RoundTripper, req *Request) {
	proxy := &httputil.ReverseProxy{
		Transport: transport,
	}
	proxy.ServeHTTP(req.HttpWriter, req.HttpRequest)
}

func New(configProcessor config.Processor, requestChan chan *Request, forceReadConfig chan struct{}) Sender {
	return &sender{
		configProcessor: configProcessor,
		requestChan:     requestChan,
		forceReadConfig: forceReadConfig,
	}
}
