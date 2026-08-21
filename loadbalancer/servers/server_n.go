package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerN struct {
	port     string
	maxRange int
}

func (s1 *ServerN) Start() {

	mh := NewMetricsHandler()
	mh.Get(s1.port)
	
	ge := gin.Default()
	ge.Handle("GET", "/api/fibonacci", func(c *gin.Context) {
		fmt.Printf("on server: %v\n", s1.port)
		out := getFibonacciSeries(s1.maxRange)
		c.JSON(200, gin.H{
			"whichServer": s1.port,
			"fibonacci":   out,
		})
	})

	ge.Handle("GET", "/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"port":   s1.port,
		})
	})

	ge.Handle("GET", "/metrics", gin.WrapH(promhttp.Handler()))

	_ = ge.Run(fmt.Sprintf(":%s", s1.port))
}

func getFibonacciSeries(max int) []int {
	out := make([]int, 0)
	prev := 0
	cur := 1
	out = append(out, 0)
	for i := 1; i < max; i++ {
		x := prev + cur
		out = append(out, x)
		prev = cur
		cur = x
	}
	return out
}

func NewServerN(port string, maxRange int) *ServerN {
	return &ServerN{
		port:     port,
		maxRange: maxRange,
	}
}
