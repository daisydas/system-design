package main

import "fmt"

type (
	TransportType string
	ServerOptions func(Options) Options
)

const (
	UDP TransportType = "udp"
	TCP TransportType = "tcp"
)

type Options struct {
	MaxConnections int
	TransportType  TransportType
	Name           string
}

type Server struct {
	Options
}

func SetMaxConnections(n int) ServerOptions {
	return func(o Options) Options {
		o.MaxConnections = n
		return o
	}
}
func SetTransportType(transportType TransportType) ServerOptions {
	return func(o Options) Options {
		o.TransportType = transportType
		return o
	}
}

func SetName(name string) ServerOptions {
	return func(o Options) Options {
		o.Name = name
		return o
	}
}

func NewServer(opts ...ServerOptions) *Server {
	options := Options{}
	for _, o := range opts {
		options = o(options)
	}
	return &Server{
		Options: options,
	}
}

func main() {
	s := NewServer(SetMaxConnections(10), SetName("test-server"), SetTransportType(UDP))
	fmt.Println(s)
}
