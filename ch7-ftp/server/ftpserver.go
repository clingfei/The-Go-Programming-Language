package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var rootDir string;

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		go handleConn(conn)
	}
}

func buildConn(c net.Conn) {
	fmt.Fprintf(c, "Connection built.\n")
}

func handleConn(c net.Conn) {
	defer c.Close()
	buildConn(c)
	curDir := getDir()
	root := curDir
	sc := bufio.NewScanner(c)
	for sc.Scan(){
		fmt.Println(sc.Text())
		args := strings.Fields(sc.Text())
		switch args[0] {
			case "cd" : handleCd(c, args[1], &curDir, root)
			case "ls" :	handleLs(c, curDir)
			case "get" : handleGet(c, args[1])
			case "send" :
				filename := args[1]
				n, _ := strconv.Atoi(args[2])
				file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 777)
				if err != nil {
					fmt.Println(err)
				}
				for n >= 0 && sc.Scan(){
					fmt.Printf("%s", sc.Text())
					n--
					fmt.Fprintf(file, "%s", sc.Text())
					fmt.Fprintf(file, "%s", '\n')
				}
				file.Close()
				fmt.Fprintln(c, "Transfer successfully!")
			case "close": return
		default : fmt.Fprintln(c, "Error!Undefined operation.")
		}
	}
}

// Handle Get Operation
func handleGet(dst io.Writer, file string) {
	if data, err := ioutil.ReadFile(file); err != nil {
		fmt.Fprintln(dst, "Error")
		fmt.Fprintln(dst, err)
	} else {
		fmt.Println(data)
		n := countLines(string(data))
		fmt.Println(n)
		fmt.Fprintln(dst, n)
		fmt.Fprintln(dst, n)
		fmt.Fprintln(dst, string(data))
	}
	data, err := ioutil.ReadFile(file)
	fmt.Fprintln(dst, file)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(dst, string(data))
	}
}

// Handle ls operation
func handleLs(dst io.Writer, curDir string) {
	ls, err := listFile(curDir)
	if err != nil {
		sendLs(dst, []string{"Error, file not found."})
	} else {
		sendLs(dst, ls)
	}
}


// Handle cd operation
func handleCd(dst io.Writer, dir string, curDir *string, root string) {
	if dir == ".." {
		if (*curDir) == root {
			fmt.Fprintln(dst, "Illegal Path")
		} else {
			for i := len(*curDir)-1; i>=0; i-- {
				if string((*curDir)[i]) == "\\" {
					*curDir = (*curDir)[:i]
					break
				}
			}
		}
	} else {
		path := (*curDir) + "\\" + dir
		_, err := listFile(path)
		if err != nil {
			sendLs(dst, []string{"Error, directory not found."})
		} else {
			(*curDir) = (*curDir) + "\\" + dir
		}
	}
	if flag := len(*curDir)==len(root); flag {
		fmt.Fprintln(dst, "Current Path: ", "\\")
	} else {
		fmt.Fprintln(dst, "Current Path: ", (*curDir)[len(root):])
	}
}

// Return the list of file specified
func listFile(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file.Name()+"\\")
		} else {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

// Get the current dir of exe
func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return  dir
}

// Send the list of file to client
func sendLs(dst io.Writer, list []string) {
	for _, file := range list {
		io.WriteString(dst, file+"\n")
	}
}

// Return the lines of file to client
func countLines(data string) (result int) {
	result = 0
	for  _, v := range data {
		if v == '\n' {
			result++
		}
	}
	return result
}