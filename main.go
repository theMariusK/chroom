package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	INFO  int = 1
	WARN  int = 2
	ERROR int = 3
)

func send_error(signal int, message string) {
	switch signal {
	case INFO:
		fmt.Printf("%s\n", message)
	case WARN:
		fmt.Printf("[WARNING] %s\n", message)
	case ERROR:
		fmt.Printf("[ERROR] %s\n", message)
		os.Exit(signal)
	default:
		fmt.Printf("Something went wrong!")
	}
}

func start_client() int {
	for {
		fmt.Println("Enter server's IP Address:")
		var ip string
		fmt.Scanf("%s\n", &ip)
		if net.ParseIP(ip) == nil {
			fmt.Printf("Entered IP Address: %s is invalid\n", ip)
			continue
		}

		fmt.Println("Enter Port number (default is 7777):")
		var port string
		fmt.Scanf("%s\n", &port)
		if port == "" {
			port = "7777"
		}

		_, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("Entered Port number: %s is invalid\n", port)
			continue
		}

		fmt.Printf("Connecting to: %s:%s...", ip, port)
		conn, err := net.Dial("tcp", ip+":"+port)
		defer conn.Close()

		msg := "Hey"
		conn.Write([]byte(msg))

		reply := make([]byte, 1024)
		conn.Read(reply)

		fmt.Println(string(reply))

		break
	}
	return 0
}

func start_server() int {
	for {
		fmt.Println("Enter Port number to listen (default is 7777):")
		var port string
		fmt.Scanf("%s\n", &port)
		if port == "" {
			port = "7777"
		}

		_, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("Entered Port number: %s is invalid\n", port)
			continue
		}

		fmt.Printf("Listening on: 127.0.0.1:%s...", port)
		conn, err := net.Listen("tcp", "127.0.0.1:"+port)
		defer conn.Close()

		for {
			server, _ := conn.Accept()
			go processConnection(server)
		}

		break
	}
	return 0
}

func processConnection(server net.Conn) {
	buffer := make([]byte, 1024)
	len, _ := server.Read(buffer)
	fmt.Printf("Got message: %s\n", buffer[:len])
	server.Write([]byte("Thank you!"))
	server.Close()
}

func main() {
	var ptype int
	fmt.Println("Welcome to chROOM!")
	fmt.Println("Please choose program type:\n[1] - Client\n[2] - Server")

	for {
		fmt.Scanf("%d\n", &ptype)

		if ptype != 1 && ptype != 2 {
			send_error(INFO, "Please enter 1 (Client) or 2 (Server)")
			continue
		}

		if ptype == 1 {
			start_client()
			break
		} else {
			start_server()
		}
	}
}
