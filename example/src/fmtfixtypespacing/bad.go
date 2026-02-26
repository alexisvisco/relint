package fmtfixtypespacing

type ( // want `FMTFIX: apply format fixes \(merge type blocks, reorder declarations\)`
	TenantInput  struct{}
	TenantOutput struct {
		Body string
	}
)
