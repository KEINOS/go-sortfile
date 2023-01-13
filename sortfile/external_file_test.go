package sortfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-sortfile/sortfile/chunk"
	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

func TestExternalFile_fail_create_new_chunk_file_reader(t *testing.T) {
	// Mock the chunk.OsOpen to return an error
	oldOsOpen := chunk.OsOpen
	defer func() { chunk.OsOpen = oldOsOpen }()

	chunk.OsOpen = func(name string) (*os.File, error) {
		return nil, errors.New("forced error")
	}

	pathDirTmp := t.TempDir()

	ptrFileIn, err := os.Open(pathDirTmp)
	require.NoError(t, err, "failed to open the temp dir during test")

	out := capturer.CaptureOutput(func() {
		ptrFileOut := os.Stdout

		// External merge sort with sizeMemoryFree as the chunk size
		err = ExternalFile(datasize.MiB, datasize.KiB, ptrFileIn, ptrFileOut)

		require.Error(t, err, "ExternalFile failed during test")
		assert.Contains(t, err.Error(), "failed to create reader for the chunk file",
			"error message should contain the reason")
	})

	require.Empty(t, out, "ExternalFile should not output anything")
}

func TestExternalFile_input_file_pointer_is_nil(t *testing.T) {
	pathDirTmp := t.TempDir()

	ptrFileIn, err := os.Open(pathDirTmp)
	require.NoError(t, err, "failed to open the input file during test")

	err = ExternalFile(datasize.MiB, datasize.KiB, nil, ptrFileIn)

	require.Error(t, err, "input file is a directory should return error")
	require.Contains(t, err.Error(), "failed to split the file into chunks",
		"error message should contain the reason")
}

func TestExternalFile_input_is_smaller_than_chunk_size(t *testing.T) {
	// Input and output file paths
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")
	pathFileExpect := filepath.Join("testdata", "sorted_chunks", "expect_out.txt")

	// Get file and memory information
	sizeFileIn, _, err := datasize.File(pathFileIn)
	require.NoError(t, err, "failed to get file size during test")

	sizeMemoryFree, err := datasize.AvailableMemory()
	require.NoError(t, err, "failed to get free memory size during test")

	// Open the file to read
	fileIn, err := os.Open(pathFileIn)
	require.NoError(t, err, "failed to open the input file during test")

	defer fileIn.Close()

	out := capturer.CaptureOutput(func() {
		fileOut := os.Stdout

		// External merge sort with sizeMemoryFree as the chunk size
		err = ExternalFile(sizeFileIn, sizeMemoryFree, fileIn, fileOut)
		require.NoError(t, err, "ExternalFile failed during test")
	})

	// Read expected output
	expectOutByte, err := os.ReadFile(pathFileExpect)
	require.NoError(t, err, "failed to read the expected output file during test")

	require.Equal(t, string(expectOutByte), out, "ExternalFile failed to sort the file")
}
