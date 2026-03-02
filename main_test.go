package main

import (
	"slices"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestPreprocessArgs_OnlyFmtfix(t *testing.T) {
	analyzers := []*analysis.Analyzer{
		{Name: "fmtfix"},
		{Name: "lint001"},
		{Name: "lint027"},
	}

	args, err := preprocessArgs([]string{"relint", "-only-fmtfix", "./..."}, analyzers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Contains(args, "-fmtfix=true") {
		t.Fatalf("expected -fmtfix=true in args, got: %v", args)
	}
	if !slices.Contains(args, "-lint001=false") {
		t.Fatalf("expected -lint001=false in args, got: %v", args)
	}
	if !slices.Contains(args, "-lint027=false") {
		t.Fatalf("expected -lint027=false in args, got: %v", args)
	}
	targetIdx := slices.Index(args, "./...")
	if targetIdx == -1 {
		t.Fatalf("expected package arg in output args, got: %v", args)
	}
	if targetIdx == 1 {
		t.Fatalf("expected injected flags before package arg, got: %v", args)
	}
	if slices.Index(args, "-fmtfix=true") > targetIdx {
		t.Fatalf("expected -fmtfix=true before package arg, got: %v", args)
	}
}

func TestPreprocessArgs_OnlyFmtfixFalse(t *testing.T) {
	analyzers := []*analysis.Analyzer{
		{Name: "fmtfix"},
		{Name: "lint001"},
	}

	args, err := preprocessArgs([]string{"relint", "-only-fmtfix=false", "./..."}, analyzers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if slices.Contains(args, "-only-fmtfix=false") {
		t.Fatalf("custom flag should be removed from args: %v", args)
	}
	if slices.Contains(args, "-lint001=false") {
		t.Fatalf("lint001 should not be disabled when only-fmtfix=false: %v", args)
	}
}

func TestPreprocessArgs_InvalidOnlyFmtfixValue(t *testing.T) {
	analyzers := []*analysis.Analyzer{
		{Name: "fmtfix"},
	}

	_, err := preprocessArgs([]string{"relint", "-only-fmtfix=nope", "./..."}, analyzers)
	if err == nil {
		t.Fatal("expected error for invalid -only-fmtfix value")
	}
}

func TestStripVersionArgs_Enabled(t *testing.T) {
	showVersion, args, err := stripVersionArgs([]string{"relint", "-version", "./..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !showVersion {
		t.Fatal("expected version flag to be enabled")
	}
	if slices.Contains(args, "-version") {
		t.Fatalf("version flag should be removed from args: %v", args)
	}
	if !slices.Contains(args, "./...") {
		t.Fatalf("expected package args to be preserved: %v", args)
	}
}

func TestStripVersionArgs_Disabled(t *testing.T) {
	showVersion, args, err := stripVersionArgs([]string{"relint", "--version=false", "./..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if showVersion {
		t.Fatal("expected version flag to be disabled")
	}
	if slices.Contains(args, "--version=false") {
		t.Fatalf("version flag should be removed from args: %v", args)
	}
	if !slices.Contains(args, "./...") {
		t.Fatalf("expected package args to be preserved: %v", args)
	}
}

func TestStripVersionArgs_InvalidValue(t *testing.T) {
	_, _, err := stripVersionArgs([]string{"relint", "-version=nope"})
	if err == nil {
		t.Fatal("expected error for invalid -version value")
	}
}
