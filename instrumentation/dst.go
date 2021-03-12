package instrumentation

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation/hook"
	"github.com/ssst0n3/go_instrumentation/stmt"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"io"
	"os"
	"path/filepath"
)

type Dst struct {
	PkgPath           string
	ImportCfgPath     string
	FileSet           *token.FileSet
	ParsedFiles       map[string]*dst.File
	ParsedFileSources map[*dst.File]string
	Stmt              stmt.Stmt
	Instrumented      []*hook.Hook
	InstrumentedFiles []*dst.File
	ImportPkg         map[*dst.FuncDecl]map[string]string
}

func NewDst(pkgPath string, importCfgPath string, stmt stmt.Stmt) *Dst {
	return &Dst{
		PkgPath:           pkgPath,
		ImportCfgPath:     importCfgPath,
		FileSet:           token.NewFileSet(),
		ParsedFiles:       map[string]*dst.File{},
		ParsedFileSources: map[*dst.File]string{},
		ImportPkg:         map[*dst.FuncDecl]map[string]string{},
		Stmt:              stmt,
	}
}

func (i *Dst) AddFile(filepath string) (err error) {
	mode := parser.ParseComments
	file, err := decorator.ParseFile(i.FileSet, filepath, nil, mode)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	i.ParsedFiles[filepath] = file
	i.ParsedFileSources[file] = filepath
	return
}

func (i *Dst) Instrument() (err error) {
	root, err := dst.NewPackage(i.FileSet, i.ParsedFiles, nil, nil)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	dstutil.Apply(root, i.Pre, i.Post)
	return nil
}

func (i *Dst) Pre(cursor *dstutil.Cursor) (result bool) {
	switch node := cursor.Node().(type) {
	case *dst.FuncDecl:
		i.InstrumentFuncDeclPre(node)
		return false
	}
	return true
}

func (i *Dst) InstrumentFuncDeclPre(funcDecl *dst.FuncDecl) {
	h := hook.New(i.PkgPath, funcDecl, i.Stmt)
	i.Instrumented = append(i.Instrumented, h)
	if funcDecl != nil && funcDecl.Body != nil {
		funcDecl.Body.List = append([]dst.Stmt{h.InstrumentationStmt}, funcDecl.Body.List...)
	}
	i.ImportPkg[funcDecl] = h.ImportPkg
}

func (i *Dst) Post(cursor *dstutil.Cursor) (result bool) {
	switch node := cursor.Node().(type) {
	case *dst.File:
		i.InstrumentFilePost(node)
	}
	return true
}

func (i *Dst) InstrumentFilePost(file *dst.File) {
	if len(i.Instrumented) == 0 {
		// Nothing got instrumented
		return
	}
	i.Instrumented = nil
	i.InstrumentedFiles = append(i.InstrumentedFiles, file)
}

func (i *Dst) WriteInstrumentedFiles(buildDirPath string) (src2dst map[string]string, err error) {
	src2dst = make(map[string]string, len(i.InstrumentedFiles))
	for _, node := range i.InstrumentedFiles {
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

func (i *Dst) WriteFile(file *dst.File, w io.Writer) (err error) {
	fileset, astFile, err := decorator.RestoreFile(file)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	// AddImport
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*dst.FuncDecl); ok {
			for path, name := range i.ImportPkg[funcDecl] {
				log.Logger.Info("uses import ", path, astutil.UsesImport(astFile, path))
				//if astutil.UsesImport(astFile, path) {
				//	continue
				//}
				log.Logger.Info("add named import ", name, path)
				astutil.AddNamedImport(fileset, astFile, name, path)
				// write to importcfg
				err := i.Stmt.ImportCfg(i.ImportCfgPath)
				if err != nil {
					awesome_error.CheckErr(err)
					return err
				}
			}
		}
	}
	return printer.Fprint(w, fileset, astFile)
}
