package ch7

import "strconv"

type Tree interface {
	String() string
}

type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	if t == nil {
		return ""
	}
	return t.left.String() + " " + strconv.Itoa(t.value) + " " + t.right.String()
}
