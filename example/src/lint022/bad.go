package handler

import "net/http"

type AuthHandler struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) { // want `LINT-022: route handler "Login" on "AuthHandler" must be in file "auth_login_handler\.go"`
}
