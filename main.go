package main

import (
	"fmt"
	"log"
	"nosepass/storage"
	"os"
	"strings"
)

func help() {
	fmt.Printf("Nosepass. Simple yet another password manager\n" +
		"Developed by Dmitry K (ex0hunt) https://github.com/ex0hunt/nosepass" +
		"\n\nnosepass insert mail/gmail.com\nnosepass get mail/gmail\nnosepass delete mail/gmail\n")
}

func _getNodeInLevel(name string, nodes []PathTreeNode) PathTreeNode {
	var res PathTreeNode
	for _, b := range nodes {
		if b.name == name {
			res = b
			break
		}
	}
	return res
}

type PathTreeNode struct {
	level     int
	name      string
	childrens []PathTreeNode
}

func passTree(listPass []string) {
	levels := make(map[int][]PathTreeNode)
	for _, path := range listPass {
		nodeList := strings.Split(path, "/")
		var parentItem PathTreeNode
		for i, node := range nodeList {
			parentItem = _getNodeInLevel(node, levels[i])
		}
	}
}

func main() {
	// Main func with args
	var appError error

	switch cmdType := os.Args[1]; cmdType {
	case "insert":
		if len(os.Args) < 2 || os.Args[1] == "help" {
			help()
			return
		}
		passName := os.Args[2]
		appError = storage.StorePassword(passName)
	case "get":
		if len(os.Args) < 2 || os.Args[1] == "help" {
			help()
			return
		}
		passName := os.Args[2]
		pwd, err := storage.GetPassword(passName)
		appError = err
		fmt.Printf("\n%s\n", pwd)
	case "delete":
		if len(os.Args) < 2 || os.Args[1] == "help" {
			help()
			return
		}
		passName := os.Args[2]
		err := storage.DeletePassword(passName)
		appError = err
	case "show":
		listPath, err := storage.ListPassword()
		if err == nil {
			passTree(listPath)
		}
		appError = err

	default:
		help()
	}

	if appError != nil {
		log.Panic(appError)
	}
}
