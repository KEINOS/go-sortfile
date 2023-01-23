package chunk

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileWriter_Close_buffer_not_empty(t *testing.T) {
	pathFileTmp := filepath.Join(t.TempDir(), t.Name()+".txt")

	fWriter, err := NewFileWriter(pathFileTmp, 32)
	require.NoError(t, err)

	defer fWriter.Close()

	fWriter.buf = []byte("hello world")

	err = fWriter.Close()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "buffer is not empty. Call Done() before Close()")
}

func TestNewFileWriter_path_is_dir(t *testing.T) {
	pathDirTmp := t.TempDir()

	fWriter, err := NewFileWriter(pathDirTmp, 32)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open the file",
		"error should contain the error reason")
	assert.Contains(t, err.Error(), pathDirTmp,
		"error should contain the wrapped error")
	require.Nil(t, fWriter, "returned writer must be nil on error")
}

func TestNewIOWriter(t *testing.T) {
	fWriter := NewIOWriter(nil, 32)

	require.NoError(t, fWriter.Close())
}
