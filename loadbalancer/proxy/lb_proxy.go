package proxy

import (
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"net/http"
	"system-design/loadbalancer/config"
	"time"
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
			go s.forwardRequest(client, req)
			currentServerIndex++
		case <-ctx.Done():
			log.Print(ctx.Err())
			break
		}
	}
}

func (s *sender) forwardRequest(client *http.Client, req *Request) {
	resp, err := client.Do(req.HttpRequest)
	if err != nil {
		conn, connErr := net.DialTimeout("tcp", req.HttpRequest.Host, 10*time.Millisecond)
		if connErr != nil {
			s.forceReadConfig <- struct{}{}
			s.requestChan <- req
			return
		}
		conn.Close()
		http.Error(req.HttpWriter, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		req.HttpWriter.Header()[k] = v
	}
	req.HttpWriter.WriteHeader(resp.StatusCode)
	_, err = io.Copy(req.HttpWriter, resp.Body)
}

func New(configProcessor config.Processor, requestChan chan *Request, forceReadConfig chan struct{}) Sender {
	return &sender{
		configProcessor: configProcessor,
		requestChan:     requestChan,
		forceReadConfig: forceReadConfig,
	}
}
