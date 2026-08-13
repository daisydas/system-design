package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerN struct {
	port     string
	maxRange int
}

func (s1 *ServerN) Start() {
	ge := gin.Default()
	ge.Handle("GET", "/api/fibonacci", func(c *gin.Context) {
		out := getFibonacciSeries(s1.maxRange)
		c.Writer.WriteHeader(http.StatusOK)
		c.JSON(200, gin.H{
			"whichServer": s1.port,
			"fibonacci":   out,
		})
	})
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
