package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
	"github.com/solsw/httphelper"
)

// ReqBody sends an [http.Request] by [http.Client] and returns the contents of the response body.
//
// If 'isError' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is JSON-decoded from the response body.
// If response body does not contain JSON-encoded 'E', but contains a string,
// this string is returned in [httphelper.Error.Message] field.
//
// Pass [generichelper.NoType] as 'E' to omit processing of [httphelper.Error.Object].
func ReqBody[E any](client *http.Client, rq *http.Request, isError func(*http.Response) bool) ([]byte, error) {
	rs, err := client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if isError != nil && isError(rs) {
		// !generichelper.IsNoType[E]() - is checked in NewError
		herr, err := httphelper.NewError[E](rs, httphelper.ErrorOptionWithObject(), httphelper.ErrorOptionWithMessage())
		if err != nil {
			return nil, err
		}
		return nil, herr
	}
	return io.ReadAll(rs.Body)
}

// ReqJson sends an [http.Request] by [http.Client] and returns
// output object of type 'O' passed JSON-encoded as the response body.
//
// If 'isError' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is JSON-decoded from the response body.
// If response body does not contain JSON-encoded 'E', but contains a string,
// this string is returned in [httphelper.Error.Message] field.
//
// Pass [generichelper.NoType] as corresponding [type argument] to omit processing of either object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func ReqJson[O, E any](client *http.Client, rq *http.Request, isError func(*http.Response) bool) (*O, error) {
	bbout, err := ReqBody[E](client, rq, isError)
	if err != nil {
		return nil, err
	}
	if generichelper.IsNoType[O]() {
		return nil, nil
	}
	var out O
	if err := json.Unmarshal(bbout, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
