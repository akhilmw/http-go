package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

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
		writeHandlerError(conn, &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Bad Request\n",
		})
		return
	}

	var buf bytes.Buffer

	handlerErr := s.handler(&buf, req)
	if handlerErr != nil {
		writeHandlerError(conn, handlerErr)
		return
	}

	body := buf.Bytes()

	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(conn, headers)
	conn.Write(body)
}
