package chunk

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLines_Dump_failed_to_create_temp_file(t *testing.T) {
	// Backup and defer restore osCreateTemp
	oldOsCreateTemp := osCreateTemp
	defer func() {
		osCreateTemp = oldOsCreateTemp
	}()

	// Mock osCreateTemp to force error
	osCreateTemp = func(dir, pattern string) (*os.File, error) {
		return nil, errors.New("forced error")
	}

	lines := NewLines()

	pathFileTmp, err := lines.Dump()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create a temporary file",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "forced error",
		"it should contain the wrapped error")
	assert.Equal(t, "", pathFileTmp,
		"the returned path should be empty on error")
}

func TestLines_Dump_failed_to_write_temp_file(t *testing.T) {
	// Backup and defer restore osCreateTemp
	oldOsCreateTemp := osCreateTemp
	defer func() {
		osCreateTemp = oldOsCreateTemp
	}()

	// Mock osCreateTemp to force return a path to a directory.
	osCreateTemp = func(dir, pattern string) (*os.File, error) {
		return os.Open(t.TempDir())
	}

	lines := NewLines()

	lines.AppendLine("charlie")
	lines.AppendLine("bob")
	lines.AppendLine("alice")

	pathFileTmp, err := lines.Dump()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write sorted lines",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "failed to dump the final output",
		"it should contain the wrapped error")
	assert.Equal(t, "", pathFileTmp,
		"the returned path should be empty on error")
}
