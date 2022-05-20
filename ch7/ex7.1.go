package ch7

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type ByteCounter int

type WordsCounter int

type LinesCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func (c *LinesCounter) Write(p []byte) (int, error) {
	return Counter(p, bufio.ScanLines)
}

func (c *WordsCounter) Write(p []byte) (int, error) {
	return Counter(p, bufio.ScanWords)
}

func Counter(p []byte, splitFunc bufio.SplitFunc) (int, error) {
	s := string(p)
	res := 0
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(splitFunc)
	for scanner.Scan() {
		res++
	}
	if scanner.Err() != nil {
		return res, scanner.Err()
	} else {
		return res, nil
	}
}

func main() {
	var w WordsCounter
	var l LinesCounter

	if words, err := w.Write([]byte("Spicy jalapeno pastrami ut ham turducken.\n Lorem sed ullamco, leberkas sint short loin strip steak ut shoulder shankle porchetta venison prosciutto turducken swine.\n Deserunt kevin frankfurter tongue aliqua incididunt tri-tip shank nostrud.\n")); err == nil {
		fmt.Printf("Words: %d\n", words)
	} else {
		fmt.Println(err)
	}
	if lines, err := l.Write([]byte("Spicy jalapeno pastrami ut ham turducken.\n Lorem sed ullamco, leberkas sint short loin strip steak ut shoulder shankle porchetta venison prosciutto turducken swine.\n Deserunt kevin frankfurter tongue aliqua incididunt tri-tip shank nostrud.\n")); err == nil {
		fmt.Printf("Lines: %d\n", lines)
	} else {
		fmt.Println(err)
	}

	var p *bytes.Buffer
	fmt.Printf("type of p: %T\n", p)
	fmt.Printf("value of p: %v\n", p)
	p = new(bytes.Buffer)
	fmt.Printf("type of p: %T\n", p)
	fmt.Printf("value of p: %v\n", p)
}
