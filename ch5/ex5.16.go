package ch5

import "fmt"

func Join(sep string, items ...string) string {
	var res string
	for _, v := range items {
		res += v + sep
	}
	return res[:len(res)-len(sep)]
}

func test() {
	fmt.Println("test start.")
	defer hello()

}

func hello() func() {
	fmt.Println("hello world")
	return func() {
		fmt.Println("hello end")
	}
}

func triple(x int) (result int) {
	defer func() { result += x }()
	result += x
	return result
}

func main() {
	//fmt.Println(Join(" ", "hello", "world"))
	test()
	fmt.Println(triple(1))
}
