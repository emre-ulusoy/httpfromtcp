package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("error reading from reader", err)
		return nil, err
	}
	reqBytes.Index()

	reqLine := strings.Split(string(reqBytes), "\r\n")[0]
	reqLineStruct, err := ParseRequestLine(reqLine)
	if err != nil {
		return nil, err
	}
	reqStruct := Request{
		RequestLine: *reqLineStruct,
	}
	return &reqStruct, nil
}

func ParseRequestLine(reqLine string) (*RequestLine, error) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	threeParts := strings.Split(reqLine, " ")
	if len(threeParts) < 3 {
		return nil, errors.New("invalid number of parts in the request line")
	}

	ret := RequestLine{
		HttpVersion:   strings.TrimLeft(threeParts[2], "HTTP/"),
		RequestTarget: threeParts[1],
		Method:        threeParts[0],
	}

	if !slices.Contains(methods, ret.Method) {
		return nil, errors.New("invalid method name")
	}

	if ret.HttpVersion != "1.1" {
		return nil, errors.New("invalid HTTP version")
	}

	return &ret, nil
}
