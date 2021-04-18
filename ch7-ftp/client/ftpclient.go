package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		args := strings.Fields(sc.Text())
		if len(args) == 0 {
			continue
		}
		cmd := args[0]
		switch cmd {
		case "ls", "cd": if _, err := fmt.Fprintln(conn, sc.Text()); err != nil {return}
		case "get": if _, err := fmt.Fprintln(conn, sc.Text()); err != nil {break}
					file, _ := os.OpenFile(args[1], os.O_CREATE | os.O_RDWR, 0644)
					scanner := bufio.NewScanner(conn)
					scanner.Scan()
					flag := scanner.Text()
					if flag == "Error" {
						scanner.Scan()
						break
					} else {
						n, _ := strconv.Atoi(flag)
						for i := 0; i <= n && scanner.Scan(); i++ {
							//fmt.Println(scanner.Text())
							_, err := fmt.Fprintln(file, scanner.Text())
							if err !=nil {
								fmt.Println(err)
							}
						}
					}
					file.Close()


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

// Return the lines of file
func countLines(data string) (result int) {
	result = 0
	for  _, v := range data {
		if v == '\n' {
			result++
		}
	}
	return result
}