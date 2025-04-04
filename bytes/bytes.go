package bytesutils

import "unsafe"

// ReverseBytes reverses the order of bytes in the slice in-place.
// Optimization: In-place operation avoids allocation.
func ReverseBytes(obj []byte) {
	for i, j := 0, len(obj)-1; i < j; i, j = i+1, j-1 {
		obj[i], obj[j] = obj[j], obj[i]
	}
}

// B2s converts a byte slice to a string without copying using unsafe operations.
// Optimization: Zero-copy conversion for maximum efficiency.
func B2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2b converts a string to a byte slice without copying using unsafe operations.
// Optimization: Zero-copy conversion for maximum efficiency.
func S2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ToLower converts all uppercase letters in the byte slice to lowercase in-place.
// Optimization: In-place modification avoids allocation.
func ToLower(b []byte) []byte {
	for i := range b {
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return b
}

// ToUpper converts all lowercase letters in the byte slice to uppercase in-place.
// Optimization: In-place modification avoids allocation.
func ToUpper(b []byte) []byte {
	for i := range b {
		if 'a' <= b[i] && b[i] <= 'z' {
			b[i] -= 32
		}
	}
	return b
}
