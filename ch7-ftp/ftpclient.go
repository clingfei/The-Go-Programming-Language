package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustRead(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
	fmt.Println("zzz")

	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
	}
}

func mustRead(dst io.Writer, src io.Reader) {
	/*result, _ := ioutil.ReadAll(src)
	fmt.Println(string(result))
*/
	io.Copy(dst, src)
}

func getFile(buffer []byte) {
	filename := os.Args[1]
	fmt.Println(filename)
	err := ioutil.WriteFile("D:\\Go_workspace\\src\\The-Go-Programming-Language\\"+"file.txt", buffer[:], 777)	//创建文件并写入
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Get Successfully!")
	}
}