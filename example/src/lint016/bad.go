package handler

import "net/http"

func InjectUser(next http.Handler) http.Handler { // want `LINT-016: middleware "InjectUser" must be in file "inject_user.go"`
	return next
}

type humaContext interface{}

type humaNext func(humaContext)

func injectUserSession(_ interface{}) func(humaContext, humaNext) { // want `LINT-016: middleware "injectUserSession" must be in file "inject_user_session.go"`
	return func(_ humaContext, _ humaNext) {}
}
