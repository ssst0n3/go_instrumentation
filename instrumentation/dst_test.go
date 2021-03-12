package instrumentation

import (
	"github.com/dave/dst/decorator"
	"github.com/davecgh/go-spew/spew"
	"github.com/ssst0n3/go_instrumentation/stmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDst_Instrument(t *testing.T) {
	dst := NewDst("", "", stmt.NewTrace())
	dst.AddFile("/home/st0n3/pentest_project/go_instrumentation/test/data/hello-world/hello-world.go")
	assert.NoError(t, dst.Instrument())
	spew.Dump(dst.InstrumentedFiles)

	_, file, err := decorator.RestoreFile(dst.InstrumentedFiles[0])
	assert.NoError(t, err)
	spew.Dump(file)
	//var buff bytes.Buffer

}
