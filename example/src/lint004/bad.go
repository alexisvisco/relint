package lint004

import "context"

func Good(ctx context.Context, x int) {}

func Bad(x int, ctx context.Context) {} // want `LINT-004: context.Context must be the first parameter`

func AlsoBad(a, b int, ctx context.Context) {} // want `LINT-004: context.Context must be the first parameter`
