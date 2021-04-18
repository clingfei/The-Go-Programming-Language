package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Counter int

func (c *Counter) Write(p []byte) (int, error) {
	input := string(p)
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan(){
		*c ++
	}
	return len(p), nil
}

func LimitReader(r io.Reader, n int64) io.Reader {
	reader := r.Read()
}
func main() {
	var c Counter
	w, _ := c.Write([]byte("aaa bbb ccc"))
	fmt.Println(c)
	fmt.Println(w)
}