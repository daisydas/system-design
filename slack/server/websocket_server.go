package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"system-design/slack/models"
	"time"
)

type Client struct {
	conn *websocket.Conn
	id   string
}

type Server struct {
	clients          map[string]*Client
	registerClient   chan Client
	unregisterClient chan Client
	broadcast        chan *models.Message
	dm               chan *models.Message
	mu               sync.RWMutex
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.registerClient:
			fmt.Println("inside 1")
			s.mu.Lock()
			s.clients[client.id] = &client
			s.mu.Unlock()

		case client := <-s.unregisterClient:
			fmt.Println("inside 2")
			s.mu.Lock()
			delete(s.clients, client.id)
			s.mu.Unlock()

		case message := <-s.broadcast:
			fmt.Println("inside broadcast")
			for _, client := range s.clients {
				client.conn.WriteMessage(websocket.TextMessage, convertToBytes(message))
			}
		case message := <-s.dm:
			fmt.Println("inside dm")
			for _, client := range s.clients {
				if client.id == message.Receiver {
					fmt.Println("sending direct message")
					client.conn.WriteMessage(websocket.TextMessage, convertToBytes(message))
				}
			}
		}
	}
}

func convertToBytes(message *models.Message) []byte {
	bytes, _ := jsoniter.Marshal(message)
	return bytes
}

func NewServer() *Server {
	return &Server{
		broadcast:        make(chan *models.Message),
		dm:               make(chan *models.Message),
		clients:          make(map[string]*Client),
		registerClient:   make(chan Client),
		unregisterClient: make(chan Client),
	}
}

func handleWebSocket(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		clientID = r.RemoteAddr + "-" + strconv.Itoa(int(time.Now().Unix()))
		fmt.Println(clientID)
	}

	server.registerClient <- Client{conn: conn, id: clientID}

}

func handleMessages(server *Server, w http.ResponseWriter, r *http.Request) {
	var wr models.WebsocketRequest
	bytes, _ := io.ReadAll(r.Body)
	jsoniter.Unmarshal(bytes, &wr)

	if wr.Type == models.BroadCast {
		server.broadcast <- &models.Message{
			Sender:    r.RemoteAddr,
			Receiver:  getClientID(server.clients, wr.Client),
			Content:   wr.Message,
			Timestamp: 0,
		}
	} else if wr.Type == models.DirectMessage {
		server.dm <- &models.Message{
			Sender:    r.RemoteAddr,
			Receiver:  getClientID(server.clients, wr.Client),
			Content:   wr.Message,
			Timestamp: 0,
		}
	}
}

func getClientID(clients map[string]*Client, client string) string {
	var count = 0
	for k, _ := range clients {
		clientNumber, _ := strconv.Atoi(client)
		if count == clientNumber {
			fmt.Println("identified client")
			fmt.Println(k)
			return k
		}
		count++
	}
	return ""
}

func main() {

	server := NewServer()
	go server.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(server, w, r)
	})
	srv := &http.Server{
		Addr: ":8089",
	}

	http.HandleFunc("/pd", func(w http.ResponseWriter, r *http.Request) {
		handleMessages(server, w, r)
	})

	srv.ListenAndServe()
}
