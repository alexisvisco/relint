package fmtfix

import "fmt"

// Bar is a type.
type Bar struct{} // want `FMTFIX: apply format fixes \(merge declaration blocks, reorder declarations\)`

// Baz is a type.
type Baz struct{}

const Alpha = "a"

const Beta = "b"

var One = 1

var Two = 2

// Foo does something.
func Foo() { fmt.Println() }
