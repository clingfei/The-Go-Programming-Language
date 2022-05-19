package ch5

import "fmt"

func max(vals ...int) int {
	res := -1 << 31
	for _, v := range vals {
		if v > res {
			res = v
		}
	}
	if res == -1<<31 {
		return 1<<31 - 1
	} else {
		return res
	}
}

func min(vals ...int) int {
	res := 1<<31 - 1
	for _, v := range vals {
		if v < res {
			res = v
		}
	}
	if res == 1<<31-1 {
		return -1 << 31
	} else {
		return res
	}
}

func max1(val int, vals ...int) int {
	res := val
	for _, v := range vals {
		if v > res {
			res = v
		}
	}
	return res
}

func min1(val int, vals ...int) int {
	res := val
	for _, v := range vals {
		if v < res {
			res = v
		}
	}
	return res
}

func main() {
	fmt.Println(max())
	fmt.Println(min())
}
