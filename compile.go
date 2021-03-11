package go_instrumentation

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation/instrumentation"
	"github.com/ssst0n3/go_instrumentation/stmt"
	"path/filepath"
	"strings"
)

func Compile(args []string, stmt stmt.Stmt, passPackageFunc func(pkgPath string) bool) (newArgs []string, err error) {
	newArgs = args
	flag, err := ParseCompileFlag(args[1:])
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Warn(flag.PkgPath)
	// TODO: hook except internal,runtime
	if passPackageFunc(flag.PkgPath) {
		return
	}
	log.Logger.Error(flag.ImportCfgPath)
	i := instrumentation.NewAst(flag.PkgPath, flag.ImportCfgPath, stmt)
	filepathList := ParseArgs(args)
	for _, src := range filepathList {
		err = i.AddFile(src)
		if err != nil {
			awesome_error.CheckErr(err)
			return nil, err
		}
	}
	instrumented, err := i.Instrument()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	src2dst, err := i.WriteInstrumentedFiles(filepath.Dir(flag.Output), instrumented)
	argsJoin := strings.Join(args, "\x00")
	for src, dest := range src2dst {
		argsJoin = strings.ReplaceAll(argsJoin, src, dest)
	}
	newArgs = strings.Split(argsJoin, "\x00")
	return
}

func ParseArgs(args []string) (filepath []string) {
	for _, arg := range args {
		if strings.HasSuffix(arg, ".go") {
			filepath = append(filepath, arg)
		}
	}
	return
}
