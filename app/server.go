package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
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
			conn.Write(prepareRESP("PONG"))
		case "echo":
			conn.Write(prepareRESPArray(args))
		case "set":
			if len(args) < 3 {
				store.Set(args[0].String(), args[1].String(), time.Duration(0))
			} else {
				if args[2].String() == "px" {
					expireIn, err := strconv.Atoi(args[3].String())
					if err != nil {
						conn.Write([]byte("-ERR invalid expire time in SET\r\n"))
						continue
					}

					store.Set(args[0].String(), args[1].String(), time.Duration(expireIn)*time.Millisecond)
				} else {
					conn.Write([]byte(fmt.Sprintf("-ERR unknown option for set: %s\r\n", args[2].String())))
					continue
				}
			}
			conn.Write(prepareRESP("OK"))
		case "get":
			value := store.Get(args[0].String())
			if value != "" {
				fmt.Println("HERE1 ", value)
				conn.Write(prepareRESPString(value))
			} else {
				conn.Write(prepareRESPErr())
			}
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
