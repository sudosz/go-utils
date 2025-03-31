package utils

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"mime"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

var (
	src    = rand.NewSource(time.Now().UnixNano())
	random = rand.New(src)
)

func Copy[T any](src, dst *T) {
	*dst = *(*T)(unsafe.Pointer(uintptr(unsafe.Pointer(src))))
}

func Reverse[T any](obj []T) []T {
	for i, j := 0, len(obj)-1; i < j; i, j = i+1, j-1 {
		obj[i], obj[j] = obj[j], obj[i]
	}
	return obj
}

func Must[T any](val T, err error) T {
	return val
}

func Last[T any, Slice ~[]T](slice Slice) T {
	return slice[len(slice)-1]
}

func RandomChoice[T any, Slice ~[]T](slice Slice) T {
	return slice[random.Int31n(int32(len(slice)))]
}

func Ptr[T any](val T) *T {
	return &val
}

func Retry(fn func() error, times int) (err error) {
	for i := 0; i < times; i++ {
		if err = fn(); err == nil {
			break
		}
	}
	return
}

func NotifyOnEnd(ch chan<- struct{}, fn ...func()) {
	for _, f := range fn {
		go func(f func()) {
			f()
			ch <- struct{}{}
		}(f)
	}
}

func GenerateUUID(ps ...string) string {
	p, s := "", ""
	if len(ps) > 0 {
		p = ps[0]
		if len(ps) > 1 {
			s = ps[1]
		}
	}
	return p + uuid.New().String() + s
}

func GetMimeTypeByFileName(filename string) string {
	ext := "." + Last(strings.Split(filename, "."))
	mimeType := ""
	if ext == "." {
		mimeType = "application/octet-stream"
	} else {
		mimeType = mime.TypeByExtension(ext)
	}
	return mimeType
}

func IsNil(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func IsTimeoutError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if errors.Is(errors.Unwrap(err), context.DeadlineExceeded) {
		return true
	}
	if os.IsTimeout(err) {
		return true
	}
	if os.IsTimeout(errors.Unwrap(err)) {
		return true
	}
	return false
}

var imageExtensions = [...]string{
	".png",
	".bmp",
	".jpg",
	".jpeg",
}

func IsImageFileExt(ext string) bool {
	for _, e := range imageExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func All[Fn ~func() bool](funcs ...Fn) bool {
	for _, f := range funcs {
		if !f() {
			return false
		}
	}
	return true
}

func Any(funcs ...func() bool) bool {
	for _, f := range funcs {
		if f() {
			return true
		}
	}
	return false
}

type RWNopeCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

type rwNopeCloser struct {
	io.Reader
	io.Writer
}

func (rw *rwNopeCloser) Close() error {
	return nil
}

func NopeCloserRW(rw io.ReadWriter) RWNopeCloser {
	return &rwNopeCloser{
		Reader: rw,
		Writer: rw,
	}
}

func IsDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("not a directory")
	}
	return nil
}
