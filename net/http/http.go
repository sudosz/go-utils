package httputils

import (
	"encoding/base64"

	bytesutils "github.com/sudosz/go-utils/bytes"
	intutils "github.com/sudosz/go-utils/ints"
)

// BasicAuthHeader generates a Basic Auth header from string login and password.
// Optimization: Uses zero-copy S2b and efficient base64 encoding.
func BasicAuthHeader(login, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		bytesutils.S2b(login+":"+password),
	)
}

// JoinHostPort combines host bytes and an integer port into a string.
// Optimization: Efficient append with Int64ToString.
func JoinHostPort(host []byte, port int) string {
	b := append(host, ':')
	b = append(b, intutils.Int64ToString(int64(port))...)
	return bytesutils.B2s(b)
}

// JoinHostPortString combines string host and integer port into a string.
// Optimization: Direct byte conversion avoids intermediate allocations.
func JoinHostPortString(host string, port int) string {
	return JoinHostPort([]byte(host), port)
}
