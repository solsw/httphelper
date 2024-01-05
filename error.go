package httphelper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
)

// Error represents a HTTP error.
// Object (if turned on in options) is deserialized from JSON read from HTTP response body (if any).
type Error[E any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Object     E      `json:"object,omitempty"`
}

// Error implements the [error] interface.
func (e *Error[B]) Error() string {
	bb, _ := json.MarshalIndent(e, "", "  ")
	return string(bb)
}

// NewError creates [Error] from [http.Response].
func NewError[E any](rs *http.Response, opts ...func(o *ErrorOptions)) (*Error[E], error) {
	herr := Error[E]{StatusCode: rs.StatusCode, Status: rs.Status}
	var options ErrorOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.withObject && !generichelper.IsNoType[E]() {
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
