package handler

import "net/http"

type AuthHandler struct{}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {}

// LoginInput must now live in handlertypes.
type LoginInput struct{} // want `LINT-023: type "LoginInput" must be declared in package "handlertypes"`
