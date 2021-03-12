package stmt

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/davecgh/go-spew/spew"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Stmt interface {
	Statement(pkgPath string, decl *dst.FuncDecl) (stmt dst.Stmt, importPkgPath map[string]string)
	ImportCfg(importCfgPath string) (err error)
}

// TODO: collect pkg binary path from importcfg; delete that file when begin
func BuildImportPkg(pkgPath string) (binaryPath string, err error) {
	r, err := regexp.Compile("WORK=(/tmp/go-build[0-9]+)")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	pkg, err := build.Import(pkgPath, "", 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	spew.Fdump(log.Logger.Out, pkg)
	// todo: determine whether use -trimpath
	args := []string{"build"}
	if len(pkg.PkgObj) == 0 {
		args = append(args, "-trimpath")
	}
	args = append(args, []string{"-work", "-buildmode", "archive", "-a", pkg.Dir}...)
	cmd := exec.Command("go", args...)
	log.Logger.Info("pack command: ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	workDir := r.FindSubmatch(output)[1]
	binaryPath = fmt.Sprintf("%s/b001/_pkg_.a", workDir)
	// TODO: read from collect
	//binaryPath = "/tmp/go-build4095183497/b027/_pkg_.a"
	return
}

func AddPackageFileIntoImportCfg(pkgPath string, binaryPath string, importCfgPath string) (err error) {
	// TODO: what about importCfgPath empty
	log.Logger.Warnf("pkg: %s, binary: %s, cfg: %s", pkgPath, binaryPath, importCfgPath)
	content, err := ioutil.ReadFile(importCfgPath)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if strings.Contains(string(content), fmt.Sprintf("packagefile %s=", pkgPath)) {
		return
	}
	pkgFile := fmt.Sprintf("packagefile %s=%s\n", pkgPath, binaryPath)
	log.Logger.Debug(importCfgPath)
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

func ParseStmtFromSource(source string) (stmt dst.Stmt, importPkgPath map[string]string, err error) {
	importPkgPath = map[string]string{}
	parse, err := decorator.ParseFile(token.NewFileSet(), "src.go", source, parser.ParseComments)
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
		var name string
		if importSpec.Name != nil {
			name = importSpec.Name.String()
		}
		importPkgPath[path] = name
	}
	for _, decl := range parse.Decls {
		switch decl.(type) {
		case *dst.FuncDecl:
			stmt = decl.(*dst.FuncDecl).Body
			return
		}
	}

	return
}
