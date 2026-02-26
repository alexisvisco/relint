package fmtfix

import "fmt"

// Foo does something.
func Foo() { fmt.Println() } // want `FMTFIX: apply format fixes \(merge type blocks, reorder declarations\)`

// Bar is a type.
type Bar struct{}

// Baz is a type.
type Baz struct{}
