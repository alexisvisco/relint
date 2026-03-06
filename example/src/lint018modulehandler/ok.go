package authhandler

import "net/http"

func MyMiddleware(next http.Handler) http.Handler {
	return next
}
