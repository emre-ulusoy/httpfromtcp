package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type Server struct {
	port     int
	isOpen   atomic.Bool
	listener net.Listener
}

type Listener struct{}

func Serve(port int) (*Server, error) {
	server := &Server{
		port: port,
	}
	server.isOpen.Store(true)

	// Listener stuff
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	server.listener = tcpListener

	go server.listen() // How does this end? When atomic bool is flipped

	return server, nil
}

func (s *Server) Close() error {
	s.listener.Close()
	s.isOpen.Store(false)
	return nil // What error could be returned?
}

func (s *Server) listen() {
	for {
		if !s.isOpen.Load() {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			return // ??
		}
		go s.handle(conn) // Not handled in a goroutine!
		fmt.Println("zzz...")
		time.Sleep(time.Second / 2)
	}

	return // What to do when server is closed
}

func (s *Server) handle(conn net.Conn) {
	response := `HTTP/1.1 200 OK
Content-Type: text/plain

Hello World!`

	conn.Write([]byte(response))

	conn.Close()
}
