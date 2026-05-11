package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var (
	maxServer = 5
	startPort = 8081
)

const localhost = "localhost"

func main() {
	ge := gin.Default()

	lb := NewLoadBalancer()
	ge.Any("/*path", func(c *gin.Context) {
		path := c.Param("path")
		fmt.Println(path)
		resp := lb.lbHandler(c)
		c.Status(resp.StatusCode)

		c.Writer.Write([]byte{byte(resp.StatusCode)})
	})

	ge.Run(":8093")
}

type loadBalancer struct {
	servers []string
	counter int64
	mutex   *sync.Mutex
	client  *http.Client
}

func (lb *loadBalancer) getServer() int64 {
	return lb.counter % int64(maxServer)
}

func NewLoadBalancer() *loadBalancer {
	servers := make([]string, maxServer)
	for i := 0; i < maxServer; i++ {
		servers[i] = strconv.Itoa(startPort + i)
	}
	return &loadBalancer{servers, 0, &sync.Mutex{}, &http.Client{}}
}

func (lb *loadBalancer) lbHandler(c *gin.Context) *http.Response {
	path := c.Request.URL.Path

	lb.mutex.Lock()
	port := lb.servers[lb.getServer()]
	lb.counter++
	lb.mutex.Unlock()

	url := fmt.Sprintf("http://%s:%s%s", localhost, port, path)
	request, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return &http.Response{StatusCode: http.StatusBadRequest,
			Body: ioutil.NopCloser(bytes.NewBufferString("bad request"))}
	}
	request.Header = c.Request.Header.Clone()

	resp, err := lb.client.Do(request)
	if err != nil {
		return &http.Response{StatusCode: http.StatusInternalServerError,
			Body: ioutil.NopCloser(bytes.NewBufferString("Internal Server Error"))}
	}
	return resp
}
