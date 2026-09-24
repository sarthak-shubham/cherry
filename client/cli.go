package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	command := os.Args[1]
	key := os.Args[2]

	switch command {
	case "get":
		getKey(key)

	case "set":
		if len(os.Args) != 4 {
			fmt.Println("usage: go run ./client set <key> <value>")
			return
		}

		setKey(key, os.Args[3])

	case "delete":
		deleteKey(key)

	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  go run ./client get <key>")
	fmt.Println("  go run ./client set <key> <value>")
	fmt.Println("  go run ./client delete <key>")
}
