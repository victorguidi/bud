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
	ENABLE   BudCLI = "start"
	DISABLED BudCLI = "stop"
	CHAT     BudCLI = "chat"
	RAG      BudCLI = "rag"
	LISTEN   BudCLI = "listen"
	KILL     BudCLI = "kill"
)

func (m BudCLI) String() string {
	return string(m)
}

func ParseCommand(args []string) error {
	switch args[1] {
	case ENABLE.String():
		SendCommand(strings.Join(args[1:], " "))
	case DISABLED.String():
		SendCommand(strings.Join(args[1:], " "))
	case RAG.String():
		SendCommand(strings.Join(args[1:], " "))
	case CHAT.String():
		SendCommand(strings.Join(args[1:], " "))
	case LISTEN.String():
		SendCommand(strings.Join(args[1:], " "))
	case KILL.String():
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
