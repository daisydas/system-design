package main

import (
	"errors"
	"fmt"
	"sync"
)

const maxConnectionPoolSize = 5

var (
	poolClosedErr = errors.New("connection pool is closed")
)

type ConnectionPool interface {
	Acquire() (Connection, error)
	Release(conn Connection)
	Close()
}

type connectionPool struct {
	lock        *sync.Mutex
	connections chan Connection
	closed      bool
}

type Connection struct {
	Display string
}

func NewConnectionPool() ConnectionPool {

	connections := make(chan Connection, maxConnectionPoolSize)
	for i := 0; i < maxConnectionPoolSize; i++ {
		connections <- Connection{
			Display: fmt.Sprintf("Connection %d", i),
		}
	}

	return &connectionPool{
		connections: connections,
		lock:        &sync.Mutex{},
	}
}

func (p *connectionPool) Acquire() (Connection, error) {
	r, ok := <-p.connections
	if !ok {
		return Connection{}, poolClosedErr
	}
	return r, nil
}

func (p *connectionPool) Release(conn Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		return
	}
	p.connections <- conn
}

func (p *connectionPool) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.closed = true
	close(p.connections)
}
