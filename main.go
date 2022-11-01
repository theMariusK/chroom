package main

import (
	"fmt"
	"net"
	"strconv"
)

func start_client() {
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

		fmt.Printf("Connecting to: %s:%s...\n", ip, port)
		conn, err := net.Dial("tcp", ip+":"+port)
		defer conn.Close()

		if err != nil {
			fmt.Println("There was an error while connecting to server!")
			continue
		}
		fmt.Println("Connection successful! Write something and press `Enter`")

		for {
			var message string
			fmt.Print("You: ")
			fmt.Scanf("%s\n", &message)
			conn.Write([]byte(message))
		}
	}
}

func start_server() {
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

		ip := "127.0.0.1"
		fmt.Printf("Listening on: %s:%s...\n", ip, port)
		conn, err := net.Listen("tcp", ip+":"+port)
		defer conn.Close()

		if err != nil {
			fmt.Println("There was an error while opening connection!")
			continue
		}

		for {
			server, err := conn.Accept()

			if err != nil {
				fmt.Println("Error occured while trying to handle the connection!")
			}

			fmt.Printf("Client %s connected!\n", conn.Addr().String())
			go handle_client(server)
			defer server.Close()
		}
	}
}

func handle_client(server net.Conn) {
	for {
		buffer := make([]byte, 1024)
		len, err := server.Read(buffer)

		if err != nil {
			fmt.Println("Client disconnected!")
			break
		}

		fmt.Printf("They: %s\n", buffer[:len])
	}
}

func main() {
	var ptype int
	fmt.Println("Welcome to chROOM!")
	fmt.Println("Please choose program type:\n[1] - Client\n[2] - Server")

	for {
		fmt.Scanf("%d\n", &ptype)

		if ptype != 1 && ptype != 2 {
			fmt.Println("Please enter 1 (Client) or 2 (Server)")
			continue
		}

		if ptype == 1 {
			start_client()
		} else {
			start_server()
		}
	}
}
