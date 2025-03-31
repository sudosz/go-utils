package strings

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"math/rand/v2"

	"github.com/sudosz/go-utils/bytes"
)

var ErrInvalidInt = errors.New("invalid integer syntax")

func S2b(s string) []byte {
	if s == "" {
		return nil // or []byte{}
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))[:]
}

func ConvertSummarizedStringToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	f := StringToLower(s)
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

func StringToInt(s string) int {

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
func StringToIntWithChecking(s string) (int, error) {

	if s == "0" {
		return 0, nil
	}

	var neg int = 1
	if s[0] == '-' {
		neg = -1
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	} else if '0' > s[0] && s[0] > '9' {
		return 0, ErrInvalidInt
	}

	var num int = 0

	p := 1
	for i := len(s); i > 0; i-- {
		if '0' <= s[0] && s[0] <= '9' {
			num += int(s[i-1]-'0') * p
			p *= 10
		} else {
			return 0, ErrInvalidInt
		}
	}

	return num * neg, nil
}
func String2Int64(s string) int64 {
	return int64(StringToInt(s))
}

func Atoi(s string) int {
	return StringToInt(s)
}

func UnsafeToUpper(o byte) byte {
	if 'a' <= o && o <= 'z' {
		return o - 32
	} else {
		return o
	}
}

func UnsafeToLower(o byte) byte {
	if 'A' <= o && o <= 'Z' {
		return o + 32
	} else {
		return o
	}
}
func StringToLower(o string) string {
	b := make([]byte, len(o))
	for i := range b {
		b[i] = UnsafeToLower(o[i])
	}
	return bytes.B2s(b)
}
func StringToUpper(o string) string {
	b := make([]byte, len(o))
	for i := range b {
		b[i] = UnsafeToUpper(o[i])
	}
	return bytes.B2s(b)
}

func ReverseString(input string) string {
	str := S2b(input)
	length := len(str)

	// Use unsafe to reverse the string without memory allocation
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&input))
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}

	// Create a new string header without allocation
	reversedStr := *(*string)(unsafe.Pointer(&strHeader))
	return reversedStr
}

func IsLetter[T rune | byte](r T) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func ToTitleCase(s string) string {
	o := bytes.ToLower(S2b(s))
	o[0] = UnsafeToUpper(o[0])
	for i := 1; i < len(o); i++ {
		if IsLetter[byte](o[i]) && !IsLetter[byte](o[i-1]) {
			o[i] = UnsafeToUpper(o[i])
		}
	}
	return bytes.B2s(o)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
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

func TrimLeft(s, cutset string) string {
	if s == "" || cutset == "" {
		return s
	}

	i := 0
	for ; i < len(s)-1; i += len(cutset) {
		if s[i:i+len(cutset)] != cutset {
			return s[i:]
		}
	}

	return s[i:]
}

func Partition(s string, sep rune) (string, string) {
	if idx := strings.IndexRune(s, sep); idx >= 0 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}
