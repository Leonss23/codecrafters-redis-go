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

	_, err = listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
