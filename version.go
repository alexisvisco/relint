package main

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var embeddedVersion string

func binaryVersion() string {
	version := strings.TrimSpace(embeddedVersion)
	if version == "" {
		return "dev"
	}
	return version
}
