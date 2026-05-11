package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const maxClients = 5

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(maxClients + 1)

	h := Initialize()
	handleShutDown(h, wg)

	port := 8081
	for i := 0; i < maxClients; i++ {
		go func() {
			oneClient(h, strconv.Itoa(port+i), wg)
		}()
	}

	wg.Wait()
	fmt.Println("done")
}

func handleShutDown(h *Handler, wg *sync.WaitGroup) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Kill command
		syscall.SIGQUIT, // Quit command
		syscall.SIGHUP,  // Hangup (terminal closed)
	)

	// Handle these in a separate goroutine
	go func() {
		sig := <-sigChan
		fmt.Println("received signal:", sig)
		h.DropTable()
		defer wg.Done()
	}()
}

func oneClient(h *Handler, port string, wg *sync.WaitGroup) {
	ge := gin.Default()
	ge.POST("/put", gin.WrapF(h.PutKeyHandler))
	ge.GET("/get", gin.WrapF(h.GetKeyHandler))
	ge.GET("/delete", gin.WrapF(h.DeleteKeyHandler))

	ge.Run(fmt.Sprintf(":%s", port))
	wg.Done()
}
