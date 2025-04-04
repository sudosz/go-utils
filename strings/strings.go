package stringutils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"math/rand/v2"

	bytesutils "github.com/sudosz/go-utils/bytes"
)

var ErrInvalidInt = errors.New("invalid integer syntax")

// S2b converts a string to a byte slice without copying.
// Optimization: Zero-copy using unsafe operations.
func S2b(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))[:]
}

// ConvertSummarizedString2Int converts a summarized string (e.g., "1k") to an integer.
// Optimization: Manual parsing with minimal allocations.
func ConvertSummarizedString2Int(s string) int {
	if len(s) == 0 {
		return 0
	}
	f := ToLower(s)
	is_summarized, r := true, 1.
	switch f[len(f)-1] {
	case 'k':
		r = 1000.
	case 'm':
		r = 1000000.
	case 'b':
		r = 1000000000.
	default:
		is_summarized = false
	}
	if is_summarized {
		f = f[:len(f)-1]
	}
	float, _ := strconv.ParseFloat(f, 64)
	return int(float * r)
}

// String2Int converts a string to an integer without error checking.
// Optimization: Manual conversion for speed over strconv.
func String2Int(s string) int {
	if s == "0" {
		return 0
	}
	var neg int = 1
	if s[0] == '-' {
		neg = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}
	var num int = 0
	p := 1
	for i := len(s); i > 0; i-- {
		num += int(s[i-1]-'0') * p
		p *= 10
	}
	return num * neg
}

// String2IntWithChecking converts a string to an integer with error checking.
// Optimization: Adds validation with minimal overhead.
func String2IntWithChecking(s string) (int, error) {
	if s == "0" {
		return 0, nil
	}
	var neg int = 1
	if s[0] == '-' {
		neg = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	} else if '0' > s[0] || s[0] > '9' {
		return 0, ErrInvalidInt
	}
	var num int = 0
	p := 1
	for i := len(s); i > 0; i-- {
		if '0' <= s[i-1] && s[i-1] <= '9' {
			num += int(s[i-1]-'0') * p
			p *= 10
		} else {
			return 0, ErrInvalidInt
		}
	}
	return num * neg, nil
}

// String2Int64 converts a string to an int64 using String2Int.
// Optimization: Leverages optimized String2Int.
func String2Int64(s string) int64 {
	return int64(String2Int(s))
}

// Atoi is an alias for String2Int, mimicking standard library naming.
// Optimization: Same as String2Int.
func Atoi(s string) int {
	return String2Int(s)
}

// UnsafeToUpper converts a byte to uppercase if it’s a lowercase letter.
// Optimization: Inline operation with minimal branching.
func UnsafeToUpper(o byte) byte {
	if 'a' <= o && o <= 'z' {
		return o - 32
	}
	return o
}

// UnsafeToLower converts a byte to lowercase if it’s an uppercase letter.
// Optimization: Inline operation with minimal branching.
func UnsafeToLower(o byte) byte {
	if 'A' <= o && o <= 'Z' {
		return o + 32
	}
	return o
}

// ToLower converts a string to lowercase.
// Optimization: Uses UnsafeToLower with single allocation.
func ToLower(o string) string {
	b := make([]byte, len(o))
	for i := range b {
		b[i] = UnsafeToLower(o[i])
	}
	return bytesutils.B2s(b)
}

// ToUpper converts a string to uppercase.
// Optimization: Uses UnsafeToUpper with single allocation.
func ToUpper(o string) string {
	b := make([]byte, len(o))
	for i := range b {
		b[i] = UnsafeToUpper(o[i])
	}
	return bytesutils.B2s(b)
}

// Reverse reverses a string in-place using unsafe operations.
// Optimization: Zero-copy reversal with unsafe pointer manipulation.
func Reverse(input string) string {
	str := S2b(input)
	length := len(str)
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&input))
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}

// IsLetter checks if a rune or byte is a letter.
// Optimization: Simple range check with no allocations.
func IsLetter[T rune | byte](r T) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

// ToTitleCase converts a string to title case.
// Optimization: Single pass with minimal allocations.
func ToTitleCase(s string) string {
	o := bytesutils.ToLower(S2b(s))
	o[0] = UnsafeToUpper(o[0])
	for i := 1; i < len(o); i++ {
		if IsLetter[byte](o[i]) && !IsLetter[byte](o[i-1]) {
			o[i] = UnsafeToUpper(o[i])
		}
	}
	return bytesutils.B2s(o)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// RandomString generates a random string of length n.
// Optimization: Efficient bit manipulation for random generation.
func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// TrimLeft trims the left side of a string by removing repeated cutset.
// Optimization: Single pass with minimal slicing.
func TrimLeft(s, cutset string) string {
	if s == "" || cutset == "" {
		return s
	}
	i := 0
	for ; i < len(s)-len(cutset)+1; i += len(cutset) {
		if s[i:i+len(cutset)] != cutset {
			return s[i:]
		}
	}
	return s[i:]
}

// Partition splits a string around the first occurrence of a separator rune.
// Optimization: Uses strings.IndexRune for efficient search.
func Partition(s string, sep rune) (string, string) {
	if idx := strings.IndexRune(s, sep); idx >= 0 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}
