package httphelper

// HttpErrorOptions contains [HttpError] options.
type HttpErrorOptions struct {
	withObject bool
}

// WithObject turns on [HttpError]'s Object JSON deserializing from HTTP response body.
func WithObject() func(o *HttpErrorOptions) {
	return func(o *HttpErrorOptions) {
		o.withObject = true
	}
}
