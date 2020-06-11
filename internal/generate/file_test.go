package generate

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func runTestInDir(t *testing.T, dir string) func() {
	// change working directory to project root
	wd, err := os.Getwd()
	assert.NilError(t, err)
	assert.NilError(t, os.Chdir(dir))
	return func() {
		assert.NilError(t, os.Chdir(wd))
	}
}
