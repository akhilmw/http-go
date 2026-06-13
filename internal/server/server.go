package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/akhilmw/http-go/internal/request"
	"github.com/akhilmw/http-go/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func writeHandlerError(w io.Writer, hErr *HandlerError) {
	body := []byte(hErr.Message)

	response.WriteStatusLine(w, hErr.StatusCode)
	headers := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(w, headers)
	w.Write(body)
}

func Serve(port int, handler Handler) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to bind to port: %v", err)
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}

			fmt.Println("error accepting connection:", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w := response.NewWriter(conn)
		body := []byte(err.Error())

		w.WriteStatusLine(response.StatusBadRequest)

		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/plain")

		w.WriteHeaders(headers)
		w.WriteBody(body)
		return
	}

	w := response.NewWriter(conn)
	s.handler(w, req)

}
