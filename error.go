package httphelper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
)

// Error represents a HTTP error.
// Object (if turned on in options) is deserialized from JSON read from HTTP response body (if any).
type Error[T any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Object     T      `json:"object,omitempty"`
}

// Error implements the [error] interface.
//
// [error]: https://pkg.go.dev/builtin#error
func (e *Error[T]) Error() string {
	bb, _ := json.Marshal(e)
	return string(bb)
}

// NewError creates [Error] from [http.Response].
func NewError[T any](rs *http.Response, opts ...func(o *ErrorOptions)) (*Error[T], error) {
	herr := Error[T]{StatusCode: rs.StatusCode, Status: rs.Status}
	var options ErrorOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.withObject && !generichelper.IsNoType[T]() {
		bb, err := io.ReadAll(rs.Body)
		if err != nil {
			return nil, err
		}
		if len(bb) > 0 {
			if err := json.Unmarshal(bb, &herr.Object); err != nil {
				return nil, err
			}
		}
	}
	return &herr, nil
}
