package importcfg

import (
	"github.com/ssst0n3/awesome_libs/log"
	"os"
	"testing"
)

func TestPackageFile_Load(t *testing.T) {
	log.Logger.Out = os.Stdout
	p := NewPackageFile("/tmp/go-build1316651728/b001/importcfg.link")
	log.Logger.Info(p.Load(true, 0))
}
