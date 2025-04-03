package http

import (
	"encoding/base64"

	"github.com/sudosz/go-utils/bytes"
	"github.com/sudosz/go-utils/ints"
)

// BasicAuthHeader generates a Basic Auth header from string login and password.
// Optimization: Uses zero-copy S2b and efficient base64 encoding.
func BasicAuthHeader(login, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		bytes.S2b(login+":"+password),
	)
}

// JoinHostPort combines host bytes and an integer port into a string.
// Optimization: Efficient append with Int64ToString.
func JoinHostPort(host []byte, port int) string {
	b := append(host, ':')
	b = append(b, ints.Int64ToString(int64(port))...)
	return bytes.B2s(b)
}

// JoinHostPortString combines string host and integer port into a string.
// Optimization: Direct byte conversion avoids intermediate allocations.
func JoinHostPortString(host string, port int) string {
	return JoinHostPort([]byte(host), port)
}
