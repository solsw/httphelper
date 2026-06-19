package httphelper

import (
	"bytes"
	"encoding/json"
	"errors"
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
	bb, err := json.Marshal(e)
	if err != nil {
		return e.Status
	}
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
		bb, err := readBody(rs.Body)
		if err != nil {
			return nil, err
		}
		if len(bb) == 0 {
			// no body to extract Object/Message from, but StatusCode/Status are still meaningful
			return &herr, nil
		}
		if err := objMsg(&herr, bb, options); err != nil {
			return nil, err
		}
	}
	return &herr, nil
}

func objMsg[T any](herr *Error[T], bb []byte, options ErrorOptions) error {
	gotObject := false
	var erro error
	if options.withObject && !generichelper.IsNoType[T]() {
		dec := json.NewDecoder(bytes.NewReader(bb))
		dec.DisallowUnknownFields()
		if erro = dec.Decode(&herr.Object); erro == nil {
			if dec.More() {
				// trailing data after the JSON value: not a clean 'T'
				erro = errors.New("unexpected trailing data after JSON value")
			} else {
				gotObject = true
			}
		}
	}
	if options.withMessage && !gotObject {
		if !utf8.Valid(bb) {
			return errors.Join(erro, errors.New("invalid UTF-8-encoded runes"))
		}
		herr.Message = string(bb)
		return nil
	}
	return erro
}
