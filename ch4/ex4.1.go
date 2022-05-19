package ch4

import (
	"crypto/sha256"
	"fmt"
)

func count(sh1 [32]byte, sh2 [32]byte) (res int) {
	for i := 0; i < 32; i++ {
		xor := ^(sh1[i] ^ sh2[i])
		res += int(xor&0x1) + int((xor&0x2)>>1) + int((xor&0x4)>>2) + int((xor&0x8)>>3) + int((xor&0x10)>>4) + int((xor&0x20)>>5) + int((xor&0x40)>>6) + int((xor&0x80)>>7)
	}
	return
}

func main() {
	sh1 := sha256.Sum256([]byte("hello world\n"))
	sh2 := sha256.Sum256([]byte("hello wwrld\n"))
	fmt.Printf("the number of different bits: %d", count(sh1, sh2))
}
