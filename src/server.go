package src

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/theMariusK/chroom/packet"
)

func StartServer() {
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

		if err != nil {
			fmt.Println("Port is taken!")
			os.Exit(1)
		}

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

			fmt.Printf("\nClient %s connected!\n", server.RemoteAddr().String())
			go handleClient(server)
			defer server.Close()

			for {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("You: ")
				message, _ := reader.ReadString('\n')
				if len(message) > 99 {
					fmt.Println("Message is too long!")
					continue
				}

				p := packet.InitPacket([]byte(message))
				packet.SendPacket(server, p)
			}
		}
	}
}

func handleClient(server net.Conn) {
	for {
		buffer := make([]byte, 1024)
		len, err := server.Read(buffer)

		if err != nil {
			fmt.Println("\nSYSTEM: Client disconnected!\n")
			break
		}

		_, msg, hash := packet.ParsePacket(buffer[:len])

		if packet.CompareChecksum(packet.GenerateChecksum([]byte(msg)), hash) {
			fmt.Printf("\nThey: %s", msg)
		}
	}
}
