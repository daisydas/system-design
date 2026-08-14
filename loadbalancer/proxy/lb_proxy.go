package proxy

import (
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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
	Completed   chan bool
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

			proxyUrl := &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   req.HttpRequest.URL.Path,
			}
			proxyHttpReq, _ := http.NewRequest(req.HttpRequest.Method, proxyUrl.String(), req.HttpRequest.Body)
			go s.forwardRequest(client, proxyHttpReq, req.HttpWriter, req.Completed)
			currentServerIndex++
		case <-ctx.Done():
			log.Print(ctx.Err())
			return
		}
	}
}

func (s *sender) forwardRequest(client *http.Client, req *http.Request, writer http.ResponseWriter, completed chan bool) {
	resp, err := client.Do(req)
	if err != nil {
		conn, connErr := net.DialTimeout("tcp", req.Host, 10*time.Millisecond)
		if connErr != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		conn.Close()
		http.Error(writer, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, val := range v {
			writer.Header().Add(k, val)
		}
	}

	// 2. Write status code SECOND
	writer.WriteHeader(200)

	// 3. Copy body LAST
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
	completed <- true
	return

}

func New(configProcessor config.Processor, requestChan chan *Request, forceReadConfig chan struct{}) Sender {
	return &sender{
		configProcessor: configProcessor,
		requestChan:     requestChan,
		forceReadConfig: forceReadConfig,
	}
}
