package chunk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSplit_input_file_not_found(t *testing.T) {
	pathDirTmp := t.TempDir()

	data, err := FileSplit(pathDirTmp, 32, nil)

	require.Error(t, err, "it should error if the path is not a file")
	require.Nil(t, data, "returned data must be nil on error")
	assert.Contains(t, err.Error(), "failed to get file size",
		"error should contain the error reason")
	assert.Contains(t, err.Error(), "failed to count lines in file",
		"returned error should contain the wrapped error")
}
