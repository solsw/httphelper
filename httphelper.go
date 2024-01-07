package httphelper

import (
	"encoding/base64"
)

// AuthBasic returns Basic authorization value.
func AuthBasic(userid, password string) string {
	// https://datatracker.ietf.org/doc/html/rfc7617
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userid+":"+password))
}
