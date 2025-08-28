package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	closed   atomic.Bool
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	portString := strconv.Itoa(port)

	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Listener: l,
	}
	go server.listen()
	return server, nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: 13\n\nHello World!\n"
	conn.Write([]byte(response))
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error setting up connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.Listener.Close()
	return err
}
