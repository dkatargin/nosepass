package store

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	//mypass insert mail/gmail.com
	binName := filepath.Base(os.Args[0])
	cmdType := filepath.Base(os.Args[1])
	passName := filepath.Base(os.Args[2])
	log.Println(binName, cmdType, passName)
}