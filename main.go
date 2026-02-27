package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/alexisvisco/relint/all"
)

func main() {
	args, err := preprocessArgs(os.Args, all.Analyzers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
	os.Args = args
	multichecker.Main(all.Analyzers...)
}

func preprocessArgs(args []string, analyzers []*analysis.Analyzer) ([]string, error) {
	if len(args) == 0 {
		return args, nil
	}

	onlyFmtfix := false
	filtered := make([]string, 0, len(args))
	filtered = append(filtered, args[0])

	for _, arg := range args[1:] {
		handled, enabled, err := parseOnlyFmtfixArg(arg)
		if err != nil {
			return nil, err
		}
		if handled {
			onlyFmtfix = enabled
			continue
		}
		filtered = append(filtered, arg)
	}

	if !onlyFmtfix {
		return filtered, nil
	}

	// Keep fmtfix enabled and disable all other analyzers.
	filtered = append(filtered, "-fmtfix=true")
	for _, analyzer := range analyzers {
		if analyzer.Name == "fmtfix" {
			continue
		}
		filtered = append(filtered, "-"+analyzer.Name+"=false")
	}

	return filtered, nil
}

func parseOnlyFmtfixArg(arg string) (handled bool, enabled bool, err error) {
	if arg == "-only-fmtfix" || arg == "--only-fmtfix" {
		return true, true, nil
	}

	const shortPrefix = "-only-fmtfix="
	const longPrefix = "--only-fmtfix="
	if strings.HasPrefix(arg, shortPrefix) || strings.HasPrefix(arg, longPrefix) {
		value := strings.TrimPrefix(strings.TrimPrefix(arg, shortPrefix), longPrefix)
		b, parseErr := strconv.ParseBool(value)
		if parseErr != nil {
			return false, false, fmt.Errorf("invalid value for -only-fmtfix: %q", value)
		}
		return true, b, nil
	}

	return false, false, nil
}
