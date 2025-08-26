package httpserver

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
)

const port = 42069

type Server struct {
	Open     atomic.Bool
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
	server.Open.Store(true)
	return server, nil
}

func (s *Server) listen() {
	for s.Open {
		s.Listener.Accept()
	}
}

func (s *Server) close() error {
	s.Open.Store(false)
}

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
