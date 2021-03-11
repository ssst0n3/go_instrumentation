package stmt

import (
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildImportPkg(t *testing.T) {
	binaryPath, err := BuildImportPkg("log")
	assert.NoError(t, err)
	log.Logger.Info(binaryPath)
}