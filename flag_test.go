package go_instrumentation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCompileFlag(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		args := []string{"-p", "main", "-o", "/tmp"}
		cmdFlags, err := ParseCompileFlag(args)
		assert.NoError(t, err)
		assert.Equal(t, "main", cmdFlags.PkgPath)
		assert.Equal(t, "/tmp", cmdFlags.Output)
	})
	t.Run("full", func(t *testing.T) {
		args := []string{
			"-o", "/tmp/go-build386194881/b005/_pkg_.a",
			"-trimpath", "/tmp/go-build386194881/b005=>",
			"-p", "internal/unsafeheader",
			"-std", "-complete", "-buildid", "pOrk0J8syj5lw9uf8eos/pOrk0J8syj5lw9uf8eos",
			"-goversion", "go1.15.7",
			"-D", "-importcfg", "/tmp/go-build386194881/b005/importcfg",
			"-pack", "-c=8", "/usr/local/go/src/internal/unsafeheader/unsafeheader.go",
		}
		cmdFlags, err := ParseCompileFlag(args)
		assert.NoError(t, err)
		assert.Equal(t, "internal/unsafeheader", cmdFlags.PkgPath)
		assert.Equal(t, "/tmp/go-build386194881/b005/_pkg_.a", cmdFlags.Output)
		assert.Equal(t, "/tmp/go-build386194881/b005/importcfg", cmdFlags.ImportCfgPath)
	})
}
