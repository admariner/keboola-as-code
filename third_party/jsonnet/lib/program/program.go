// Package program provides API for AST pre-processing (desugaring, static analysis).
package program

import (
	"github.com/google/go-jsonnet/ast"

	"github.com/keboola/keboola-as-code/third_party/jsonnet/lib/parser"
)

// SnippetToAST converts a Jsonnet code snippet to a desugared and analyzed AST.
func SnippetToAST(diagnosticFilename ast.DiagnosticFileName, importedFilename, snippet string) (ast.Node, error) {
	node, _, err := parser.SnippetToRawAST(diagnosticFilename, importedFilename, snippet)
	if err != nil {
		return nil, err
	}
	if err := PreprocessAst(&node); err != nil {
		return nil, err
	}
	return node, nil
}

func PreprocessAst(node *ast.Node) error {
	err := desugarAST(node)
	if err != nil {
		return err
	}
	err = analyze(*node)
	if err != nil {
		return err
	}
	return nil
}