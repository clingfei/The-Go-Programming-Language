package ch4

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"strconv"
)

func hash(input string, v int) {
	if v == 0 {
		fmt.Println(sha256.Sum256([]byte(input)))
	} else if v == 1 {
		fmt.Println(sha512.Sum384([]byte(input)))
	} else {
		fmt.Println(sha512.Sum512([]byte(input)))
	}
}

func main() {
	//var s, sep string
	v := 0
	var err error
	//var hash []byte
	var input string
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-n" {
			if i == len(os.Args)-1 {
				fmt.Println("wrong parameters")
				return
			}
			v, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Println("wrong parameters")
			}
		}
	}
	fmt.Println("input")
	_, err = fmt.Scanln(&input)
	fmt.Println(input)
	fmt.Println(v)
	hash(input, v)
}
