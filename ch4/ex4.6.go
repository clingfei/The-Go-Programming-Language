package ch4

import (
	"fmt"
	"unicode"
)

func delSpace(str []byte) []byte {
	i, j := 0, 1
	for j < len(str) {
		if unicode.IsSpace(rune(str[i])) && unicode.IsSpace(rune(str[j])) {
			j++
		} else {
			str[i+1] = str[j]
			i++
			j++
		}
	}
	return str[:i+1]
}

func main() {
	str := []byte{'h', 'e', ' ', ' ', 'o'}
	str = delSpace(str)
	for _, v := range str {
		fmt.Printf("%c", v)
	}
}
