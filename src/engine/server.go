package engine

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

type ServerProperties struct {
	Host string
	Port string
}

func NewServerEngine(host, port string) *ServerProperties {
	return &ServerProperties{
		Host: host,
		Port: port,
	}
}

func (e *Engine) StartServer() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", e.Host, e.Port))
	if err != nil {
		log.Printf("Failed to bind to port %s", e.Port)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go e.HandleConn(conn)
	}
}

func (e *Engine) HandleConn(conn net.Conn) {
	defer conn.Close()

	// This is so we can read everytime the user sends a message
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		buf = bytes.Trim(buf, "\x00") // \x00 is NULL

		if len(buf) > 0 {
			// conn.Write([]byte(e.ParseCommand(buf[:n], conn)))
			conn.Write(e.CliParser(string(buf[:n])))
		}
	}
}
