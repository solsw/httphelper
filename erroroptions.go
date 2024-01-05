package httphelper

// ErrorOptions contains [Error] options.
type ErrorOptions struct {
	withObject bool
}

// ErrorOptionWithObject turns on [Error]'s Object JSON deserialization from HTTP response body.
func ErrorOptionWithObject() func(o *ErrorOptions) {
	return func(o *ErrorOptions) {
		o.withObject = true
	}
}
