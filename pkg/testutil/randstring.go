package testutil

// this code is take from github.com/erikh/spin/pkg/testutil and is licensed
// MIT at the time of copy.

import (
	"crypto/rand"
)

// RandomString generates a random string which has a length of the provided
// amount with a guaranteed minimum. It chooses from a-zA-Z0-9. It's very
// simple.
func RandomString(length, min uint) string {
	buf := make([]byte, length+1)
	n, err := rand.Reader.Read(buf)
	if err != nil {
		panic(err)
	}

	if uint(n) != length+1 {
		panic("short read on random device")
	}

	// we use the first random byte to determine the length of the rest of the
	// string. We cap it by subtracting the minimum before the mod then adding it
	// back in.

	l := length
	if min < length {
		l = length - min
	}

	l = uint(buf[0]) % l

	if min < length {
		length = l + min
	}

	str := []byte{}

	// Because we're using the first byte above, the offsets must be calculated.
	for _, b := range buf[1 : length+1] {
		c := b % 62
		// above 52 is a number; calculate from '0'. above 26 is a capital letter,
		// so calculate the offset from 'A'. Remember, in ascii, latin alphabet is
		// listed in order, so we can just add to get the letters.
		if c >= 52 {
			str = append([]byte(str), byte('0')+byte(c-52))
		} else if c >= 26 {
			str = append([]byte(str), byte('A')+byte(c-26))
		} else {
			str = append([]byte(str), byte('a')+byte(c))
		}
	}

	return string(str)
}
