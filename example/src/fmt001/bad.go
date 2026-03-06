package fmt001

type Foo struct{}

type Bar struct{} // want `FMT-001: consecutive declarations should be merged into a declaration block`

const A = "a"

const B = "b" // want `FMT-001: consecutive declarations should be merged into a declaration block`

var X = 1

var Y = 2 // want `FMT-001: consecutive declarations should be merged into a declaration block`
