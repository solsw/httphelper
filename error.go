package httphelper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/solsw/generichelper"
)

// Error represents a HTTP error.
// Object (if turned on in options) is deserialized from JSON read from HTTP response body (if any).
type Error[T any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Object     T      `json:"object,omitempty"`
	Message    string `json:"message,omitempty"`
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
	if options.withObject || options.withMessage {
		bb, err := io.ReadAll(rs.Body)
		if err != nil {
			return nil, err
		}
		if len(bb) == 0 {
			return nil, ErrEmptyResponseBody
		}
		_, err = objMsg(&herr, bb, options)
		if err != nil {
			return nil, err
		}
	}
	return &herr, nil
}

func objMsg[T any](herr *Error[T], bb []byte, options ErrorOptions) (*Error[T], error) {
	var erro, errm error
	if options.withObject && !generichelper.IsNoType[T]() {
		erro = json.Unmarshal(bb, &herr.Object)
	}
	if options.withMessage && (erro != nil || generichelper.IsZeroValue(herr.Object)) {
		if !utf8.Valid(bb) {
			errm = errors.New("invalid UTF-8-encoded runes")
		} else {
			herr.Message = string(bb)
			return herr, nil
		}
	}
	return herr, errors.Join(erro, errm)
}
