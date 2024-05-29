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
		case k:
			log.Println(k)
			v.Call(c[1:])
			if conn != nil {
				fmt.Fprintf(conn, "Worker %s called\n", k)
			}
		case "start":
			if c[1] == k {
				go v.Run()
				fmt.Fprintf(conn, "STARTED WORKER %s\n", k)
			}
		case "stop":
			if c[1] == k {
				go v.Stop()
				fmt.Fprintf(conn, "STOPPED WORKER %s\n", k)
			}
		case "kill":
			if c[1] == k {
				v.Kill()
				delete(Workers, k)
			}
		}
	}
}
