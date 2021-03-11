package filter

import (
	"github.com/davecgh/go-spew/spew"
	log2 "github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"go/build"
	"os/exec"
	"testing"
)

func TestBypassGoSrcPackage(t *testing.T) {
	{
		pkg, err := build.Import("log", "", build.FindOnly)
		assert.NoError(t, err)
		spew.Dump(pkg)

		cmd := exec.Command("go", "build", "-work", "-buildmode", "archive", "-a", pkg.Dir)
		log2.Logger.Info(cmd.String())
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		log2.Logger.Info(string(output))
	}
}
