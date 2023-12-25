package httphelper

import (
	"encoding/json"
	"io"
	"net/http"
)

// HttpError represents a HTTP error.
// Object (if turned on) is read from a HTTP response body, if any.
type HttpError[O any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Object     O      `json:"object,omitempty"`
}

// Error implements the [error] interface.
func (e *HttpError[B]) Error() string {
	bb, _ := json.MarshalIndent(e, "", "  ")
	return string(bb)
}

// NewHttpError creates HttpError from [http.Response].
func NewHttpError[O any](rs *http.Response, opts ...func(o *HttpErrorOptions)) (*HttpError[O], error) {
	he := HttpError[O]{StatusCode: rs.StatusCode, Status: rs.Status}
	var options HttpErrorOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.withObject {
		bb, _ := io.ReadAll(rs.Body)
		var obj O
		if len(bb) > 0 {
			if err := json.Unmarshal(bb, &obj); err != nil {
				return nil, err
			}
		}
		he.Object = obj
	}
	return &he, nil
}
