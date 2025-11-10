package middlewares

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareStac(mds ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mds) - 1; i >= 0; i-- {
			next = mds[i](next)
		}

		return next
	}
}
