package config

import (
	"net/http"
	"sync"
	"system-design/loadbalancer/config/models"
)

type Processor interface {
	Read() models.Configuration
	Write(models.Configuration)
	Equal(models.Configuration, models.Configuration) bool
	GetClients() map[string]*http.Client
}

type configProcessor struct {
	mutex     *sync.RWMutex
	config    models.Configuration
	clientMap map[string]*http.Client
}

func (p *configProcessor) Read() models.Configuration {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.config
}

func (p *configProcessor) GetClients() map[string]*http.Client {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.clientMap
}

func (p *configProcessor) Write(cfg models.Configuration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.config = cfg
	updateClients(p.clientMap, cfg)
	return
}

func updateClients(clientMap map[string]*http.Client, cfg models.Configuration) {
	for i := 0; i < len(cfg.CurrentAPIServers); i++ {
		if _, ok := clientMap[cfg.CurrentAPIServers[i]]; !ok {
			continue
		}
		client := &http.Client{
			Transport: getTransport(),
			Timeout:   10,
		}
		clientMap[cfg.CurrentAPIServers[i]] = client
	}
}

func getTransport() http.RoundTripper {
	return &http.Transport{
		MaxIdleConns:    5,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 5,
	}
}

func (p *configProcessor) Equal(cfg1 models.Configuration, cfg2 models.Configuration) bool {
	if len(cfg1.CurrentAPIServers) != len(cfg2.CurrentAPIServers) || len(cfg1.UnhealthyServers) != len(cfg2.UnhealthyServers) {
		return false
	}
	for i := 0; i < len(cfg1.CurrentAPIServers); i++ {
		if !isPresent(cfg1.CurrentAPIServers[i], cfg2.CurrentAPIServers) {
			return false
		}
	}
	for i := 0; i < len(cfg1.UnhealthyServers); i++ {
		if !isPresent(cfg1.UnhealthyServers[i], cfg2.UnhealthyServers) {
			return false
		}
	}
	return true
}

func isPresent(s string, servers []string) bool {
	for _, server := range servers {
		if server == s {
			return true
		}
	}
	return false
}

func NewProcessor(config models.Configuration, mutex *sync.RWMutex) Processor {
	return &configProcessor{
		config:    config,
		mutex:     mutex,
		clientMap: make(map[string]*http.Client),
	}
}
