package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	fmt.Println("><><>< ", idx)

	// Incomplete data
	if idx < 0 {
		return 0, false, nil
	}
	// DONE
	if idx == 0 {
		return 0, true, nil
	}

	colonIdx := bytes.Index(data[:idx], []byte(":"))
	key := data[:colonIdx]
	val := data[colonIdx+1 : idx]

	// Check if there's space b/w the key and the colon
	if key2 := bytes.TrimSuffix(key, []byte(" ")); len(key2) < len(key) {
		return 0, false, errors.New("Invalid spacing in header")
	}
	// Trim whitespace
	key = bytes.TrimSpace(key)
	val = bytes.TrimSpace(val)
	key = bytes.ToLower(key)
	bytes.ContainsAny(key, chars)
	strings.ContainsAny

	h[string(key)] = string(val)
	return idx + 2, false, nil
}
