package fmt004

type MyInterface interface {
	Foo()
	Bar() // want `FMT-004: interface methods must be separated by exactly one blank line`

	Baz()
}

type GoodInterface interface {
	Foo()

	Bar()

	Baz()
}

type GoodCommentedInterface interface {
	Foo()

	// Bar documents the next method and should not count as extra spacing.
	Bar()

	// Baz is also correctly separated by a single blank line.
	Baz()
}
