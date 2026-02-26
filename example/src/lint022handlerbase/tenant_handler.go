package handler

type TenantHandler struct{}

func (h *TenantHandler) Tenant() {} // ok - route name equals handler base name, so tenant_handler.go is valid
