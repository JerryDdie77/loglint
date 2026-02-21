package analyzer

import (
	"go/ast"
	"go/token"
)

// IsLogCall determines if the given AST node represents a log function call
func IsLogCall(node ast.Node) (*ast.CallExpr, *ast.BasicLit, bool) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return nil, nil, false
	}

	// Check for package.method pattern
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, nil, false
	}

	// Check package name (log, slog, or zap)
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil, nil, false
	}

	switch pkg.Name {
	case "log", "slog", "zap":
	default:
		return nil, nil, false
	}

	// First argument must exist and be string literal
	if len(call.Args) == 0 {
		return nil, nil, false
	}

	arg := call.Args[0]
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil, nil, false
	}

	return call, lit, true
}
