package health

import (
	"context"
	"fmt"
	"log"
	"net"
	"system-design/loadbalancer/config"
	"system-design/loadbalancer/config/models"
	"time"
)

type Checker interface {
	HealthCheck(ctx context.Context)
}

type healthChecker struct {
	configProcessor config.Processor
	configChanged   <-chan struct{}
}

func (hc *healthChecker) HealthCheck(ctx context.Context) {
	targetTime := hc.updateConfig()
	for {
		select {
		case <-time.After(time.Until(targetTime)):
			targetTime = hc.updateConfig()
		case <-hc.configChanged:
			targetTime = hc.updateConfig()
		case <-ctx.Done():
			log.Print(ctx.Err())
			break
		}
	}
}

func checkHealthOfServers(hosts []string, prevUnHealthy []string) ([]string, []string) {
	healthy := make([]string, 0)
	unhealthy := make([]string, 0)
	for _, host := range hosts {
		if checkHost(host) {
			healthy = append(healthy, host)
		} else {
			unhealthy = append(unhealthy, host)
		}
	}

	for _, host := range prevUnHealthy {
		if checkHost(host) {
			healthy = append(healthy, host)
		} else {
			unhealthy = append(unhealthy, host)
		}
	}
	return healthy, unhealthy
}

func checkHost(host string) bool {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		fmt.Printf("Server %s is down or unreachable: %v\n", host, err)
		return false
	}
	defer conn.Close()
	return true
}

func (hc *healthChecker) updateConfig() time.Time {
	configuration := hc.configProcessor.Read()
	healthy, unhealthy := checkHealthOfServers(configuration.CurrentAPIServers, configuration.UnhealthyServers)
	if !hc.configProcessor.Equal(configuration, models.Configuration{CurrentAPIServers: healthy, UnhealthyServers: unhealthy}) {
		configuration.UnhealthyServers = unhealthy
		configuration.CurrentAPIServers = healthy
		select {
		case <-hc.configChanged:
			hc.updateConfig()
		default:
			hc.configProcessor.Write(configuration)
		}
	}
	return time.Now().Add(5 * time.Minute)
}

func NewHealthChecker(processor config.Processor, configChanged <-chan struct{}) Checker {
	return &healthChecker{

		configProcessor: processor,
		configChanged:   configChanged,
	}
}
