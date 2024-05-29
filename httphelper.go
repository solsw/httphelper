package httphelper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

// JsonReader creates [bytes.Reader] containing JSON-encoded 'in'
// to use as body with [http.NewRequest] or [http.NewRequestWithContext].
// If 'in' is nil, nil is returned.
func JsonReader(in any) (*bytes.Reader, error) {
	if in == nil {
		return nil, nil
	}
	jin, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jin), nil
}

// JsonBody creates [http] Body containing JSON-encoded 'in'.
// If 'in' is nil, [http.NoBody] is returned.
func JsonBody(in any) (io.ReadCloser, error) {
	if in == nil {
		return http.NoBody, nil
	}
	body, err := JsonReader(in)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(body), nil
}

// AuthBasic returns Basic authorization value.
func AuthBasic(userid, password string) string {
	// https://datatracker.ietf.org/doc/html/rfc7617
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userid+":"+password))
}
