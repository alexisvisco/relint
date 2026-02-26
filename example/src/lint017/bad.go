package handler

import "net/http"

func RequireAuth(next http.Handler) http.Handler { // want `LINT-017: middleware "RequireAuth" must be in file "require_auth.go"`
	return next
}

type humaContext interface{}

type humaNext func(humaContext)

func requireUserSession(_ interface{}) func(humaContext, humaNext) { // want `LINT-017: middleware "requireUserSession" must be in file "require_user_session.go"`
	return func(_ humaContext, _ humaNext) {}
}
