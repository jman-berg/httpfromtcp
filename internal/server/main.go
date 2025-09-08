package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jman-berg/httpfromtcp/internal/request"
	"github.com/jman-berg/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode   int
	ErrorMessage string
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, response.StatusCode(he.StatusCode))
	messageBytes := []byte(he.ErrorMessage)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

type Server struct {
	closed   atomic.Bool
	Listener net.Listener
	Handler  Handler
}

func Serve(port int, h Handler) (*Server, error) {
	portString := strconv.Itoa(port)

	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Listener: l,
		Handler:  h,
	}
	go server.listen()
	return server, nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	r, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode:   int(response.StatusBadRequest),
			ErrorMessage: err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	if hErr := s.Handler(buf, r); hErr != nil {
		hErr.Write(conn)
		return
	}

	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusOK)
	h := response.GetDefaultHeaders(len(b))
	if err := response.WriteHeaders(conn, h); err != nil {
		log.Printf("Error writing errors: %v", err)
	}
	conn.Write(b)
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
