# httphelper
[![Go Reference](https://pkg.go.dev/badge/github.com/solsw/httphelper.svg)](https://pkg.go.dev/github.com/solsw/httphelper)
[![GitHub](https://img.shields.io/badge/github--green?logo=github)](https://github.com/solsw/httphelper)

Helpers for Go's [`net/http`](https://pkg.go.dev/net/http) package: generic request senders, a
typed HTTP error, status-code checks, and small request-building utilities.

## Install

```sh
go get github.com/solsw/httphelper
```

```go
import "github.com/solsw/httphelper"
```

## Overview

| Symbol | Purpose |
| --- | --- |
| `ReqJSON[O, E]` | Send a request and JSON-decode the response body into `*O`. |
| `ReqD9e[O, E]` | Like `ReqJSON`, but with a caller-supplied deserializer. |
| `ReqBody[E]` | Send a request and return the raw response body. |
| `Error[T]` | A typed HTTP error carrying status, an optional decoded object, and/or a message. |
| `NewError[T]` | Build an `Error[T]` from an `*http.Response`. |
| `IsNotStatusOK` / `IsNotStatus2xx` | Status-code predicates, handy as the `isErr` argument. |
| `JSONReader` | Build a request body `io.Reader` from a value via JSON. |
| `AuthBasic` | Build an HTTP Basic `Authorization` header value. |

The two generic type parameters that recur throughout are:

- **`O`** — the type to decode a **successful** response body into.
- **`E`** — the type to decode an **error** response body into (stored in `Error.Object`).

To skip processing of either, pass [`generichelper.NoType`](https://pkg.go.dev/github.com/solsw/generichelper#NoType)
as that type argument.

## Sending requests

All senders take an `*http.Client`, an `*http.Request`, and an `isErr func(*http.Response) bool`.
When `isErr` is non-nil and reports `true`, the sender returns an [`*Error[E]`](#error-handling)
instead of decoding a success value.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/solsw/generichelper"
	"github.com/solsw/httphelper"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// APIError is the shape this API returns for error responses.
type APIError struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
}

func getUser(ctx context.Context) (*User, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.example.com/user/1", nil)
	if err != nil {
		return nil, err
	}
	// Decode success bodies into *User, error bodies into APIError.
	return httphelper.ReqJSON[User, APIError](http.DefaultClient, rq, httphelper.IsNotStatus2xx)
}

func main() {
	u, err := getUser(context.Background())
	if err != nil {
		var herr *httphelper.Error[APIError]
		if errors.As(err, &herr) {
			fmt.Printf("HTTP %d: %+v\n", herr.StatusCode, herr.Object)
		}
		return
	}
	fmt.Println(u.Name)
}
```

### Sending a JSON request body

Use `JSONReader` to turn a value into a request body:

```go
body, err := httphelper.JSONReader(User{Name: "Alice", Age: 30})
if err != nil {
	return err
}
rq, err := http.NewRequest(http.MethodPost, url, body)
```

`JSONReader(nil)` returns `(nil, nil)`, so an absent body is handled naturally.

### Custom deserialization

`ReqD9e` is the general form — supply any `func(data []byte, v any) error`
(e.g. a YAML or protobuf-JSON decoder). `ReqJSON` is simply `ReqD9e` with `json.Unmarshal`.

```go
out, err := httphelper.ReqD9e[Config, generichelper.NoType](
	client, rq, httphelper.IsNotStatus2xx, yaml.Unmarshal)
```

### Raw body

`ReqBody` returns the response body bytes (still applying `isErr`):

```go
bb, err := httphelper.ReqBody[generichelper.NoType](client, rq, httphelper.IsNotStatus2xx)
```

## Error handling

`Error[T]` implements the `error` interface; its `Error()` method returns the value JSON-encoded.

```go
type Error[T any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Object     T      `json:"object,omitempty"`
	Message    string `json:"message,omitempty"`
}
```

When a sender produces an error response, it builds the `Error` with both object and message
processing enabled. The body is interpreted as follows:

- **Decodes cleanly as `T`** (strict: unknown fields and trailing data are rejected) → stored in `Object`.
- **Does not decode as `T`** but is valid UTF-8 text → stored in `Message`.
- **Empty body** → neither is set, but `StatusCode`/`Status` are still populated.

You can also build an error directly from a response with explicit options:

```go
herr, err := httphelper.NewError[APIError](resp,
	httphelper.ErrorOptionWithObject(),
	httphelper.ErrorOptionWithMessage())
```

With no options, `NewError` does not read the body and returns an `Error` carrying only
`StatusCode` and `Status`.

## Status checks

```go
httphelper.IsNotStatusOK(resp)  // resp.StatusCode != 200
httphelper.IsNotStatus2xx(resp) // resp.StatusCode is outside 200–299
```

Both have the `func(*http.Response) bool` signature expected by the `isErr` parameter.

## Basic authentication

```go
rq.Header.Set("Authorization", httphelper.AuthBasic("Aladdin", "open sesame"))
// Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
```

## Limiting response size

To guard against unbounded or hostile responses, bodies are read through a size limit.
Reading more than `MaxResponseBodySize` bytes returns `ErrResponseBodyTooLarge`.

```go
httphelper.MaxResponseBodySize = 1 << 20 // 1 MiB; default is 10 MiB
```
