package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint
}

const base = 32 << (^uint(0) >> 63)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/base, uint(x%base)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

//Add addds the non-negative  value x to the set
func (s *IntSet) Add(x int) {
	word, bit := x/base, uint(x%base)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

//Add a slice of the non-negative value to the set
func (s *IntSet) AddAll(words ...int) {
	for _, i := range words {
		s.Add(i)
	}
}

//IntersectWith sets s to the intersect  of s and t
func (s *IntSet) IntersectWith(t *IntSet) {
	if len(s.words) > len(t.words) {
		s.words = s.words[:len(t.words)]
	}
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		}
	}
}

//UnionWith sets s to the union of s and t
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//DifferenceWith sets s to the difference of s and t
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= s.words[i] & tword
		}
	}
}
func (s *IntSet) SymmetricDifference(t *IntSet) {
	var s1 IntSet
	for _, v := range s.words {
		s1.words = append(s1.words, v)
	}
	s.UnionWith(t)
	s1.IntersectWith(t)
	s.DifferenceWith(&s1)
}

//String returns the set as a string of the form "{1,2,3}"
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < base; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", base*i+j)
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
	word, bit := x/base, uint(x%base)
	if word >= len(s.words) {
		fmt.Printf("Error, there is not element %v\n", x)
	} else {
		s.words[word] ^= (1 << bit)
	}
}

//remove all elements from the set
func (s *IntSet) Clear() {
	s.words = nil
}

//return a copy of the set
func (s *IntSet) Copy() *IntSet {
	var result IntSet
	result.words = make([]uint, len(s.words))
	for i := 0; i < len(s.words); i++ {
		result.words[i] = s.words[i]
	}
	return &result
}

//count the number of 1 in uint element
func popCount(x uint) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func main() {
	var x, y IntSet
	x.Add(1)
	x.Add(14)
	x.Add(9)
	x.AddAll(3, 4, 5)
	//fmt.Println(x.String())
	//r := x.Copy()
	//fmt.Println(r.String())
	//x.Remove(14)
	//fmt.Println(x.String())
	//fmt.Println(x.Len())
	//x.Clear()
	//fmt.Println(x.String())
	y.Add(9)
	y.Add(42)
	fmt.Println(x.String())

	x.IntersectWith(&y)
	fmt.Println(x.String())
	x.AddAll(9, 43, 58)
	fmt.Println(x.String())
	x.DifferenceWith(&y)
	fmt.Println(x.String())
	x.AddAll(10, 34, 9, 42)
	fmt.Println(x.String())
	x.SymmetricDifference(&y)
	fmt.Println(x.String())
}
