package datasize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSizeFile_failed_to_open_file(t *testing.T) {
	t.Parallel()

	pathFile := filepath.Join("..", "testdata", "unknown.txt")

	sizeFile, numLines, err := File(pathFile)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open file")
	require.Empty(t, sizeFile, "size of file should be zero on error")
	require.Zero(t, numLines, "number of lines should be zero on error")
}

func TestSizeFile_failed_to_get_file_stat(t *testing.T) {
	pathFileTmp := filepath.Join("..", "testdata", "empty.txt")

	oldOsOpen := OsOpen
	defer func() {
		OsOpen = oldOsOpen
	}()

	// Mock os.Open to return a file pointer to an closed file.
	OsOpen = func(name string) (*os.File, error) {
		filePtr, err := os.Create(pathFileTmp)

		require.NoError(t, err)

		filePtr.Close()

		return filePtr, nil
	}

	sizeFile, numLines, err := File(pathFileTmp)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get file stat")
	require.Empty(t, sizeFile, "size of file should be zero on error")
	require.Zero(t, numLines, "number of lines should be zero on error")
}
