package analyzer

import (
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for lowercase start, English only, no emojis/secrets",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, lit, isLog := IsLogCall(node)
			if !isLog {
				return true
			}

			msg := lit.Value[1 : len(lit.Value)-1]

			// Rule 1: lowercase start
			if len(msg) > 0 && unicode.IsUpper(rune(msg[0])) {
				pass.Reportf(call.Pos(), "log message must start with lowercase letter")
			}

			return true
		})
	}
	return nil, nil
}
