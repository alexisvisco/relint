package fmtfixcomments

import "fmt"

type ( // want `FMTFIX: apply format fixes \(merge type blocks, reorder declarations\)`
	Foo struct{}
	Bar struct{}
)

// KeepComments keeps body comments.
func KeepComments() {
	// keep: comment before statement
	fmt.Println("a") // keep: inline comment

	/* keep: block comment */
	fmt.Println("b")
}
