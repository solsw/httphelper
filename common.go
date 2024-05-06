package httphelper

import (
	"errors"
	"net/http"
)

var ErrEmptyResponseBody = errors.New("empty response body")

// IsNotStatusOK determines whether r.StatusCode is not [http.StatusOK].
func IsNotStatusOK(rs *http.Response) bool {
	return rs.StatusCode != http.StatusOK
}
