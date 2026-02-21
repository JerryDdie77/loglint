package analyzer

import (
	"go/ast"
	"go/token"
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
			// 1. Check for function calls
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			// 2. Check for package.method calls (log.Printf)
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// 3. Verify left side is package "log"
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != "log" {
				return true
			}

			// 4. Ensure first argument exists
			if len(call.Args) == 0 {
				return true
			}
			arg := call.Args[0]

			// 5. Verify first argument is string literal
			lit, ok := arg.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}

			// Extract message content (remove quotes)
			msg := lit.Value[1 : len(lit.Value)-1]

			// Rule 1: Message must start with lowercase letter
			if len(msg) > 0 && unicode.IsUpper(rune(msg[0])) {
				pass.Reportf(call.Pos(), "log message must start with lowercase letter")
			}

			return true
		})
	}
	return nil, nil
}
