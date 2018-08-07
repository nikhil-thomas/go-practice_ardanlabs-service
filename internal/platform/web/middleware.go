package web

// Middleware defines a middleware type which wraps a handler
type Middleware func(Handler) Handler

// wrapMiddleware wraps a handler with some middleware
func wrapMiddleware(handler Handler, mw []Middleware) Handler {

	for i := len(mw) - 1; i > 0; i-- {
		if mw[i] != nil {
			handler = mw[i](handler)
		}
	}
	return handler
}
