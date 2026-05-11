package main

import (
	"sync"
)

type ConnectionPool struct {
	mutex       *sync.Mutex
	channel     chan int
	connections []*Conn
}

func (cp *ConnectionPool) Release(conn *Conn) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.connections = append(cp.connections, conn)
	cp.channel <- 0
}

func (cp *ConnectionPool) Close() {
	for _, conn := range cp.connections {
		_ = conn.db.Close()
	}
	close(cp.channel)
}

func (cp *ConnectionPool) Get() *Conn {
	<-cp.channel
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	conn := cp.connections[0]
	cp.connections = cp.connections[1:]

	return conn
}

func NewConnectionPool(maxConnections int) *ConnectionPool {
	channel := make(chan int, maxConnections)
	connections := make([]*Conn, 0, maxConnections)
	for i := 0; i < maxConnections; i++ {
		channel <- 0
		connections = append(connections, newConnection())
	}

	return &ConnectionPool{
		channel:     channel,
		mutex:       &sync.Mutex{},
		connections: connections,
	}
}
