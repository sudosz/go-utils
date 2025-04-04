package netutils

import (
	"bufio"
	gbytes "bytes"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	bytesutils "github.com/sudosz/go-utils/bytes"
	intutils "github.com/sudosz/go-utils/ints"
	poolutils "github.com/sudosz/go-utils/pool"
	stringutils "github.com/sudosz/go-utils/strings"
)

// AuthProvider defines a function that provides authentication strings.
type AuthProvider func() string

// ProxyCredentialsProvider defines an interface for proxy credential retrieval.
type ProxyCredentialsProvider interface {
	GetProxyCredentials() ([]byte, []byte)
}

const (
	http10OKResponse   = "HTTP/1.0 200 OK\r\n"
	http11OKResponse   = "HTTP/1.1 200 OK\r\n"
	http2OKResponse    = "\x00\x00\x86\x04\x00\x00\x00"
	maxResponseLen     = len(http11OKResponse)
	connectPrefix      = "CONNECT "
	httpVersion        = " HTTP/1.1\r\n"
	hostPrefix         = "Host: "
	headerSep          = ": "
	proxyAuthPrefix    = "Proxy-Authorization: "
	crlf               = "\r\n"
	basicAuthPrefix    = "Basic "
	basicAuthPrefixLen = len(basicAuthPrefix)
	okStartResponse    = "HTTP/"
	okEndResponse      = " 200 OK\r\n\r\n"
	okStatusTotalLen   = len(okStartResponse) + 3 + len(okEndResponse)
)

var (
	bufPool    = poolutils.NewLRULimitedBufferPool(1<<10, 1<<7, 1*time.Minute)
	HopHeaders = [...]string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Connection",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)

// BasicAuthHeader generates a Basic Authentication header from byte slice credentials.
// Optimization: Pre-calculates buffer size for single allocation.
func BasicAuthHeader(username, password []byte) string {
	credLen := len(username) + len(password) + 1
	totalLen := basicAuthPrefixLen + ((credLen+2)/3)*4
	buf := make([]byte, totalLen)
	copy(buf, basicAuthPrefix)
	base64.StdEncoding.Encode(buf[basicAuthPrefixLen:], append(append(username, ':'), password...))
	return bytesutils.B2s(buf)
}

// CachedAuth returns a function caching the auth header for a specified duration.
// Optimization: Uses TryLock to minimize contention on cache updates.
func CachedAuth(provider ProxyCredentialsProvider, dur time.Duration) func() string {
	var (
		lru   time.Time
		lcred string
		mux   sync.Mutex
	)
	return func() string {
		if lru.IsZero() || time.Since(lru) > dur {
			if mux.TryLock() {
				lcred = BasicAuthHeader(provider.GetProxyCredentials())
				lru = time.Now()
				mux.Unlock()
			}
		}
		return lcred
	}
}

// BasicAuthHeaderStr generates a Basic Auth header from string credentials.
// Optimization: Leverages zero-copy S2b for inputs.
func BasicAuthHeaderStr(username, password string) string {
	return BasicAuthHeader(bytesutils.S2b(username), bytesutils.S2b(password))
}

// SimpleAuth returns a function generating auth headers without caching.
// Optimization: Minimal overhead with direct call to BasicAuthHeader.
func SimpleAuth(provider ProxyCredentialsProvider) func() string {
	return func() string {
		return BasicAuthHeader(provider.GetProxyCredentials())
	}
}

// JoinHostPort combines host and port bytes into a single string.
// Optimization: Single allocation with append.
func JoinHostPort(host []byte, port []byte) string {
	b := append(host, ':')
	b = append(b, port...)
	return bytesutils.B2s(b)
}

// JoinHostIntPort combines host bytes and an integer port into a string.
// Optimization: Uses efficient Int64ToBytes for port conversion.
func JoinHostIntPort(host []byte, port int) string {
	return JoinHostPort(host, intutils.Int64ToBytes(int64(port)))
}

// JoinStrHostIntPort combines a string host and integer port into a string.
// Optimization: Zero-copy S2b for host.
func JoinStrHostIntPort(host string, port int) string {
	return JoinHostIntPort(bytesutils.S2b(host), port)
}

// JoinStrHostStrPort combines string host and port into a single string.
// Optimization: Zero-copy S2b for both inputs.
func JoinStrHostStrPort(host string, port string) string {
	return JoinHostPort(bytesutils.S2b(host), bytesutils.S2b(port))
}

// StatusOKBytes generates an HTTP OK status line for the given version.
// Optimization: Pre-allocates buffer with exact size.
func StatusOKBytes(major, minor int) []byte {
	buf := make([]byte, 0, okStatusTotalLen)
	buf = append(buf, okStartResponse...)
	buf = append(buf, '0'+byte(major), '.', '0'+byte(minor))
	buf = append(buf, okEndResponse...)
	return buf
}

// IsHTTPOK checks if the buffer contains an HTTP OK response.
// Optimization: Early length check and switch for efficiency.
func IsHTTPOK(buf []byte, isConnect bool) bool {
	if isConnect && len(buf) == 0 {
		return true
	}
	if len(buf) < maxResponseLen {
		return false
	}
	switch {
	case bytesutils.B2s(buf[:len(http10OKResponse)]) == http10OKResponse:
		return true
	case bytesutils.B2s(buf[:len(http11OKResponse)]) == http11OKResponse:
		return true
	case bytesutils.B2s(buf[:len(http2OKResponse)]) == http2OKResponse:
		return true
	default:
		return false
	}
}

// IsHTTPOKConn reads from a reader to check for an HTTP OK response.
// Optimization: Uses pooled buffer to reduce allocations.
func IsHTTPOKConn(r io.Reader, isConnect bool) (bool, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<5), 1<<9)
	scanner.Split(bufio.ScanLines)
	buf := bufPool.Get().(*gbytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	for scanner.Scan() {
		line := scanner.Bytes()
		line = append(line, '\r', '\n')
		buf.Write(line)
		if len(line) == 2 {
			return IsHTTPOK(buf.Bytes(), isConnect), nil
		}
	}
	if scanner.Err() != nil {
		return false, scanner.Err()
	}
	return IsHTTPOK(buf.Bytes(), isConnect), nil
}

