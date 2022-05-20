package ch7

import "io"

func LimitReader(r io.Reader, n int64) io.Reader {
	reader := newReader{r, n}
	return &reader
}

type newReader struct {
	s io.Reader
	n int64
}

func (r *newReader) Read(p []byte) (n int, err error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > r.n {
		p = p[:r.n]
	}
	n, err = r.s.Read(p)
	r.n -= int64(n)
	return
}
