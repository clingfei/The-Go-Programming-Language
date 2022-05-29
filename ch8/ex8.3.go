package ch8

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	tcpConn := conn.(*net.TCPConn)

	done := make(chan struct{})
	go func() {
		//这里是从连接中读数据
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)

	tcpConn.CloseWrite()

	//应该只关闭写的部分
	<-done // wait for background goroutine to finish
	tcpConn.CloseRead()
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
