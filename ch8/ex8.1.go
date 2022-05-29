package ch8

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Port string

func (p *Port) Set(s string) error {
	*p = Port(s)
	return nil
}

func (p *Port) String() string {
	return string(*p)
}

func PortFlag(name string, value string, usage string) *Port {
	f := Port(value)
	flag.CommandLine.Var(&f, name, usage)
	return &f
}

var p = PortFlag("port", "8000", "port number")

func main() {
	flag.Parse()
	port := "localhost:" + p.String()
	fmt.Println(port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
