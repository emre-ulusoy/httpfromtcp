package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/emre-ulusoy/httpfromtcp/internal/request"
	"github.com/emre-ulusoy/httpfromtcp/internal/response"
	"github.com/emre-ulusoy/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(handlerFunc, port)
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

func handlerFunc(writer io.Writer, req *request.Request) *server.HandlerError {
	handlerErr := &server.HandlerError{}
	targetPath := req.RequestLine.RequestTarget
	if targetPath == "/yourproblem" {
		handlerErr.StatusCode = response.StatusCode400
		handlerErr.Message = "Your problem is not my problem\n"
		return handlerErr

	} else if targetPath == "/myproblem" {
		handlerErr.StatusCode = response.StatusCode500
		handlerErr.Message = "Woopsie, my bad\n"
		return handlerErr

	} else {
		writer.Write([]byte("All good, frfr\n"))
		return nil
	}
}
