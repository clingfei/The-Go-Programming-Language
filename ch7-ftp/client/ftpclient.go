package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustRead(os.Stdout, conn)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		args := strings.Fields(sc.Text())
		if len(args) == 0 {
			continue
		}
		cmd := args[0]
		switch cmd {
		case "ls", "cd", "get": if _, err := fmt.Fprintln(conn, sc.Text()); err != nil {return}
		case "send": if len(args) < 2 {
						fmt.Println("Not enough arguments.")
					} else {
						data, err := ioutil.ReadFile(args[1])
						if err != nil {
							log.Println("Read error!")
							log.Println(err)
						} else {
							fmt.Println(sc.Text())
							if _, err := fmt.Fprintln(conn, sc.Text(), countLines(string(data))); err != nil {return}
							if _, err := fmt.Fprintln(conn, string(data)); err != nil {return}
						}
					}
		case "close": fmt.Fprintln(conn, args[0])
					return
		}
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
	}
}

func mustRead(dst io.Writer, src io.Reader) {
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

func countLines(data string) (result int) {
	result = 0
	for  _, v := range data {
		if v == '\n' {
			result++
		}
	}
	return result
}