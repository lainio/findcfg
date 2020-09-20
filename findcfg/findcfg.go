package findcfg

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `finds and prints a control flow graph of a particular function

The findcfg analysis reports a functions CFG if it exists for
the particular name.`

var Analyzer = &analysis.Analyzer{
	Name:             "findcfg",
	Doc:              Doc,
	Run:              run,
	RunDespiteErrors: true,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		ctrlflow.Analyzer,
	},
}

var name string // -name flag

func init() {
	Analyzer.Flags.StringVar(&name, "name", name, "name of the function to find")
}

func run(pass *analysis.Pass) (interface{}, error) {
	var funcDecl *ast.FuncDecl

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Name.Name == name {
				funcDecl = n
			}
		}
	})

	if funcDecl != nil {
		CFGs := pass.ResultOf[ctrlflow.Analyzer].(*ctrlflow.CFGs)
		if cfg := CFGs.FuncDecl(funcDecl); cfg != nil {
			pass.Reportf(funcDecl.Pos(), "\n%s", cfg.Format(pass.Fset))
		}
	}

	return nil, nil
}
