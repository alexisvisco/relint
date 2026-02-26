package fmt005

type (
	TenantInput  struct{}
	TenantOutput struct { // want `FMT-005: type specs in a type block must be separated by exactly one blank line`
		Body string
	}
)