// RawConnectRequestBytes builds a raw CONNECT request with optional authentication.
// Optimization: Uses pooled buffer for minimal allocations.
func RawConnectRequestBytes(address string, proxyAuth func() string) []byte {
	buf := bufPool.Get().(*gbytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.WriteString(connectPrefix)
	buf.WriteString(address)
	buf.WriteString(httpVersion)
	buf.WriteString(hostPrefix)
	buf.WriteString(address)
	buf.WriteString(crlf)
	if proxyAuth != nil {
		buf.WriteString(proxyAuthPrefix)
		buf.WriteString(proxyAuth())
		buf.WriteString(crlf)
	}
	buf.WriteString(crlf)
	return buf.Bytes()
}

const (
	copyBufSize = 1 << 8
)

// Flush flushes the writer if it implements http.Flusher.
// Optimization: Type assertion with minimal overhead.
func Flush(w any) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

var copyBufPool = &poolutils.LimitedPool{
	New: func() any {
		var b [copyBufSize]byte
		return &b
	},
	N: 1 << 6,
}

// CopyBody copies data from src to dst, flushing if possible.
// Optimization: Uses pooled buffer and atomic counter for efficiency.
func CopyBody(dst io.Writer, src io.Reader) (int64, error) {
	buf := copyBufPool.Get().(*[copyBufSize]byte)
	defer copyBufPool.Put(buf)
	var written atomic.Uint32
	for {
		nr, err := src.Read(buf[:])
		if nr > 0 {
			nw, err := dst.Write(buf[:nr])
			written.Add(uint32(nw))
			if nw > 0 {
				Flush(src)
			}
			if err != nil {
				return 0, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return int64(written.Load()), nil
}

// HeaderDeleter defines an interface for deleting headers.
type HeaderDeleter interface {
	Del(string)
}

// DelHopHeaders removes hop-by-hop headers from the provided header.
// Optimization: Linear iteration over fixed array.
func DelHopHeaders(header HeaderDeleter) {
	for _, headerName := range HopHeaders {
		header.Del(headerName)
	}
}

// CopyHTTPHeaders copies headers from src to dst.
// Optimization: Direct slice copy for efficiency.
func CopyHTTPHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst[k] = vv[:]
	}
}

// BuildRequestBytes constructs raw bytes for an HTTP request.
// Optimization: Uses pooled buffer to minimize allocations.
func BuildRequestBytes(req *http.Request) (_ []byte, err error) {
	buf := bufPool.Get().(*gbytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.WriteString(req.Method)
	buf.WriteByte(' ')
	buf.WriteString(req.URL.Path)
	buf.WriteByte(' ')
	buf.WriteString(req.Proto)
	buf.WriteString(crlf)
	buf.WriteString(hostPrefix)
	h := req.Host
	if h == "" {
		h = req.URL.Host
		if h == "" && len(req.Header["Host"]) > 0 {
			h = req.Header["Host"][0]
		}
	}
	buf.WriteString(h)
	buf.WriteString(crlf)
	for k, vv := range req.Header {
		for _, v := range vv {
			buf.WriteString(k)
			buf.WriteString(headerSep)
			buf.WriteString(v)
			buf.WriteString(crlf)
		}
	}
	buf.WriteString(crlf)
	if req.GetBody != nil {
		req.Body, err = req.GetBody()
		if err != nil {
			return nil, err
		}
	}
	if req.Body != nil {
		_, err = io.Copy(buf, req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}
	return buf.Bytes(), nil
}

// IsConnClosedErr checks if the error indicates a closed connection.
// Optimization: Uses errors.Is for efficient error matching.
func IsConnClosedErr(err error) bool {
	switch {
	case errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}

// PartitionIP4 splits an IPv4 address into its four parts.
// Optimization: Uses efficient string partitioning.
func PartitionIP4(ip string) (string, string, string, string) {
	part1, part234 := stringutils.Partition(ip, '.')
	part2, part34 := stringutils.Partition(part234, '.')
	part3, part4 := stringutils.Partition(part34, '.')
	return part1, part2, part3, part4
}
