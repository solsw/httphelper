package httphelper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// JsonBody sets [http.Request]'s body to JSON-encoded 'in', if not nil.
func JsonBody(rq *http.Request, in any) (*http.Request, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}
	jin, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	rq.Body = io.NopCloser(bytes.NewReader(jin))
	return rq, nil
}

// AuthBasic returns Basic authorization value.
func AuthBasic(userid, password string) string {
	// https://datatracker.ietf.org/doc/html/rfc7617
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userid+":"+password))
}
