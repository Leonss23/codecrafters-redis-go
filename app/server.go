package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var connNum = 0

func main() {
	protocol := "tcp"
	ip := "0.0.0.0"
	port := 6379
	host := fmt.Sprint(ip, ":", port)

	listener, err := net.Listen(protocol, host)
	if err != nil {
		fmt.Println("Failed to bind to port", port, "\nError:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}
		connNum += 1

		go HandleFunction(conn)
	}
}

func HandleFunction(conn net.Conn) {
	defer conn.Close()
	var buf [256]byte

	fmt.Printf("Processing connection #%v\n", connNum)
	time.Sleep(time.Second)

	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Failed to read data from connection.\nError:", err)
			break
		}
		received := buf[:n]
		fmt.Printf("Received: %v bytes\n```\n%v```\n", n, string(received))

		response := []byte("+PONG\r\n")
		n, err = conn.Write(response)
		if err != nil {
			fmt.Println("Failed to write to connection.\nError:", err)
			break
		}
		fmt.Printf("Sent: %v bytes\n```\n%v```\n", n, string(response))
	}
}
