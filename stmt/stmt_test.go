package stmt

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/davecgh/go-spew/spew"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"testing"
)

func TestParseStmtFromSource(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		source := `package stmt
import (
	instrument_log "github.com/sirupsen/logrus"
	"os"
)

func __stmt() {
	logger := instrument_log.New()
	file, _ := os.OpenFile("/tmp/instrumentation", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger.SetOutput(file)
	logger.Info()
	file.Close()
}`
		stmt, pkgPath, err := ParseStmtFromSource(source)
		assert.NoError(t, err)
		log.Logger.Info(pkgPath)
		spew.Dump(stmt)
	})
	t.Run("struct", func(t *testing.T) {
		source := `package stmt
import (
	instrument_log "github.com/sirupsen/logrus"
	"os"
)

type A struct{
	test string
}

func (a *A)Stmt() {
	logger := instrument_log.New()
	file, _ := os.OpenFile("/tmp/instrumentation", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger.SetOutput(file)
	logger.Info()
	file.Close()
}`
		parse, err := decorator.ParseFile(token.NewFileSet(), "src.go", source, parser.ParseComments)
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		log.Logger.Info(parse.Decls[2].(*dst.FuncDecl).Name)
		log.Logger.Info(parse.Decls[2].(*dst.FuncDecl).Recv.List[0].Type)
		spew.Dump(parse.Decls[2].(*dst.FuncDecl).Recv.List[0].Type)
	})
}

func TestBuildImportPkg(t *testing.T) {
	{
		binaryPath, err := BuildImportPkg("log")
		assert.NoError(t, err)
		log.Logger.Info(binaryPath)
	}
	{
		binaryPath, err := BuildImportPkg("github.com/ssst0n3/awesome_libs/log")
		assert.NoError(t, err)
		log.Logger.Info(binaryPath)
	}
}
