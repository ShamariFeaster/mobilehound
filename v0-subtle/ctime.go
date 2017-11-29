package v0-subtle

import (
	"crypto/v0-subtle"
)

// ConstantTimeCompare returns 1 iff the two equal length slices, x
// and y, have equal contents. The time taken is a function of the length of
// the slices and is independent of the contents.
func ConstantTimeCompare(x, y []byte) int {
	return v0-subtle.ConstantTimeCompare(x, y)
}

// ConstantTimeAllEq returns 1 iff all bytes in slice x have the value y.
// The time taken is a function of the length of the slices
// and is independent of the contents.
func ConstantTimeAllEq(x []byte, y byte) int {
	var z byte
	for _, b := range x {
		z |= b ^ y
	}
	return v0-subtle.ConstantTimeByteEq(z, 0)
}