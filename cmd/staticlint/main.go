package main

import (
	"strings"

	"github.com/ZnNr/go-musthave-metrics.git/internal/staticlint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// This package defines the main function for an analysis driver
// with several analyzers from packages go/analysis and staticcheck.io
// Usage:
// `go run cmd/staticlint/main.go <analyzers> <files>`

// Analyzers from package golang.org/x/tools/go/analysis/passes:
// printf - enable printf analysis
// structtag - enable structtag analysis
// shift - enable shift analysis

// Analyzers from package staticcheck.io:
// SA* - see docs (https://staticcheck.dev/docs/checks/#SA)
// S1039 - check unnecessary use of fmt.Sprint
// ST1008 - check if a function’s error value is the last return value
// QF1003 - check if the if/else-if chain needs to be converted to tagged switch

// Custom analyzers
// exitcheck - check for os.Exit from main functions of package main

// Examples:
// Checking all files in the current folder with all analyzers
// go run cmd/staticlint/main.go ./...
func main() {
	// the standard static analyzers of package golang.org/x/tools/go/analysis/passes
	// and a custom analyzer which denies the usage of os.Exit in main function of package main
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		staticlint.ExitFromMainAnalyzer,
	}
	someAdditionalChecks := map[string]struct{}{
		"S1039":  {}, // Unnecessary use of fmt.Sprint
		"ST1008": {}, // A function’s error value should be its last return value
		"QF1003": {}, // Convert if/else-if chain to tagged switch
	}

	// all analyzers SA class of package staticcheck.io
	// and some additional analyzers of other classes from package staticcheck.io
	for _, v := range staticcheck.Analyzers {
		if _, ok := someAdditionalChecks[v.Analyzer.Name]; ok || strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}
	}
	multichecker.Main(
		checks...,
	)
}
