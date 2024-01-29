package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
	"github.com/solsw/httphelper"
)

// BodyBody performs REST request-response sequence with 'in' request body
// and returns the contents of the response body.
//
// If 'isError' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is then JSON-decoded from the response body.
// Pass [generichelper.NoType] as 'E' to omit processing of [httphelper.Error.Object].
func BodyBody[E any](ctx context.Context, client *http.Client, method, url string,
	header http.Header, in io.Reader, isError func(*http.Response) bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, in)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if isError != nil && isError(resp) {
		// !generichelper.IsNoType[E]() - is checked in NewError
		herr, err := httphelper.NewError[E](resp, httphelper.ErrorOptionWithObject())
		if err != nil {
			return nil, err
		}
		return nil, herr
	}
	return io.ReadAll(resp.Body)
}

// BodyJson performs REST request-response sequence with 'in' request body
// and output object of type 'O' passed JSON-encoded as the response body.
//
// If 'isError' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is then JSON-decoded from the response body.
//
// Pass [generichelper.NoType] as corresponding [type argument] to omit processing of either object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func BodyJson[O, E any](ctx context.Context, client *http.Client, method, url string,
	header http.Header, in io.Reader, isError func(*http.Response) bool) (*O, error) {
	bbOut, err := BodyBody[E](ctx, client, method, url, header, in, isError)
	if err != nil {
		return nil, err
	}
	if generichelper.IsNoType[O]() {
		return nil, nil
	}
	var out O
	if err := json.Unmarshal(bbOut, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// JsonJson performs REST request-response sequence with input and output objects
// passed JSON-encoded as the request and the response body respectively.
// 'I' - type of the input object. 'O' - type of the output object.
//
// If 'isError' is not nil and returns 'true', [httphelper.Error] is returned.
// [httphelper.Error.Object] of type 'E' is then JSON-decoded from the response body.
//
// Pass [generichelper.NoType] as corresponding [type argument] to omit processing of either object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func JsonJson[I, O, E any](ctx context.Context, client *http.Client, method, url string,
	header http.Header, in *I, isError func(*http.Response) bool) (*O, error) {
	var body io.Reader
	if in != nil && !generichelper.IsNoType[I]() {
		bbIn, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bbIn)
	}
	return BodyJson[O, E](ctx, client, method, url, header, body, isError)
}
