package bytes

import "unsafe"

func ReverseBytes(obj []byte) {
	for i, j := 0, len(obj)-1; i < j; i, j = i+1, j-1 {
		obj[i], obj[j] = obj[j], obj[i]
	}
}

func B2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func S2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func ToLower(b []byte) []byte {
	for i := range b {
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return b
}

func ToUpper(b []byte) []byte {
	for i := range b {
		if 'a' <= b[i] && b[i] <= 'z' {
			b[i] -= 32
		}
	}
	return b
}