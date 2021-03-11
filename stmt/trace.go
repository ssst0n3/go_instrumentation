package stmt

import (
	"fmt"
	"github.com/ssst0n3/awesome_libs"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type Trace struct {
	Stmt          ast.Stmt
	ImportPkgPath []string
}

func NewTrace() *Trace {
	return &Trace{
	}
}

func GetParamsFromFuncDecl(funcDecl *ast.FuncDecl) (params []string) {
	for _, param := range funcDecl.Type.Params.List {
		params = append(params, param.Names[0].String())
	}
	return
}

func (t *Trace) Statement(decl *ast.FuncDecl) (stmt ast.Stmt) {
	params := GetParamsFromFuncDecl(decl)
	format := ""
	var formatArgs []*ast.Ident
	for _, param := range params {
		format += fmt.Sprintf(" %s=%%v ", param)
		formatArgs = append(formatArgs, &ast.Ident{
			Name: param,
		})
	}
	format = fmt.Sprintf("[TRACE] %s(%s)\\n", decl.Name, format)
	args := []ast.Expr{&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}}
	for _, arg := range formatArgs {
		args = append(args, arg)
	}
	source := awesome_libs.Format(`package stmt
import "log"

func __stmt() {.left}
	log.Printf("{.format}")
}`, awesome_libs.Dict{
		"format": format,
		"args":   strings.Join(params, ","),
		"left":   "{",
	})
	var err error
	t.Stmt, t.ImportPkgPath, err = ParseStmtFromSource(source)
	log.Logger.Warn(t.Stmt)
	log.Logger.Warn(t.ImportPkgPath)
	awesome_error.CheckFatal(err)
	return t.Stmt
}

func (t *Trace) StatementOld(decl *ast.FuncDecl) (stmt ast.Stmt) {
	params := GetParamsFromFuncDecl(decl)
	format := ""
	var formatArgs []*ast.Ident
	for _, param := range params {
		format += fmt.Sprintf(" %s=%%v ", param)
		formatArgs = append(formatArgs, &ast.Ident{
			Name: param,
		})
	}
	format = fmt.Sprintf("[TRACE] %s(%s)\n", decl.Name, format)
	args := []ast.Expr{&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}}
	for _, arg := range formatArgs {
		args = append(args, arg)
	}
	t.Stmt = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "log.Printf",
					},
					Args: args,
					//Args: []ast.Expr{
					//	&ast.BasicLit{
					//		Kind:  token.STRING,
					//		Value: strconv.Quote("hello"),
					//	},
					//},
				},
			},
		},
	}
	t.ImportPkgPath = []string{"log"}
	return t.Stmt
}

func (t *Trace) ImportCfg(importCfgPath string) (err error) {
	for _, pkgPath := range t.ImportPkgPath {
		// TODO: build import pkg
		binaryPath, err := BuildImportPkg(pkgPath)
		if err != nil {
			return err
		}
		// TODO: add import pkg
		err = AddPackageFileIntoImportCfg(pkgPath, binaryPath, importCfgPath)
		if err != nil {
			return err
		}
	}
	return
}
