package config

import (
	"net"
	"net/http"
	"sync"
	"system-design/loadbalancer/config/models"
	"time"
)

const (
	Create = "create"
	Update = "update"
)

type Processor interface {
	Read() models.Configuration
	Write(models.Configuration, string)
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

	newMap := make(map[string]*http.Client)
	for k, v := range p.clientMap {
		newMap[k] = v
	}
	return newMap
}

func (p *configProcessor) Write(cfg models.Configuration, action string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if action == Create {
		p.config = cfg
	} else if action == Update {
		p.config.CurrentAPIServers = append(p.config.CurrentAPIServers, cfg.CurrentAPIServers...)
		p.config.Algorithm = cfg.Algorithm
	}

	for i := 0; i < len(cfg.CurrentAPIServers); i++ {
		client, _ := p.clientMap[cfg.CurrentAPIServers[i]]
		if client == nil {
			client = &http.Client{
				Transport: getTransport(),
			}
		}
		p.clientMap[cfg.CurrentAPIServers[i]] = client
	}

	for k, _ := range p.clientMap {
		if !isActiveServer(cfg.CurrentAPIServers, k) {
			delete(p.clientMap, k)
		}
	}

	return
}

func isActiveServer(currentServers []string, mapHost string) bool {
	for _, server := range currentServers {
		if server == mapHost {
			return true
		}
	}
	return false
}

func getTransport() http.RoundTripper {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		MaxIdleConns:    5,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 5 * time.Second,
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

func NewProcessor(config models.Configuration, mutex *sync.RWMutex, clientMap map[string]*http.Client) Processor {
	return &configProcessor{
		config:    config,
		mutex:     mutex,
		clientMap: clientMap,
	}
}
