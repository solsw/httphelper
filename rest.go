package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/solsw/generichelper"
)

// RestInOut performs REST request-response sequence.
// Input object 'in' of type 'I' is JSON serialized to request body.
// Output object of type 'O' is JSON deserialized from response body, if there is no error, and returned.
// If response's HTTP status code is not [http.StatusOK], [HttpError] is returned.
// [HttpError]'s Object of type 'E' is JSON deserialized from response body.
// Pass [generichelper.NoType] as corresponding [type argument] to skip processing of any object.
//
// [type argument]: https://go.dev/ref/spec#Instantiations
func RestInOut[I, O, E any](ctx context.Context, client *http.Client, method, url string, header http.Header, in *I) (*O, error) {
	var body io.Reader
	if !generichelper.IsNoType[I]() {
		bbIn, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bbIn)
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
		herr, err := NewHttpError[E](resp, WithObject())
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
