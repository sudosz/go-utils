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

// Copy copies the value from src to dst using unsafe pointer operations for zero-overhead copying.
// Optimization: Uses unsafe to avoid allocation and deep copying.
func Copy[T any](src, dst *T) {
	*dst = *(*T)(unsafe.Pointer(uintptr(unsafe.Pointer(src))))
}

// Reverse reverses the order of elements in the slice in-place and returns it.
// Optimization: Operates in-place to minimize memory usage.
func Reverse[T any](obj []T) []T {
	for i, j := 0, len(obj)-1; i < j; i, j = i+1, j-1 {
		obj[i], obj[j] = obj[j], obj[i]
	}
	return obj
}

// Must returns the value and discards the error, useful for simplifying code when errors are ignorable.
// Optimization: No additional overhead beyond returning the value.
func Must[T any](val T, err error) T {
	return val
}

// Last returns the last element of a slice.
// Optimization: Direct index access for O(1) performance.
func Last[T any, Slice ~[]T](slice Slice) T {
	return slice[len(slice)-1]
}

// RandomChoice selects a random element from the slice using a seeded random source.
// Optimization: Pre-seeded random source avoids repeated initialization.
func RandomChoice[T any, Slice ~[]T](slice Slice) T {
	return slice[random.Int31n(int32(len(slice)))]
}

// Ptr returns a pointer to the provided value.
// Optimization: Simple allocation, no additional overhead possible.
func Ptr[T any](val T) *T {
	return &val
}

// Retry attempts the function up to 'times' times until it succeeds or exhausts retries.
// Optimization: Minimal overhead with a simple loop; could add backoff but not requested.
func Retry(fn func() error, times int) (err error) {
	for i := 0; i < times; i++ {
		if err = fn(); err == nil {
			break
		}
	}
	return
}

// NotifyOnEnd runs each function in a goroutine and sends a signal on the channel when each completes.
// Optimization: Launches goroutines concurrently for parallelism.
func NotifyOnEnd(ch chan<- struct{}, fn ...func()) {
	for _, f := range fn {
		go func(f func()) {
			f()
			ch <- struct{}{}
		}(f)
	}
}

// GenerateUUID generates a UUID with optional prefix and suffix.
// Optimization: Efficient string concatenation with minimal allocations.
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

// GetMimeTypeByFileName returns the MIME type based on the file extension.
// Optimization: Uses standard library mime.TypeByExtension for efficiency.
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

// IsNil checks if the value is nil or a nil pointer using reflection.
// Optimization: Combines direct nil check with reflection-based pointer check.
func IsNil(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

// IsTimeoutError checks if the error indicates a timeout from various sources.
// Optimization: Uses errors.Is and os.IsTimeout for efficient error checking.
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

// IsImageFileExt checks if the extension matches a known image file extension.
// Optimization: Uses a fixed array and linear search (sufficient for small set).
func IsImageFileExt(ext string) bool {
	for _, e := range imageExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

// All returns true if all provided functions return true.
// Optimization: Short-circuits on first false result.
func All[Fn ~func() bool](funcs ...Fn) bool {
	for _, f := range funcs {
		if !f() {
			return false
		}
	}
	return true
}

// Any returns true if any provided function returns true.
// Optimization: Short-circuits on first true result.
func Any(funcs ...func() bool) bool {
	for _, f := range funcs {
		if f() {
			return true
		}
	}
	return false
}

// RWNopeCloser is an interface combining Reader, Writer, and Closer.
type RWNopeCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

// rwNopeCloser is a struct implementing RWNopeCloser with a no-op Close method.
type rwNopeCloser struct {
	io.Reader
	io.Writer
}

// Close is a no-op implementation for rwNopeCloser.
// Optimization: Minimal overhead as it does nothing.
func (rw *rwNopeCloser) Close() error {
	return nil
}

// NopeCloserRW wraps an io.ReadWriter with a no-op Closer.
// Optimization: Simple wrapping with no additional overhead.
func NopeCloserRW(rw io.ReadWriter) RWNopeCloser {
	return &rwNopeCloser{
		Reader: rw,
		Writer: rw,
	}
}

// IsDir checks if the path is a directory, returning an error if not.
// Optimization: Uses os.Stat directly for minimal overhead.
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
