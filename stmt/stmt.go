package stmt

import (
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type Stmt interface {
	Statement(decl *ast.FuncDecl) (stmt ast.Stmt)
	ImportCfg(importCfgPath string) (err error)
}

func BuildImportPkg(pkgPath string) (binaryPath string, err error) {
	r, err := regexp.Compile("WORK=(/tmp/go-build[0-9]+)")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	pkg, err := build.Import("log", "", build.FindOnly)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	cmd := exec.Command("go", "build", "-work", "-buildmode", "archive", "-a", pkg.Dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	workDir := r.FindSubmatch(output)[1]
	binaryPath = fmt.Sprintf("%s/b001/_pkg_.a", workDir)
	return
}

func AddPackageFileIntoImportCfg(pkgPath string, binaryPath string, importCfgPath string) (err error) {
	// TODO: what about importCfgPath empty
	pkgFile := fmt.Sprintf("packagefile %s=%s", pkgPath, binaryPath)
	log.Logger.Error(importCfgPath)
	f, err := os.OpenFile(importCfgPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(pkgFile)
	awesome_error.CheckErr(err)
	return
}

func ParseStmtFromSource(source string) (stmt ast.Stmt, importPkgPath []string, err error) {
	parse, err := parser.ParseFile(token.NewFileSet(), "src.go", source, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, importSpec := range parse.Imports {
		path, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			awesome_error.CheckErr(err)
			return nil, nil, err
		}
		importPkgPath = append(importPkgPath, path)
	}
	stmt = parse.Decls[1].(*ast.FuncDecl).Body
	return
}
