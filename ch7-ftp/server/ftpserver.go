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
	sc := bufio.NewScanner(c)
	for sc.Scan(){
		fmt.Println(sc.Text())
		args := strings.Fields(sc.Text())
		switch args[0] {
			case "cd" : handleCd(c, args[1], &curDir)
			case "ls" :	handleLs(c, curDir)
			case "get" : handleGet(c, args[1], curDir)
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
func handleGet(dst io.Writer, file string, dir string) {
	var path string
	if file[:2] != "C:" && file[:2] != "D:" {
		path = dir + "\\" + file
	} else {
		path = file
		for i := len(file)-1; i>=0; i-- {
			if file[i] == '\\' {
				file = file[i:]
				break
			}
		}
	}
	fileObj, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(dst, err)
		fmt.Fprintf(dst, "Error, read failed.\n")
	} else {
		io.WriteString(dst, "ok")
		io.WriteString(dst, file + "\n")
		io.Copy(dst, fileObj)
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
func handleCd(dst io.Writer, dir string, curDir *string) {
	if dir[:2] == "C:" || dir[:2] == "D:" {
		list, err := listFile(dir)
		if err != nil {
			sendLs(dst, []string{"Error, file not found."})
		} else {
			sendLs(dst, list)
		}
	} else {
		if dir == ".." {
			for i := len(*curDir)-1; i>=0; i-- {
				if string((*curDir)[i]) == "\\" {
					*curDir = (*curDir)[:i]
					break
				}
			}
			//handleLs(dst, *curDir)
		} else if dir == "." {
			//handleLs(dst, *curDir)
		} else {
			path := (*curDir) + "\\" + dir
			_, err := listFile(path)
			if err != nil {
				sendLs(dst, []string{"Error, directory not found."})
			} else {
				(*curDir) = (*curDir) + "\\" + dir
			}
			//handleLs(dst, *curDir)
		}
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