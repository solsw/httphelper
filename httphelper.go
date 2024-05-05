package httphelper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// JsonBody creates [http.Request.Body] with JSON-encoded 'in', if not nil.
func JsonBody(in any) (io.ReadCloser, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}
	jin, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(jin)), nil
}

// AuthBasic returns Basic authorization value.
func AuthBasic(userid, password string) string {
	// https://datatracker.ietf.org/doc/html/rfc7617
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userid+":"+password))
}
