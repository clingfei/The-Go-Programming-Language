package ch7

import (
	"fmt"
	"io"
)

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	bc := ByteCounter{w, 0}
	//Write函数是对类型*ByteCounter即ByteCounter类型的指针而不是ByteCounter类型定义的，因此只有*ByteCounter可以被视为io.Writer
	return &bc, &bc.cnt
}

type ByteCounter struct {
	w   io.Writer
	cnt int64
}

func (c ByteCounter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.cnt += int64(n)
	return n, err
}

type Point struct {
	x int
	y int
}

func (t *Point) add() {
	t.y += 1
	t.x += 1
}

func (t Point) sub() {
	t.y += 1
	t.x += 1
}

func main() {
	t := Point{1, 2}
	t.add()
	t.sub()
	p := Point{1, 2}
	p.sub()
	fmt.Println(t)
	fmt.Println(p)
}
