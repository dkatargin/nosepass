package storage

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
)

func StorePassword(dstPath string) {
	fmt.Print("Input password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("Error write password: " + string(password))
	}
	fmt.Printf("\npassword %s stored to %s\n", password, dstPath)
}
