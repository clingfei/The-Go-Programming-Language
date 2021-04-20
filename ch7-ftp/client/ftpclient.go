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
		case "ls", "cd": if _, err := fmt.Fprintln(conn, sc.Text()); err != nil { return }
		case "get": if err := get(conn, args, sc); err != nil { break }
		case "send": if err := send(conn, args, sc); err != nil { return }
		case "close": fmt.Fprintln(conn, args[0]); 	return
		default: fmt.Println("Unsupproted operations! Please try again.")
		}
	}
}

// Send files to server
func send(conn net.Conn, args []string, sc *bufio.Scanner) (error){
	if len(args) < 2 {
		fmt.Println("Not enough arguments.")
	} else {
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			log.Println("Read error!")
			log.Println(err)
		} else {
			fmt.Println(sc.Text())
			if _, err := fmt.Fprintln(conn, sc.Text(), countLines(string(data))); err != nil { return err }
			if _, err := fmt.Fprintln(conn, string(data)); err != nil { return err }
		}
	}
	return nil
}

// Get files from server.
func get(conn net.Conn, args []string, sc *bufio.Scanner) (err error){
	if _, err := fmt.Fprintln(conn, sc.Text()); err != nil { return err }
	file, _ := os.OpenFile(args[1], os.O_CREATE | os.O_RDWR, 0644)
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	flag := scanner.Text()
	if flag == "Error" {
		scanner.Scan()
		return err
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
	return nil
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