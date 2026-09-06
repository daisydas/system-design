package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	r "system-design/locks/redis"
	"time"
)

type Server struct {
	lock   r.Data
	client *http.Client
}

func main() {
	srv := Server{
		lock: r.Data{
			Client: redis.NewClient(&redis.Options{
				Addr:        "localhost:6379",
				DialTimeout: 10 * time.Second,
				ReadTimeout: 10 * time.Second,
			}),
		},
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          2,
				IdleConnTimeout:       5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
			},
		},
	}
	http.Handle("/get", &srv)
	_ = http.ListenAndServe(":8082", nil)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.lock.GetKeyFromRedis(r.Context(), "localhost:8082")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/read", nil)

	resp, err := s.client.Do(request)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		s.lock.RemoveLock(r.Context(), "localhost:8082")
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
	s.lock.RemoveLock(r.Context(), "localhost:8082")

}
