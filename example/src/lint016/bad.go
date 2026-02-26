package handler

import "net/http"

func InjectUser(next http.Handler) http.Handler { // want `LINT-016: middleware "InjectUser" must be in file "inject_user.go"`
	return next
}
