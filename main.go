package main

import (
	"fmt"

	"github.com/theMariusK/chroom/src"
)

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
			src.StartClient()
		} else {
			src.StartServer()
		}
	}
}
