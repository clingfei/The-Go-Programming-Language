package ch4

import "fmt"

func emit(strings []string) []string {
	i, j := 0, 1
	for ; j < len(strings); j++ {
		if strings[i] == strings[j] {
			continue
		} else {
			strings[i+1] = strings[j]
			i++
		}
	}
	return strings[:i+1]
}

func main() {
	strings := []string{"hello", "hello", "test"}
	strings = emit(strings)
	for _, v := range strings {
		fmt.Println(v)
	}
}
