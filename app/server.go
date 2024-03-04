package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func main() {
	run()
}

const DEFAULT_REDIS_PORT = 6379

func run() {
	port := DEFAULT_REDIS_PORT
	host := fmt.Sprintf(":%v", port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Failed to bind to port", port, "\nError:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on port %v\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		go handleFunction(conn)
	}
}

func handleFunction(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte

	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Failed to read data from connection.\nError:", err)
			break
		}
		received := buf[:n]
		fmt.Printf("Received: %v bytes\n```\n%v```\n", n, string(received))

		var response string

		args, errMsg := parseRedisCommand(string(received))
		if errMsg != "" {
			response = errMsg
		}
		switch strings.ToUpper(args[0]) {
		case "ECHO":
			response = makeResponse(args[1])

		case "PING":
			response = makeResponse("PONG")

		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed to write to connection.\nError:", err)
			break
		}
	}
}

func parseRedisCommand(input string) ([]string, string) {
	lines := strings.Split(strings.TrimSuffix(input, "\r\n"), "\r\n")

	linesLen := len(lines)

	for i := 0; i < linesLen; i++ {
		line := lines[i]
		fmt.Printf("[%v] %s\n", i, line)

	}

	var args []string

	switch lines[0][0] {
	default:
		return nil, "Invalid Redis command, root element should be \"*\"(array)"
	case '*':
		argCount, _ := strconv.Atoi(lines[0][1:])

		// TODO: parse each argument in the array
		for i := 0; i < argCount; i++ {
			args = append(args, lines[2+i*2])
		}
		fmt.Println("Args:", args)
	}
	return args, ""
}

func makeResponse(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}
