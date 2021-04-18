package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
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
	fmt.Fprintf(c, "Connection built")
}

func handleConn(c net.Conn) {
	defer c.Close()
	var buf [512]byte
	buildConn(c)
	curDir := getDir()
	for {
		n, err := c.Read(buf[0:])

		if err != nil {
			fmt.Println(err)
			return
		} else {
			s := string(buf[0 : n-2])
			fmt.Println(s)
			context := strings.Fields(s)
			switch context[0] {
			case "cd" : handleCd(c, context[1], &curDir)
			case "ls" :	handleLs(c, curDir)
			case "get" : handleGet(c, context[1], curDir)
			case "send" : handleSend(c, context[1])
			}
		}
	}
}

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

func handleSend(dst io.Writer, file string) {

}

func handleLs(dst io.Writer, curDir string) {
	ls, err := listFile(curDir)
	if err != nil {
		sendLs(dst, []string{"Error, file not found."})
	} else {
		sendLs(dst, ls)
	}
}

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
			(*curDir) = (*curDir) + "\\" + dir
			//handleLs(dst, *curDir)
		}
	}
}

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

func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return  dir
}

func sendLs(dst io.Writer, list []string) {
	for _, file := range list {
		io.WriteString(dst, file+"\n")
	}
}