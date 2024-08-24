package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ApplyToFunc(h http.HandlerFunc, middlewares ...Middleware) http.Handler {
	return Stack(middlewares...)(h)
}

func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	return Stack(middlewares...)(h)
}

func Stack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			top := middlewares[i]
			next = top(next)
		}

		return next
	}
}
