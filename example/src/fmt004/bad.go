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
