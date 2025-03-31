package httphelper

import (
	"errors"
	"net/http"
)

var ErrEmptyResponseBody = errors.New("empty response body")

// IsNotStatusOK determines whether rs.StatusCode is not [http.StatusOK].
func IsNotStatusOK(rs *http.Response) bool {
	return rs.StatusCode != http.StatusOK
}

// IsNotStatus2xx determines whether rs.StatusCode is not a 2xx succesful one.
func IsNotStatus2xx(rs *http.Response) bool {
	return !(http.StatusOK <= rs.StatusCode && rs.StatusCode < http.StatusMultipleChoices)
}
