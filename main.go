package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/alexisvisco/relint/all"
)

func main() {
	multichecker.Main(all.Analyzers...)
}
