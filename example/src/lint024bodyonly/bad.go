package handler

type InvitationTokenBodyOutput struct {
	Client InvitationTokenBodyClientOutput
}

// ok for LINT-024: body-only helper type, validated by LINT-026.
type InvitationTokenBodyClientOutput struct {
	ID string
}
