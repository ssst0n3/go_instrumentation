package importcfg

import (
	"encoding/json"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

type PackageFile struct {
	CfgPath    string
	PkgPath    string
	BinaryPath string
}

func NewPackageFile(cfgPath string) PackageFile {
	return PackageFile{
		CfgPath: cfgPath,
	}
}

var WorkDir string
var PathImportMap string
var ImportMap = map[string]string{}

var AlreadyParsed = map[string]bool{}

func (p PackageFile) Load(recursive bool, depth int) (packageFiles map[string]string) {
	if depth == 0 {
		AlreadyParsed = map[string]bool{}
	}
	AlreadyParsed[p.CfgPath] = true
	packageFiles = map[string]string{}
	r, err := regexp.Compile("packagefile ([a-zA-Z0-9./_-]+)=(/tmp/go-build[0-9]+/[b0-9]+/_pkg_.a)")
	content, err := ioutil.ReadFile(p.CfgPath)
	awesome_error.CheckFatal(err)
	for _, find := range r.FindAllStringSubmatch(string(content), -1) {
		packageFile := PackageFile{
			CfgPath:    path.Dir(find[2]) + "/importcfg",
			PkgPath:    find[1],
			BinaryPath: find[2],
		}
		//log.Logger.Infof("%s %s %s", packageFile.PkgPath, packageFile.CfgPath, packageFile.BinaryPath)
		packageFiles[packageFile.PkgPath] = packageFile.BinaryPath
		if recursive {
			if _, ok := AlreadyParsed[packageFile.CfgPath]; ok {
				continue
			}
			for k, v := range packageFile.Load(recursive, depth+1) {
				packageFiles[k] = v
			}
		}
	}
	return
}

func Collect(importCfgPath string) {
	if len(importCfgPath) > 0 {
		WorkDir = path.Dir(path.Dir(importCfgPath))
		PathImportMap = WorkDir + "/import_map"
		if _, err := os.Stat(PathImportMap); !os.IsNotExist(err) {
			content, err := ioutil.ReadFile(PathImportMap)
			awesome_error.CheckFatal(err)
			awesome_error.CheckFatal(json.Unmarshal(content, &ImportMap))
		}

		packageFiles := NewPackageFile(importCfgPath).Load(true, 0)
		for pkgPath, binaryPath := range packageFiles {
			ImportMap[pkgPath] = binaryPath
		}
		Dump(importCfgPath)
	}
}

func Dump(importCfgPath string) {
	marshal, err := json.Marshal(ImportMap)
	awesome_error.CheckFatal(err)
	awesome_error.CheckFatal(ioutil.WriteFile(PathImportMap, marshal, 0644))
}
