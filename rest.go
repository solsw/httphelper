package httphelper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
)

// ReqBody sends the provided [http.Request] by the provided [http.Client]
// and returns the contents of the response body.
//
// If 'isErr' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is JSON-decoded from the response body.
// If response body does not contain JSON-encoded 'E', but contains a string,
// this string is returned in [httphelper.Error.Message] field.
//
// Pass [generichelper.NoType] as 'E' to omit processing of [httphelper.Error.Object].
func ReqBody[E any](cl *http.Client, rq *http.Request, isErr func(*http.Response) bool) ([]byte, error) {
	rs, err := cl.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if isErr != nil && isErr(rs) {
		// !generichelper.IsNoType[E]() - is checked in NewError
		herr, err := NewError[E](rs, ErrorOptionWithObject(), ErrorOptionWithMessage())
		if err != nil {
			return nil, err
		}
		return nil, herr
	}
	return io.ReadAll(rs.Body)
}

// ReqD9e sends the provided [http.Request] by the provided [http.Client] and returns
// output object of type 'O' deserialized by 'd9e' from the response body.
//
// If 'isErr' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is JSON-decoded from the response body.
// If response body does not contain JSON-encoded 'E', but contains a string,
// this string is returned in [httphelper.Error.Message] field.
//
// Pass [generichelper.NoType] as corresponding [type argument] to omit processing of either object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func ReqD9e[O, E any](cl *http.Client, rq *http.Request, isErr func(*http.Response) bool,
	d9e func(data []byte, v any) error) (*O, error) {
	body, err := ReqBody[E](cl, rq, isErr)
	if err != nil {
		return nil, err
	}
	if generichelper.IsNoType[O]() {
		return nil, nil
	}
	var o O
	if err := d9e(body, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// ReqJson sends the provided [http.Request] by the provided [http.Client] and returns
// output object of type 'O' deserialized by [json.Unmarshal] from the response body.
//
// If 'isErr' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is JSON-decoded from the response body.
// If response body does not contain JSON-encoded 'E', but contains a string,
// this string is returned in [httphelper.Error.Message] field.
//
// Pass [generichelper.NoType] as corresponding [type argument] to omit processing of either object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func ReqJson[O, E any](cl *http.Client, rq *http.Request, isErr func(*http.Response) bool) (*O, error) {
	return ReqD9e[O, E](cl, rq, isErr, json.Unmarshal)
}
