package go_instrumentation

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation/importcfg"
	"github.com/ssst0n3/go_instrumentation/stmt"
)

func Link(args []string) (newArgs []string) {
	flag, err := ParseCompileFlag(args[1:])
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if len(flag.ImportCfgPath) > 0 {
		log.Logger.Error(flag.ImportCfgPath)
		p := importcfg.NewPackageFile(flag.ImportCfgPath)
		//binaryPath, err := stmt.BuildImportPkg("log")
		//if err != nil {
		//	return
		//}
		// TODO: add import pkg
		origin := p.Load(false, 0)
		log.Logger.Warnf("origin: %v", origin)
		for pkgPath, binaryPath := range p.Load(true, 0) {
			log.Logger.Warnf("pkgPath: %v", pkgPath)
			if _, ok := origin[pkgPath]; ok {
				continue
			}
			err = stmt.AddPackageFileIntoImportCfg(pkgPath, binaryPath, flag.ImportCfgPath)
			if err != nil {
				return
			}
		}

	}
	return args
}
