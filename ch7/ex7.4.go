package ch7

import (
	"io"
)

type StringReader string

func NewReader(s string) io.Reader {
	sr := StringReader(s)
	return &sr
}

func (s *StringReader) Read(p []byte) (n int, err error) {
	n = copy(p, *s)
	return n, io.EOF
}
