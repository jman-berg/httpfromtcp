package server

import (
	"github.com/jman-berg/httpfromtcp/internal/response"
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

	response.WriteStatusLine(conn, 200)
	h := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, h); err != nil {
		log.Printf("Error writing errors: %v", err)
	}
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
