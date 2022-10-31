package main

import (
	"fmt"
	"os"
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
	fmt.Println("CLIENT")
	return 0
}

func start_server() int {
	fmt.Println("SERVER")
	return 0
}

func main() {
	var ptype int
	fmt.Println("Welcome to chROOM!")
	fmt.Println("Please choose program type:\n[1] - Client\n[2] - Server")

	for {
		fmt.Scanf("%d", &ptype)

		if ptype != 1 && ptype != 2 {
			fmt.Println(ptype)
			send_error(INFO, "Please enter 1 (Client) or 2 (Server)")
			continue
		}

		if ptype == 1 {
			start_client()
		} else {
			start_server()
		}
	}
}
