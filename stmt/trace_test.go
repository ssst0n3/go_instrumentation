package stmt

import (
	"go/ast"
	"testing"
)

func TestTrace_StatementFromSource(t *testing.T) {
	trace := NewTrace()
	trace.StatementFromSource(&ast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: &ast.Ident{Name: "test"},
		Type: &ast.FuncType{
			Func:    0,
			Params:  &ast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Results: nil,
		},
		Body: nil,
	})
}