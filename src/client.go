package src

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/theMariusK/chroom/packet"
)

func StartClient() {
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

		go handleServer(conn)

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("You: ")
			message, _ := reader.ReadString('\n')
			if len(message) > 99 {
				fmt.Println("Message is too long!")
				continue
			}

			p := packet.InitPacket([]byte(message))
			packet.SendPacket(conn, p)
		}
	}
}

func handleServer(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("\nSYSTEM: Server disconnected!\n")
			break
		}

		_, msg, hash := packet.ParsePacket(buffer[:len])

		if packet.CompareChecksum(packet.GenerateChecksum([]byte(msg)), hash) {
			fmt.Printf("\nThey: %s", msg)
		}
	}
}
