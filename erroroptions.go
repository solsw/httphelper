package httphelper

// ErrorOptions contains [Error] options.
type ErrorOptions struct {
	withObject  bool
	withMessage bool
}

// ErrorOptionWithObject turns on [Error]'s Object JSON deserialization from HTTP response body.
func ErrorOptionWithObject() func(o *ErrorOptions) {
	return func(o *ErrorOptions) {
		o.withObject = true
	}
}

// ErrorOptionWithMessage turns on [Error]'s Message string deserialization from HTTP response body.
func ErrorOptionWithMessage() func(o *ErrorOptions) {
	return func(o *ErrorOptions) {
		o.withMessage = true
	}
}
