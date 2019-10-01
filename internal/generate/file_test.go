package generate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func runTestInDir(t *testing.T, dir string) func() {
	// change working directory to project root
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	return func() {
		require.NoError(t, os.Chdir(wd))
	}
}
