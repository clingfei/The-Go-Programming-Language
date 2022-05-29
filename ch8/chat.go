package ch8

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client chan<- string

type clientInfo struct {
	rchan chan<- string
	cli   client
	name  string
}

//我需要维护当前已经加入聊天的用户名，那么就需要在entering和leaving时删除或添加某些东西

var (
	entering = make(chan clientInfo)
	leaving  = make(chan clientInfo)
	messages = make(chan string)
)

var tick <-chan time.Time

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]string)
	for {
		select {
		case msg := <-messages:
			//message写入管道
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli.cli] = cli.name
			rch := cli.rchan
			for _, v := range clients {
				rch <- v
			}
		case cli := <-leaving:
			delete(clients, cli.cli)
			close(cli.cli)
		}
	}
}

func handleConn(conn net.Conn) {
	tick = time.Tick(5 * time.Minute)
	//为每个client创建一个用于写的管道
	ch := make(chan string)
	go clientWriter(conn, ch, tick)

	rch := make(chan string)
	go clientReader(conn, rch, tick)

	who := conn.RemoteAddr().String()
	var name []byte
	read, err := conn.Read(name)
	if err != nil {
		return
	}
	ch <- "You are" + string(name)
	messages <- string(name) + " has arrived"
	entering <- clientInfo{rchan: rch, cli: ch, name: string(name)}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	leaving <- clientInfo{cli: ch, name: who}
	messages <- who + " has left"
	conn.Close()
}

//用来向conn中发送消息，消息来自于channel ch
//range依次从channel中接收数据，直到没有数据为止
//我好像悟了，这个程序分别在两个不同的进程中运行
func clientWriter(conn net.Conn, ch <-chan string, tick <-chan time.Time) {
	select {
	case msg := <-ch:
		_, _ = fmt.Fprintln(conn, msg)
	case <-tick:
		_ = conn.Close()
	}
	/*	for msg := range ch {
		_, _ = fmt.Fprintln(conn, msg)
	}*/
}

func clientReader(conn net.Conn, ch <-chan string, tick <-chan time.Time) {
	select {
	case msg := <-ch:
		fmt.Fprintln(conn, msg)
	case <-tick:
		conn.Close()
	}
	/*	for msg := range ch {
		_, _ = fmt.Fprintln(conn, msg)
	}*/
}
