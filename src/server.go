package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type ServerProperties struct {
	context.Context
	Host string
	Port string
}

func NewServerEngine(ctx context.Context, host, port string) *ServerProperties {
	return &ServerProperties{
		ctx,
		host,
		port,
	}
}

func (e *ServerProperties) StartServer() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", e.Host, e.Port))
	if err != nil {
		log.Printf("Failed to bind to port %s", e.Port)
		os.Exit(1)
	}
	log.Printf("Engine Socket Listening on: %s", e.Port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go e.HandleConn(conn)
	}
}

func (s *ServerProperties) HandleConn(conn net.Conn) {
	defer conn.Close()

	// This is so we can read everytime the user sends a message
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		buf = bytes.Trim(buf, "\x00") // \x00 is NULL
		buf = bytes.Trim(buf, "\r\n")

		if len(buf) > 0 {
			s.CliParser(string(buf[:n]), conn)
		}
	}
}

func (s *ServerProperties) CliParser(cmd string, conn net.Conn) {
	for k, v := range Workers {
		c := strings.Split(cmd, " ")
		switch c[0] {
		case "help", "-h":
			fmt.Fprint(conn, helpCommand())
			return
		case k:
			v.Call(c[1:])
			if conn != nil {
				fmt.Fprintf(conn, "Worker %s called\n", k)
			}
			return
		case "start", "-s":
			if c[1] == k {
				go v.Run()
				fmt.Fprintf(conn, "STARTED WORKER %s\n", k)
			}
			return
		case "stop", "-S":
			if c[1] == k {
				go v.Stop()
				fmt.Fprintf(conn, "STOPPED WORKER %s\n", k)
			}
			return
		case "kill", "-k":
			if c[1] == k {
				v.Kill()
				delete(Workers, k)
			}
			return
		default:
			fmt.Fprintf(conn, "COULD NOT FIND THE COMMAND, IS THE WORKER %s RUNNING?\n", k)
			return
		}
	}
}

func helpCommand() string {
	var help strings.Builder
	help.WriteString("Usage:\n")
	help.WriteString("  <command> [flags]\n")
	help.WriteString("\n")
	help.WriteString("Available commands:\n")
	help.WriteString("  start -s <worker>     Start a new Worker.\n")
	help.WriteString("  stop -S <worker>      Stop a Worker.\n")
	help.WriteString("  kill -k <worker>      Kill a Worker.\n")
	help.WriteString("  listen -l             Starts a Worker listener.\n")
	help.WriteString("  help -h               Print this help message.\n")
	// Add other commands and their brief descriptions here
	help.WriteString("Worker Section\n")
	help.WriteString("  Workers can be called by calling bud <worker> <params to worker>\n")
	help.WriteString(`  Example: bud chat "Where is the capital of New Zealand?"` + "\n")
	// fmt.Println("")
	// fmt.Println("Flags:")
	// fmt.Println("  -h, --help  Print this help message.")
	// // Add other flags and their descriptions here
	// fmt.Println("  -<flag>     Describe the specific flag in detail.")
	return help.String()
}
