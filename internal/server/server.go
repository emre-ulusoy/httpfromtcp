package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/emre-ulusoy/httpfromtcp/internal/request"
	"github.com/emre-ulusoy/httpfromtcp/internal/response"
)

// Server is an HTTP 1.1 server
type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // where do we get port now?
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

// WARN: THE ERRORS WILL GET WRITTEN TO THE CONNECTION
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handler called")

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCode400,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	fmt.Println("about to write status line, headers, and body")
	b := buf.Bytes()
	err = response.WriteStatusLine(conn, response.StatusCode200)
	if err != nil {
		fmt.Println("status line err: ", err)
	}

	headers := response.GetDefaultHeaders(len(b))
	fmt.Printf("default headers: %+v", headers)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Println("headers err: ", err)
	}

	n, err := conn.Write(b)
	fmt.Printf("%d bytes written for body", n)
	if err != nil {
		fmt.Println("error writing the body", err)
	}

	fmt.Println("handler done")
	return
}

func (he HandlerError) Write(w io.Writer) {
	fmt.Println("about to write status line, headers, and body")
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		fmt.Println("status line err: ", err)
	}

	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	fmt.Printf("default headers: %+v", headers)
	err = response.WriteHeaders(w, headers)
	if err != nil {
		fmt.Println("headers err: ", err)
	}

	n, err := w.Write(messageBytes)
	fmt.Printf("%d bytes written for body", n)
	if err != nil {
		fmt.Println("error writing the body", err)
	}
}
