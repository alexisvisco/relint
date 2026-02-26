package fmt002

func Foo() {}

type Bar struct{} // want `FMT-002: declarations must be in order: type, const, var, func`
