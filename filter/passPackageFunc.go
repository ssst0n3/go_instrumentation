package filter

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"go/build"
	"strings"
)

func Demo(pkgPath string) (pass bool) {
	pass = true
	if pkgPath == "main" || pkgPath == "hello-world/pkg" {
		pass = false
	}
	return
}

func BypassRuntimeAndInternal(pkgPath string) (pass bool) {
	if strings.HasPrefix("runtime", pkgPath) {
		pass = true
	}
	if strings.HasPrefix("internal", pkgPath) {
		pass = true
	}
	return
}

func BypassGoSrcPackage(pkgPath string) (pass bool) {
	if len(pkgPath) == 0 {
		pass = true
		return
	}
	if strings.HasPrefix("runtime", pkgPath) {
		pass = true
		return
	}
	if strings.HasPrefix("internal", pkgPath) {
		pass = true
		return
	}
	if strings.Contains(pkgPath, "golang.org") {
		pass = true
		return
	}
	pkg, err := build.Import(pkgPath, "", build.FindOnly)
	awesome_error.CheckWarning(err)
	pass = pkg.Goroot
	return
}

func BypassTooManyDetails(pkgPath string) bool {
	if BypassGoSrcPackage(pkgPath) {
		return true
	}
	if strings.Contains(pkgPath, "github.com/golang/protobuf") {
		return true
	}
	if strings.Contains(pkgPath, "github.com/checkpoint-restore/go-criu") {
		return true
	}
	if pkgPath == "github.com/urfave/cli" {
		return true
	}
	if strings.Contains(pkgPath, "github.com/sirupsen/logrus") {
		return true
	}
	return false
}
