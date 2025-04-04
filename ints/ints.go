package intutils

import (
	"strconv"

	bytesutils "github.com/sudosz/go-utils/bytes"
)

const maxIntBufferSize = 19

// Int2Hex converts an integer to a 4-digit hexadecimal string with leading zeros.
// Optimization: Uses strconv.AppendUint for efficient conversion.
func Int2Hex(n int) string {
	b := []byte{'0', '0', '0', 3 + 4: 0}
	b = strconv.AppendUint(b[:3], uint64(n)&0xFFFF, 16)
	return string(b[len(b)-4:])
}

// Int2Bytes converts an integer to a byte slice of specified size.
// Optimization: Manual bit shifting avoids unnecessary allocations.
func Int2Bytes(i, size int) []byte {
	num := int64(i)
	bytes := make([]byte, size)
	idx := len(bytes) - 1
	for ; idx >= 0; idx-- {
		bytes[idx] = byte(num & 0xFF)
		num >>= 8
	}
	return bytes[idx+1:]
}

// Int64ToBytes converts an int64 to a byte slice representing its string form.
// Optimization: Uses a fixed-size buffer and manual conversion for speed.
func Int64ToBytes(i int64) []byte {
	if i == 0 {
		return bytesutils.S2b("0")
	}
	var buf [maxIntBufferSize]byte
	idx := maxIntBufferSize - 1
	negative := 1
	if i < 0 {
		negative = 0
		i = -i
	}
	for i > 0 {
		buf[idx] = byte(i%10) + '0'
		i /= 10
		idx--
	}
	if negative == 0 {
		buf[idx] = '-'
	}
	return buf[idx+negative:]
}

// Int64ToString converts an int64 to a string using Int64ToBytes.
// Optimization: Zero-copy conversion via B2s.
func Int64ToString(i int64) string {
	return bytesutils.B2s(Int64ToBytes(i))
}

// Int2String converts an int to a string using Int64ToString.
// Optimization: Leverages optimized Int64ToString.
func Int2String(i int) string {
	return Int64ToString(int64(i))
}

// Itoa is an alias for Int2String for compatibility with standard library naming.
// Optimization: Same as Int2String.
func Itoa(i int) string {
	return Int64ToString(int64(i))
}
