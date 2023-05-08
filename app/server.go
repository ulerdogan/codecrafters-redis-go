package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Program has started!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	store := NewStore()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle(conn, store)
	}
}

func handle(conn net.Conn, store *Store) {
	defer conn.Close()

	for {
		value, err := DecodeRESP(bufio.NewReader(conn))
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}

		command := value.command
		args := value.Array()

		switch command {
		case "ping":
			conn.Write(prepareRESPString("PONG"))
		case "echo":
			conn.Write(prepareRESPArray(args))
		case "set":
			store.Set(args[0].String(), args[1].String())
			conn.Write(prepareRESPString("OK"))
		case "get":
			value := store.Get(args[0].String())
			conn.Write(prepareRESPString(value))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
