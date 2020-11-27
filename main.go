package main

import (
	"fmt"
	"nosepass/storage"
	"os"
)

func main() {
	//nosepass insert mail/gmail.com
	if len(os.Args) < 3 || os.Args[1] == "help" {
		fmt.Printf("Nosepass help\n\nnosepass insert mail/gmail.com\n")
		return
	}
	cmdType := os.Args[1]
	passName := os.Args[2]

	if cmdType == "insert" {
		storage.StorePassword(passName)
	}

	fmt.Printf("command type is %s\npassword name is %s\n", cmdType, passName)
}
