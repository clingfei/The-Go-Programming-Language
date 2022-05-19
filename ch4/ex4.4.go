package ch4

import "fmt"

func rotate(arr []int, r int) []int {
	res := make([]int, len(arr), len(arr))
	for i := 0; i < len(arr); i++ {
		res[(i+r)%len(arr)] = arr[i]
	}
	return res
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	arr = rotate(arr, 1)
	for _, v := range arr {
		fmt.Print(v)
	}
}
