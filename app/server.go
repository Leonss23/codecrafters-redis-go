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

type RedisDB map[string]string

func run() {
	db := make(RedisDB)
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

		go handleFunction(conn, db)
	}
}

func handleFunction(conn net.Conn, rdb RedisDB) {
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
		} else {
			response = handleCommand(args, rdb)
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed to write to connection.\nError:", err)
			break
		}
	}
}

func handleCommand(args []string, db RedisDB) string {
	command := args[0]
	argsLen := len(args)

	switch strings.ToLower(command) {
	default:
		return makeResponse("+Invalid or unsupported Redis command.")
	case "ping":
		return makeResponse("+PONG")
	case "echo":
		if argsLen < 2 {
			return makeResponse("+ECHO command requires an argument")
		}
		return makeResponse(fmt.Sprintf("+%s", args[1]))
	case "set":
		if argsLen < 3 {
			return makeResponse("+SET command requires KEY and VALUE arguments")
		}
		key, value := args[1], args[2]
		db[key] = value
		return makeResponse("+OK")
	case "get":
		key := args[1]
		value, exists := db[key]
		if exists {
			length := fmt.Sprintf("$%v", len(value))
			return makeResponse(length, value)
		} else {
			return makeResponse("$-1")
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

func makeResponse(s ...string) string {
	var sb strings.Builder
	for _, str := range s {
		sb.WriteString(str)
		sb.WriteString("\r\n")
	}
	response := sb.String()
	fmt.Println(response)
	return response
}
