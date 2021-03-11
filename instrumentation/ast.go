package instrumentation

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation/hook"
	"github.com/ssst0n3/go_instrumentation/stmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"io"
	"os"
	"path/filepath"
)

type Ast struct {
	PkgPath           string
	ImportCfgPath     string
	FileSet           *token.FileSet
	ParsedFiles       map[string]*ast.File
	ParsedFileSources map[*ast.File]string
	Instrumented      []*hook.Ast
	InstrumentedFiles []*ast.File
	Stmt              stmt.Stmt
}

func NewAst(pkgPath string, importCfgPath string, stmt stmt.Stmt) *Ast {
	return &Ast{
		PkgPath:           pkgPath,
		ImportCfgPath:     importCfgPath,
		FileSet:           token.NewFileSet(),
		ParsedFiles:       make(map[string]*ast.File),
		ParsedFileSources: make(map[*ast.File]string),
		Stmt:              stmt,
	}
}

func (i *Ast) AddFile(filepath string) (err error) {
	file, err := parser.ParseFile(i.FileSet, filepath, nil, parser.ParseComments)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	//spew.Fdump(log.Logger.Out, file)
	i.ParsedFiles[filepath] = file
	i.ParsedFileSources[file] = filepath
	return
}

func (i *Ast) Instrument() ([]*ast.File, error) {
	//root, err := ast.NewPackage(i.FileSet, i.ParsedFiles, nil, nil)
	//if err != nil {
	//	awesome_error.CheckErr(err)
	//	return nil, err
	//}
	root := &ast.Package{
		Name:  i.PkgPath,
		Files: i.ParsedFiles,
	}
	astutil.Apply(root, i.Pre, i.Post)
	return i.InstrumentedFiles, nil
}

func (i *Ast) Pre(cursor *astutil.Cursor) (result bool) {
	switch node := cursor.Node().(type) {
	case *ast.FuncDecl:
		i.InstrumentFuncDeclPre(node)
		return false
	}
	return true
}

func (i *Ast) InstrumentFuncDeclPre(funcDecl *ast.FuncDecl) {
	//log.Logger.Infof("function name: %s; params name: %v", funcDecl.Name, funcDecl.Type.Params.List[0].Names)
	h := hook.NewAst(funcDecl, i.Stmt)
	i.Instrumented = append(i.Instrumented, h)
	funcDecl.Body.List = append([]ast.Stmt{h.InstrumentationStmt}, funcDecl.Body.List...)
}

func (i *Ast) Post(cursor *astutil.Cursor) (result bool) {
	switch node := cursor.Node().(type) {
	case *ast.File:
		i.InstrumentFilePost(node)
	}
	return true
}

func (i *Ast) InstrumentFilePost(file *ast.File) {
	if len(i.Instrumented) == 0 {
		// Nothing got instrumented
		return
	}

	i.Instrumented = nil

	log.Logger.Warn(i.PkgPath, astutil.UsesImport(file, "log"))
	if !astutil.UsesImport(file, "log") {
		astutil.AddImport(i.FileSet, file, "log")
	}

	err := i.Stmt.ImportCfg(i.ImportCfgPath)
	if err != nil {
		os.Exit(-1)
	}
	i.InstrumentedFiles = append(i.InstrumentedFiles, file)
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	printer.Fprint(out, i.FileSet, file)
	out.Flush()
	log.Logger.Infof("%s", b.String())
}

func (i *Ast) WriteInstrumentedFiles(buildDirPath string, instrumentedFiles []*ast.File) (src2dst map[string]string, err error) {
	src2dst = make(map[string]string, len(instrumentedFiles))
	for _, node := range instrumentedFiles {
		//spew.Fdump(log.Logger.Out, node)
		src := i.ParsedFileSources[node]
		filename := filepath.Base(src)
		dest := filepath.Join(buildDirPath, filename)
		output, err := os.Create(dest)
		if err != nil {
			awesome_error.CheckErr(err)
			return nil, err
		}
		defer output.Close()
		// https://golang.org/cmd/compile/#hdr-Compiler_Directives
		_, err = output.WriteString(fmt.Sprintf("//line %s:1\n", src))
		if err != nil {
			awesome_error.CheckErr(err)
			return nil, err
		}
		if err := i.WriteFile(node, output); err != nil {
			awesome_error.CheckErr(err)
			return nil, err
		}
		src2dst[src] = dest
	}
	return
}

func (i *Ast) WriteFile(file *ast.File, out io.Writer) (err error) {
	return printer.Fprint(out, i.FileSet, file)
}
