package lint018

import "net/http"

func MyMiddleware(next http.Handler) http.Handler { // want `LINT-018: middleware function "MyMiddleware" outside handler package must be named "Middleware"`
	return next
}

func Middleware(next http.Handler) http.Handler { // ok
	return next
}
