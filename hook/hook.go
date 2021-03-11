package hook

import (
	"github.com/ssst0n3/go_instrumentation/stmt"
	"go/ast"
)

type Ast struct {
	InstrumentationStmt ast.Stmt
}

func NewAst(funcDecl *ast.FuncDecl, stmt stmt.Stmt) *Ast {
	return &Ast{
		InstrumentationStmt: stmt.Statement(funcDecl),
	}
}
