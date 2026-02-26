package handler

import "net/http"

type AuthHandler struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {}

// LoginInput is declared in bad.go instead of auth_login_handler.go
type LoginInput struct{} // want `LINT-023: type "LoginInput" must be declared in route file "auth_login_handler\.go"`
