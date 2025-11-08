package main

import (
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
	server, err := server.Serve(port, handlerFunc)
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

func handlerFunc(w *response.Writer, req *request.Request) {
	targetPath := req.RequestLine.RequestTarget
	if targetPath == "/yourproblem" {
		request400(w, req)
	} else if targetPath == "/myproblem" {
		request500(w, req)
	} else {
		request200(w, req)
	}
	return
}

func request400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}

func request500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}

func request200(w *response.Writer, _ *request.Request) {
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}
