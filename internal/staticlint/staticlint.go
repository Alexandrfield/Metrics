package staticlint

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitFromMainAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit from main functions of package main",
	Run:  run,
}

func ispectFunc(decl ast.Decl) {
	ast.Inspect(decl, func(n ast.Node) bool {
		// только вызовы функций
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		s, ok := c.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		// только функции Exit пакета os
		if s.Sel.Name == "Exit" && fmt.Sprintf("%s", s.X) == "os" {
			fmt.Printf("%s os.Exit from main function of main packages is denied", s.Sel.String())
		}
		return true
	})
}
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// только пакеты main
		if file.Name.Name != "main" {
			continue
		}
		for _, decl := range file.Decls {
			// только функции main
			funcName, ok := decl.(*ast.FuncDecl)
			if ok && funcName.Name.Name == "main" {
				ispectFunc(decl)
			}
		}
	}
	return nil, nil
}
