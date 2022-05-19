package ch5

import "fmt"

func returnNonzero() (res int) {
	defer func() {
		if p := recover(); p != nil {
			res = 1
		}
	}()
	panic(1)
}

func main() {
	//fmt.Println(returnNonzero())
	a := []int{0}
	test(&a)
	for _, v := range a {
		fmt.Println(v)
	}
}

func test(a *[]int) {
	*a = append(*a, 1)
}
func grayCode(n int) []int {
	/*	if n == 1 {
			return []int{0, 1}
		}
		gray := grayCode(n - 1)
		res := make([]int, 1 << n)
		for i := 0; i < len(gray); i++ {
			res[i] = gray[i];
			res[len(res) - i - 1] = gray[i] ^ (1 << (n - 1))
		}
		return res*/
	res := make([]int, 1<<n)
	gray(n, &res)
	return res
}

func gray(n int, res *[]int) {
	if n == 1 {
		(*res)[0], (*res)[1] = 0, 1
		return
	}
	gray(n-1, res)
	length := 1 << (n - 1)
	for i := 0; i < length; i++ {
		(*res)[(length)<<1-i-1] = (*res)[i] ^ length
	}
}
