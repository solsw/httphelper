package httphelper

// httpErrorOptions contains [HttpError] options.
type httpErrorOptions struct {
	withObject bool
}

// HttpErrorOption turns on some [HttpError] option.
type HttpErrorOption func() func(o *httpErrorOptions)

// WithObject turns on [HttpError] Object reading from a HTTP response body.
func WithObject() func(o *httpErrorOptions) {
	return func(o *httpErrorOptions) {
		o.withObject = true
	}
}
