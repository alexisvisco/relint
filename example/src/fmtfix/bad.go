package fmtfix

import "fmt"

// Bar is a type.
type Bar struct{} // want `FMTFIX: apply format fixes \(merge type blocks, reorder declarations\)`

// Baz is a type.
type Baz struct{}

// Foo does something.
func Foo() { fmt.Println() }
