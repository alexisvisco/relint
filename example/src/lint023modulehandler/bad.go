package authhandler

import "net/http"

type AuthHandler struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {}

type LoginInput struct{} // want `LINT-023: type "LoginInput" must be declared in route file "login\.go"`
