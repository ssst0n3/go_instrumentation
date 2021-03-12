package stmt

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/davecgh/go-spew/spew"
	"github.com/ssst0n3/awesome_libs"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation/importcfg"
	"go/token"
	"strconv"
	"strings"
)

type Trace struct {
	Stmt          dst.Stmt
	ImportPkgPath map[string]string
	PackPath      map[string]string
}

func NewTrace() *Trace {
	return &Trace{
		ImportPkgPath: map[string]string{},
		PackPath:      map[string]string{},
	}
}

func GetParamsFromFuncDecl(funcDecl *dst.FuncDecl) (params []string) {
	for _, param := range funcDecl.Type.Params.List {
		if len(param.Names) > 0 && param.Names[0].String() != "_" {
			params = append(params, param.Names[0].String())
		}
	}
	return
}

func (t *Trace) Statement(pkgPath string, decl *dst.FuncDecl) (stmt dst.Stmt, importPkgPath map[string]string) {
	params := GetParamsFromFuncDecl(decl)
	format := ""
	var formatArgs []*dst.Ident
	for _, param := range params {
		format += fmt.Sprintf(" %s=%%v ", param)
		formatArgs = append(formatArgs, &dst.Ident{
			Name: param,
		})
	}
	funcName := decl.Name.String()
	log.Logger.Info("decl.Recv", decl.Recv)
	spew.Fdump(log.Logger.Out, decl.Recv)
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		var objTypeName string
		switch decl.Recv.List[0].Type.(type) {
		case *dst.StarExpr:
			objTypeName = "&" + decl.Recv.List[0].Type.(*dst.StarExpr).X.(*dst.Ident).String()
		case *dst.Ident:
			objTypeName = decl.Recv.List[0].Type.(*dst.Ident).String()
		}
		funcName = fmt.Sprintf("%s.%s", objTypeName, decl.Name)
	}
	format = fmt.Sprintf("[TRACE] [%s] %s(%s)", pkgPath, funcName, format)
	args := []dst.Expr{&dst.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}}
	for _, arg := range formatArgs {
		args = append(args, arg)
	}
	source := awesome_libs.Format(`package stmt
import (
	instrument_log "github.com/sirupsen/logrus"
	instrument_os "os"
)

func __stmt() {.left}
	logger := instrument_log.New()
	file, _ := instrument_os.OpenFile("/tmp/instrumentation", instrument_os.O_CREATE|instrument_os.O_WRONLY|instrument_os.O_APPEND, 0644)
	logger.SetOutput(file)
	logger.Infof("{.format}",{.args})
	file.Close()
}`, awesome_libs.Dict{
		"format": format,
		"args":   strings.Join(params, ","),
		"left":   "{",
	})
	//log.Logger.Warn(source)
	var err error
	t.Stmt, t.ImportPkgPath, err = ParseStmtFromSource(source)
	//spew.Fdump(log.Logger.Out, t.Stmt)
	//log.Logger.Warn(t.ImportPkgPath)
	awesome_error.CheckFatal(err)
	return t.Stmt, t.ImportPkgPath
}

func (t *Trace) StatementNew(pkgPath string, decl *dst.FuncDecl) (stmt dst.Stmt, importPkgPath map[string]string) {
	params := GetParamsFromFuncDecl(decl)
	format := ""
	var formatArgs []*dst.Ident
	for _, param := range params {
		format += fmt.Sprintf(" %s=%%v ", param)
		formatArgs = append(formatArgs, &dst.Ident{
			Name: param,
		})
	}
	format = fmt.Sprintf("[TRACE] %s %s(%s)\n", pkgPath, decl.Name, format)
	args := []dst.Expr{&dst.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}}
	for _, arg := range formatArgs {
		args = append(args, arg)
	}
	t.Stmt = &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.ExprStmt{
				X: &dst.CallExpr{
					Fun: &dst.Ident{
						//Name: "instrument_log__.Printf",
						Name: "instrument_log__.Logger.Infof",
					},
					Args: args,
				},
			},
		},
	}
	//t.ImportPkgPath["fmt"] = "instrument_log__"
	t.ImportPkgPath["github.com/ssst0n3/awesome_libs/log"] = "instrument_log__"
	return t.Stmt, t.ImportPkgPath
}

func (t *Trace) ImportCfg(importCfgPath string) (err error) {
	for pkgPath := range t.ImportPkgPath {
		if _, has := t.PackPath[pkgPath]; !has {
			if _, has := importcfg.ImportMap[pkgPath]; !has {
				binaryPath, err := BuildImportPkg(pkgPath)
				if err != nil {
					return err
				}
				t.PackPath[pkgPath] = binaryPath
				importcfg.ImportMap[pkgPath] = binaryPath
			} else {
				t.PackPath[pkgPath] = importcfg.ImportMap[pkgPath]
			}
		}
		// TODO: add import pkg
		err = AddPackageFileIntoImportCfg(pkgPath, t.PackPath[pkgPath], importCfgPath)
		if err != nil {
			return err
		}
	}
	return
}
