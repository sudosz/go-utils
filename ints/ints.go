package ints

import (
	"strconv"

	"github.com/sudosz/go-utils/bytes"
)

const maxIntBufferSize = 19

func Int2Hex(n int) string {
	b := []byte{'0', '0', '0', 3 + 4: 0}
	b = strconv.AppendUint(b[:3], uint64(n)&0xFFFF, 16)
	return string(b[len(b)-4:])
}

func Int2Bytes(i, size int) []byte {

	num := int64(i)
	bytes := make([]byte, size)
	idx := len(bytes) - 1
	for ; idx > 0; idx-- {
		bytes[idx] = byte(num & 0xFF)
		num >>= 8
	}
	return bytes[idx+1:]
}

func Int64ToBytes(i int64) []byte {

	if i == 0 {
		return bytes.S2b("0")
	}

	// Preallocate a buffer with a capacity large enough for the maximum integer value.
	var buf [maxIntBufferSize]byte

	// Start from the end of the buffer.
	idx := maxIntBufferSize - 1

	// Handle negative numbers.
	negative := 1
	if i < 0 {
		negative = 0
		i = -i
	}

	// Convert the integer to a string in reverse order.
	for i > 0 {
		buf[idx] = byte(i%10) + '0'
		i /= 10
		idx--
	}

	// Add the negative sign if necessary.
	if negative == 0 {
		buf[idx] = '-'
	}

	return buf[idx+negative:]
}

func Int64ToString(i int64) string {
	return bytes.B2s(Int64ToBytes(i))
}

func Int2String(i int) string {
	return Int64ToString(int64(i))
}

func Itoa(i int) string {
	return Int64ToString(int64(i))
}
