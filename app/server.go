package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	HandleRequest(conn)
}

func HandleRequest(conn net.Conn) {
	var buf [256]byte

	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Failed to read data from connection.\nError:", err)
			break
		}
		dataIn := buf[:n]
		fmt.Printf("Connection data [%v]:\n%v\n", n, string(dataIn))

		response := []byte("+PONG\r\n")
		n, err = conn.Write(response)
		if err != nil {
			fmt.Println("Failed to write to connection.\nError:", err)
			break
		}
		fmt.Println("Wrote", n, "bytes to connection")
	}
}
