package ch8

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)

	ch := make(chan struct{})
	go func() {
		for input.Scan() {
			ch <- struct{}{}
		}
	}()

	for {
		select {
		case <-time.After(10 * time.Second):
			return
		case <-ch:
			go echo(c, input.Text(), 1*time.Second)
		}
	}
	//var wg sync.WaitGroup
	//for input.Scan() {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		echo(c, input.Text(), 1*time.Second)
	//	}()
	//}
	//// NOTE: ignoring potential errors from input.Err()
	//c.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
