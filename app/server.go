package main

import (
	"fmt"
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		go func() {
			for {
				buffer := make([]byte, 1024)
				if n, err := conn.Read(buffer); err != nil {
					if n == 0 {
						fmt.Println("Connection closed")
						return
					}

					fmt.Println("Error reading message: ", err.Error())
					os.Exit(1)
				}

				conn.Write(prepareRESP("PONG"))
			}
		}()
	}
}

func prepareRESP(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}
