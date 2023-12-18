package httphelper

import (
	"encoding/json"
	"io"
	"net/http"
)

// HttpError represents an HTTP error.
// ErrorBody is read from an HTTP response body.
type HttpError[B any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	ErrorBody  B      `json:"error_body,omitempty"`
}

// Error implements the [error] interface.
func (e *HttpError[B]) Error() string {
	bb, _ := json.MarshalIndent(e, "", "  ")
	return string(bb)
}

// NewHttpError creates HttpError from [http.Response].
func NewHttpError[B any](rs *http.Response) (*HttpError[B], error) {
	bb, _ := io.ReadAll(rs.Body)
	var b B
	if err := json.Unmarshal(bb, &b); err != nil {
		return nil, err
	}
	return &HttpError[B]{
			StatusCode: rs.StatusCode,
			Status:     rs.Status,
			ErrorBody:  b,
		},
		nil
}
