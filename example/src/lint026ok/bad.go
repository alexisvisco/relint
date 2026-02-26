package handler

type InvitationTokenBodyOutput struct {
	Client InvitationTokenBodyClientOutput
}

type InvitationTokenBodyClientOutput struct {
	ID string
}

type TenantBodyInput struct {
	Item TenantBodyItem
}

type TenantBodyItem struct{} // ok for LINT-026: used outside body struct, not body-only

func UseTenantBodyItem(_ TenantBodyItem) {}
