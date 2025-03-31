package net

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

	"github.com/sudosz/go-utils/ints"
	"github.com/sudosz/go-utils/pool"
	"github.com/sudosz/go-utils/bytes"
	"github.com/sudosz/go-utils/strings"
)

type AuthProvider func() string

type ProxyCredentialsProvider interface {
	GetProxyCredentials() ([]byte, []byte)
}

const (
	http10OKResponse = "HTTP/1.0 200 OK\r\n"
	http11OKResponse = "HTTP/1.1 200 OK\r\n"
	http2OKResponse  = "\x00\x00\x86\x04\x00\x00\x00"
	maxResponseLen   = len(http11OKResponse)

	connectPrefix   = "CONNECT "
	httpVersion     = " HTTP/1.1\r\n"
	hostPrefix      = "Host: "
	headerSep       = ": "
	proxyAuthPrefix = "Proxy-Authorization: "
	crlf            = "\r\n"

	basicAuthPrefix    = "Basic "
	basicAuthPrefixLen = len(basicAuthPrefix)

	okStartResponse  = "HTTP/"
	okEndResponse    = " 200 OK\r\n\r\n"
	okStatusTotalLen = len(okStartResponse) + 3 + len(okEndResponse) // +3 for "major.minor"
)

var (
	bufPool = pool.NewLRULimitedBufferPool(1<<10, 1<<7, 1*time.Minute)

	// Hop-by-hop headers. These are removed when sent to the backend.
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
	HopHeaders = [...]string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Connection",
		"Te", // canonicalized version of "TE"
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)

func BasicAuthHeader(username, password []byte) string {
	credLen := len(username) + len(password) + 1
	totalLen := basicAuthPrefixLen + ((credLen+2)/3)*4 // Base64 encoding increases length
	buf := make([]byte, totalLen)

	copy(buf, basicAuthPrefix)

	base64.StdEncoding.Encode(buf[basicAuthPrefixLen:], append(append(username, ':'), password...))

	return bytes.B2s(buf)
}

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

func BasicAuthHeaderStr(username, password string) string {
	return BasicAuthHeader(bytes.S2b(username), bytes.S2b(password))
}

func SimpleAuth(provider ProxyCredentialsProvider) func() string {
	return func() string {
		return BasicAuthHeader(provider.GetProxyCredentials())
	}
}

func JoinHostPort(host []byte, port []byte) string {
	b := append(host, ':')
	b = append(b, port...)
	return bytes.B2s(b)
}

func JoinHostIntPort(host []byte, port int) string {
	return JoinHostPort(host, ints.Int64ToBytes(int64(port)))
}

func JoinStrHostIntPort(host string, port int) string {
	return JoinHostIntPort(bytes.S2b(host), port)
}

func JoinStrHostStrPort(host string, port string) string {
	return JoinHostPort(bytes.S2b(host), bytes.S2b(port))
}

func StatusOKBytes(major, minor int) []byte {
	buf := make([]byte, 0, okStatusTotalLen)
	buf = append(buf, okStartResponse...)
	buf = append(buf, '0'+byte(major), '.', '0'+byte(minor))
	buf = append(buf, okEndResponse...)
	return buf
}

func IsHTTPOK(buf []byte, isConnect bool) bool {
	if isConnect && len(buf) == 0 {
		return true
	}
	if len(buf) < maxResponseLen {
		return false
	}
	switch {
	case bytes.B2s(buf[:len(http10OKResponse)]) == http10OKResponse:
		return true
	case bytes.B2s(buf[:len(http11OKResponse)]) == http11OKResponse:
		return true
	case bytes.B2s(buf[:len(http2OKResponse)]) == http2OKResponse:
		return true
	default:
		return false
	}
}

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

func Flush(w any) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

var copyBufPool = &pool.LimitedPool{
	New: func() any {
		var b [copyBufSize]byte
		return &b
	},
	N: 1 << 6,
}

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

type HeaderDeleter interface {
	Del(string)
}

func DelHopHeaders(header HeaderDeleter) {
	for _, headerName := range HopHeaders {
		header.Del(headerName)
	}
}

func CopyHTTPHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst[k] = vv[:]
	}
}

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
	}
	req.Body.Close()

	return buf.Bytes(), nil
}

func IsConnClosedErr(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}

func PartitionIP4(ip string) (string, string, string, string) {
	part1, part234 := strings.Partition(ip, '.')
	part2, part34 := strings.Partition(part234, '.')
	part3, part4 := strings.Partition(part34, '.')
	return part1, part2, part3, part4
}
