package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}


func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

//UnionWith sets s to the union of s and t
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i< len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//String returns the set as a string of the form "{1,2,3}"
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//return the number of elements
func (s *IntSet) Len() (result int) {
	for _, word := range s.words {
		result += popCount(word)
	}
	return result
}

//remove x from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint64(x%64)
	if word >= len(s.words) {
		fmt.Printf("Error, there is not element %v\n", x)
	} else {
		fmt.Println(word)
		s.words[word+1] ^= bit
	}
}

//remove all elements from the set
func (s *IntSet) Clear() {
	s.words = nil
}

//return a copy of the set
func (s *IntSet) Copy() *IntSet {
	var result IntSet
	result.words = make([]uint64, len(s.words))
	for i := 0; i < len(s.words); i++ {
		result.words[i] = s.words[i]
	}
	return &result
}

//count the number of 1 in uint64 element
func popCount(n uint64) int {
	return int(pc[byte(n>>0*8)] +
		pc[byte(n>>1*8)] +
		pc[byte(n>>2*8)] +
		pc[byte(n>>3*8)] +
		pc[byte(n>>4*8)] +
		pc[byte(n>>5*8)] +
		pc[byte(n>>6*8)] +
		pc[byte(n>>7*8)])
}

func main() {
	var x, y IntSet
	x.Add(1)
	x.Add(14)
	x.Add(9)
	fmt.Println(x.String())
	r := x.Copy()
	fmt.Println(r.String())
	x.Clear()
	fmt.Println(x.String())
	x.Remove(14)
	fmt.Println(x.String())
	fmt.Println(x.Len())
	y.Add(9)
	y.Add(42)
	x.UnionWith(&y)
	fmt.Println(x.String())
	fmt.Println(x.Has(9), x.Has(123))
}