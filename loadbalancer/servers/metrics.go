package servers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"os"
	"time"
)

type MetricsHandler interface {
	Get(port string)
}

type metricsHandler struct {
}

func (mh *metricsHandler) Get(port string) {

	cpu, mem, health := initialize(port)

	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		log.Printf("Failed to get process: %v", err)
	}
	go updateMetrics(port, proc, cpu, mem, health)

}

func initialize(port string) (*prometheus.GaugeVec, *prometheus.GaugeVec, *prometheus.GaugeVec) {

	cpuGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_cpu_percent_" + port,
			Help: "CPU usage percentage for this server",
		},
		[]string{port},
	)

	memoryGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_memory_bytes_" + port,
			Help: "Memory usage in bytes for this server",
		},
		[]string{port},
	)

	healthGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_health_" + port,
			Help: "Health status of this server (1=healthy, 0=unhealthy)",
		},
		[]string{port},
	)
	prometheus.MustRegister(cpuGauge, memoryGauge, healthGauge)
	return cpuGauge, memoryGauge, healthGauge
}

func updateMetrics(port string, proc *process.Process, cpu, mem, health *prometheus.GaugeVec) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		memInfo, err := proc.MemoryInfo()
		if err != nil {
			continue
		}

		cpu.WithLabelValues(port).Set(cpuPercent)
		mem.WithLabelValues(port).Set(float64(memInfo.RSS))
		health.WithLabelValues(port).Set(1)
	}
}

func NewMetricsHandler() MetricsHandler {
	return &metricsHandler{}
}
