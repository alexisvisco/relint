package authhandler

import "net/http"

type AuthHandler struct{}

type LoginInput struct{}

type LoginOutput struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {}
