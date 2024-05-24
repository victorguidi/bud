package main

import (
	"log"
	"net"
	"os"
	"strings"
)

type BudCLI string

// Constants for common HTTP methods
const (
	ASK     BudCLI = "ask"
	ASKBASE BudCLI = "askbase"
	DIR     BudCLI = "dir"
	QUIT    BudCLI = "quit"
)

func (m BudCLI) String() string {
	return string(m)
}

func ParseCommand(args []string) error {
	switch args[1] {
	case DIR.String():
		SendCommand(strings.Join(args[1:], " "))
	case ASKBASE.String():
		SendCommand(strings.Join(args[1:], " "))
	case ASK.String():
		SendCommand(strings.Join(args[1:], " "))
	}
	return nil
}

func SendCommand(cmd string) {
	conn, err := net.Dial("tcp", "0.0.0.0:9876")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte(cmd))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	log.Println(string(buf[:n]))
}

func main() {
	args := os.Args
	ParseCommand(args)
}
