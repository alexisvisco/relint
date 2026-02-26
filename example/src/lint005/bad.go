package lint005

func Good(a, b, c, d int) {}

func Bad(a, b, c, d, e int) {} // want `LINT-005: function "Bad" has 5 parameters, consider using a BadParams struct`

func AlsoBad(a int, b string, c bool, d float64, e int) {} // want `LINT-005: function "AlsoBad" has 5 parameters, consider using a AlsoBadParams struct`
