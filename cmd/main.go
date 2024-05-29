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
	HELP     BudCLI = "help"
)

// Map to associate BudCLI values with their aliases
var cliAliases = map[BudCLI][]string{
	ENABLE:   {"start", "-s"},
	DISABLED: {"stop", "-x"},
	CHAT:     {"chat", "-c"},
	RAG:      {"rag", "-r"},
	LISTEN:   {"listen", "-l"},
	KILL:     {"kill", "-k"},
	HELP:     {"help", "-h"},
}

// Method to check if a string matches any of the aliases
func IsMatch(input string) (BudCLI, bool) {
	for key, aliases := range cliAliases {
		for _, alias := range aliases {
			if input == alias {
				return key, true
			}
		}
	}
	return "", false
}

func (m BudCLI) String() string {
	return string(m)
}

func ParseCommand(args []string) error {
	if command, matched := IsMatch(args[1]); matched {
		switch command.String() {
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
		case HELP.String():
			SendCommand(strings.Join(args[1:], " "))
		}
	} else {
		SendCommand(HELP.String())
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
