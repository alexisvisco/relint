package handler // want `LINT-015: file "bad.go" in store/service/handler package must contain exactly one exported store/service/handler method, found 2`

type AuthHandler struct{}

func (h *AuthHandler) Login() {}

func (h *AuthHandler) Logout() {}

func Helper() {} // ok - exported function but not a store/service/handler method
