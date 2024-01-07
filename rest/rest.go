package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/solsw/generichelper"
	"github.com/solsw/httphelper"
)

// InOut performs REST request-response sequence with input and output objects
// passed JSON-encoded as request and response body respectively.
// 'I' - type of input object. 'O' - type of output object.
// If there is no error, output object is returned.
// If response's HTTP status code is not [http.StatusOK], [httphelper.Error] is returned.
// [httphelper.Error]'s Object of type 'E' is JSON-decoded from response body.
// Pass [generichelper.NoType] as corresponding [type argument] to skip processing of any object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func InOut[I, O, E any](ctx context.Context, client *http.Client, method, url string, header http.Header, in *I) (*O, error) {
	var body io.Reader
	if in != nil && !generichelper.IsNoType[I]() {
		if generichelper.TypeOf[I]().Kind() == reflect.String {
			var str any = *in
			body = strings.NewReader(str.(string))
		} else {
			bbIn, err := json.Marshal(in)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(bbIn)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// !generichelper.IsNoType[E]() - is checked in NewError
		herr, err := httphelper.NewError[E](resp, httphelper.ErrorOptionWithObject())
		if err != nil {
			return nil, err
		}
		return nil, herr
	}
	bbOut, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !generichelper.IsNoType[O]() {
		var out O
		if err := json.Unmarshal(bbOut, &out); err != nil {
			return nil, err
		}
		return &out, nil
	}
	return nil, nil
}
