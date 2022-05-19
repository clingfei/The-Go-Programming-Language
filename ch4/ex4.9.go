package ch4

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println("open file error.")
		os.Exit(-1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}(file)

	counts := make(map[string]int)
	input := bufio.NewScanner(file)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		counts[input.Text()]++
	}
	for k, v := range counts {
		fmt.Printf("%s : %d\n", k, v)
	}
}
