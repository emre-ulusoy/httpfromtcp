package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/emre-ulusoy/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCode200 StatusCode = iota
	StatusCode400
	StatusCode500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusCode200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK"))
		if err != nil {
			return err
		}

	case StatusCode400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request"))
		if err != nil {
			return err
		}

	case StatusCode500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
		if err != nil {
			return err
		}

	default:
		_, err := w.Write([]byte(""))
		if err != nil {
			return err
		}
	}

	fmt.Println(statusCode) // FIX: remove
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		fmt.Println("write err")
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders["Content-Length"] = strconv.Itoa(contentLen)
	defaultHeaders["Connection"] = "close"
	defaultHeaders["Content-Type"] = "text/plain"

	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
