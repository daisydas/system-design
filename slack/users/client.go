package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"system-design/slack/models"
)

type Client struct {
	conn      *websocket.Conn
	id        string
	serverURL string
	done      chan bool
}

func NewClient(serverURL string, clientID string) *Client {
	return &Client{
		id:        clientID,
		serverURL: serverURL,
		done:      make(chan bool),
	}
}

func (c *Client) Connect(wg *sync.WaitGroup) error {
	url := fmt.Sprintf("%s?id=%s", c.serverURL, c.id)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.conn = conn
	c.ReceiveMessages()
	defer wg.Done()
	return nil
}

func (c *Client) Disconnect() {
	c.conn.Close()
	c.done <- true
}

func (c *Client) ReceiveMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			_, msgs, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			c.handleMessage(msgs)
		}
	}
}

func (c *Client) handleMessage(msgs []byte) {
	var msg models.Message
	err := json.Unmarshal(msgs, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg.Content)
}

func main() {
	serverUrl := "ws://localhost:8089/ws"
	wg := &sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {
		useri := NewClient(serverUrl, strconv.Itoa(i))
		go useri.Connect(wg)
	}

	wg.Wait()
}
