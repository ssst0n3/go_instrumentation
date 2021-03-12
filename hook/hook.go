package hook

import (
	"github.com/dave/dst"
	"github.com/ssst0n3/go_instrumentation/stmt"
)

type Hook struct {
	InstrumentationStmt dst.Stmt
	ImportPkg           map[string]string
}

func New(pkgPath string, funcDecl *dst.FuncDecl, stmt stmt.Stmt) *Hook {
	statement, importPkg := stmt.Statement(pkgPath, funcDecl)
	return &Hook{
		InstrumentationStmt: statement,
		ImportPkg:           importPkg,
	}
}
