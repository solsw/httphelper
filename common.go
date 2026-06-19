package httphelper

import (
	"errors"
	"io"
	"net/http"
)

// MaxResponseBodySize limits the number of bytes read from an HTTP response body.
// Reading a body larger than this returns [ErrResponseBodyTooLarge].
var MaxResponseBodySize int64 = 10 << 20 // 10 MiB

// ErrResponseBodyTooLarge is returned when a response body exceeds [MaxResponseBodySize].
var ErrResponseBodyTooLarge = errors.New("response body too large")

// readBody reads up to [MaxResponseBodySize] bytes from r,
// returning [ErrResponseBodyTooLarge] if the limit is exceeded.
func readBody(r io.Reader) ([]byte, error) {
	bb, err := io.ReadAll(io.LimitReader(r, MaxResponseBodySize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(bb)) > MaxResponseBodySize {
		return nil, ErrResponseBodyTooLarge
	}
	return bb, nil
}

// IsNotStatusOK determines whether rs.StatusCode is not [http.StatusOK].
func IsNotStatusOK(rs *http.Response) bool {
	return rs.StatusCode != http.StatusOK
}

// IsNotStatus2xx determines whether rs.StatusCode is not a 2xx successful one.
func IsNotStatus2xx(rs *http.Response) bool {
	return !(http.StatusOK <= rs.StatusCode && rs.StatusCode < http.StatusMultipleChoices)
}
