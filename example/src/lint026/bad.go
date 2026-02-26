package handler

type InvitationTokenBodyOutput struct {
	Client ClientOutput
}

type ClientOutput struct { // want `LINT-026: body-only struct "ClientOutput" must be prefixed with "InvitationTokenBody" and suffixed with "Output"`
	ID string
}
