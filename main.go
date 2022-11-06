package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strconv"
)

// Packet generationW

type packet struct {
	data_length []byte
	data        []byte
	hash        []byte
}

func create_size(id int) []byte {
	arr := make([]byte, 2)

	k := 1
	for n := id % 10; id > 0; {
		arr[k] = byte(n)
		k--
		id = id / 10
		n = id % 10
	}

	return arr
}

func parse_packet(p []byte) (int, string, []byte) {
	var length int
	var msg string

	length = int(p[0])*10 + int(p[1])

	for i := 2; i < length+2; i++ {
		msg += string(p[i])
	}

	return length, msg, p[length+2:]
}

func init_packet(msg []byte) packet {
	length := create_size(len(msg))
	p := packet{data_length: length, data: msg, hash: generate_checksum(msg)}
	return p
}

func generate_checksum(msg []byte) []byte {
	hash := sha256.New()
	hash.Write(msg)
	h := hash.Sum(nil)
	return h
}

func send_packet(conn net.Conn, p packet) {
	buffer := append(p.data_length, p.data...)
	buffer = append(buffer, p.hash...)
	conn.Write(buffer)
}

func compare_checksum(ch1 []byte, ch2 []byte) bool {
	for i := 0; i < len(ch1); i++ {
		if ch1[i] != ch2[i] {
			return false
		}
	}

	return true
}

// Client-side

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

		go handle_server(conn)

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("You: ")
			message, _ := reader.ReadString('\n')
			if len(message) > 99 {
				fmt.Println("Message is too long!")
				continue
			}

			p := init_packet([]byte(message))
			send_packet(conn, p)
		}
	}
}

func handle_server(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Server disconnected!")
			break
		}

		_, msg, hash := parse_packet(buffer[:len])

		if compare_checksum(generate_checksum([]byte(msg)), hash) {
			fmt.Printf("\nThey: %s", msg)
		}
	}
}

// Server-side

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

		fmt.Printf("Listening on: %s port...\n", port)
		conn, err := net.Listen("tcp", ":"+port)
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

			fmt.Printf("Client %s connected!\n", server.RemoteAddr().String())
			go handle_client(server)
			defer server.Close()

			for {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("You: ")
				message, _ := reader.ReadString('\n')
				if len(message) > 99 {
					fmt.Println("Message is too long!")
					continue
				}

				p := init_packet([]byte(message))
				send_packet(server, p)
			}
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

		_, msg, hash := parse_packet(buffer[:len])

		if compare_checksum(generate_checksum([]byte(msg)), hash) {
			fmt.Printf("\nThey: %s", msg)
		}
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
