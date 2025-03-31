package http

import (
	"encoding/base64"

	"github.com/sudosz/go-utils/bytes"
	"github.com/sudosz/go-utils/ints"
)

func BasicAuthHeader(login, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		bytes.S2b(login + ":" + password),
	)
}

func JoinHostPort(host []byte, port int) string {
	b := append(host, ':')
	b = append(b, ints.Int64ToString(int64(port))...)
	return bytes.B2s(b)
}

func JoinHostPortString(host string, port int) string {
	return JoinHostPort([]byte(host), port)
}
