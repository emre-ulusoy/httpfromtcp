package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not listen on port 42069: %s\n", err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}

		fmt.Printf("Connection accepted, address is %s\n", connection.RemoteAddr().String())

		strCh := getLinesChannel(connection)
		for string := range strCh {
			fmt.Printf("%v\n", string)
			fmt.Println(string)
		}
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	currentLineContents := ""
	lineChan := make(chan string)

	go func() {
		for {
			buffer := make([]byte, 8, 8)
			n, err := conn.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lineChan <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					close(lineChan)
					conn.Close()
					fmt.Println("connection closed")
					return
				}
				fmt.Printf("error: %s\n", err.Error())
				close(lineChan)
				conn.Close()
				return
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lineChan <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return lineChan
}
