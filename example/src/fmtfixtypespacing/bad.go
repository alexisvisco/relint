package fmtfixtypespacing

type ( // want `FMTFIX: apply format fixes \(merge declaration blocks, reorder declarations\)`
	TenantInput  struct{}
	TenantOutput struct {
		Body string
	}
)
