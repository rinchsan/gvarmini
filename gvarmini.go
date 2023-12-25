package gvarmini

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "gvarmini is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gvarmini",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		decl, ok := n.(*ast.GenDecl)
		if !ok {
			return
		}
		if decl.Tok != token.VAR {
			return
		}
		for _, spec := range decl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if len(valueSpec.Names) > 3 {
				diag := analysis.Diagnostic{
					Pos:            valueSpec.Pos(),
					Message:        "多いよ",
					SuggestedFixes: []analysis.SuggestedFix{generateFix(pass, decl, valueSpec)},
				}
				pass.Report(diag)
				return
			}
		}
	})

	return nil, nil
}

func generateFix(pass *analysis.Pass, decl *ast.GenDecl, valueSpec *ast.ValueSpec) analysis.SuggestedFix {
	fix := analysis.SuggestedFix{
		Message: "多いよ",
	}
	newDecl := &ast.GenDecl{
		TokPos: decl.TokPos,
		Tok:    decl.Tok,
	}
	for i := range valueSpec.Names {
		vSpec := &ast.ValueSpec{
			Names: []*ast.Ident{valueSpec.Names[i]},
			Type:  valueSpec.Type,
		}
		newDecl.Specs = append(newDecl.Specs, vSpec)
	}
	var buf bytes.Buffer
	format.Node(&buf, pass.Fset, newDecl)
	fix.TextEdits = append(fix.TextEdits, analysis.TextEdit{
		Pos:     decl.Pos(),
		End:     decl.End(),
		NewText: buf.Bytes(),
	})
	return fix
}
