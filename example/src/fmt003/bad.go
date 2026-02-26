package fmt003

func StartBlank() { // want `FMT-003: function "StartBlank" body must not start with a blank line`

	x := 1
	_ = x
}

func EndBlank() {
	y := 2
	_ = y

} // want `FMT-003: function "EndBlank" body must not end with a blank line`
