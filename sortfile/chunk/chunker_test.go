package chunk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunker_input_is_nil(t *testing.T) {
	chunkList, err := Chunker(nil, 0, 0, nil)

	require.Error(t, err, "it should error on nil input")
	require.Contains(t, err.Error(), "input file is nil",
		"error message should contain the error reason")
	require.Nil(t, chunkList, "chunk list should be nil on error")
}

func TestChunker_wailed_to_write_chunk(t *testing.T) {
	// Backup and defer restore osCreateTemp
	oldOsCreateTemp := osCreateTemp
	defer func() {
		osCreateTemp = oldOsCreateTemp
	}()

	// Mock osCreateTemp to force return a path to a directory.
	osCreateTemp = func(dir, pattern string) (*os.File, error) {
		return nil, errors.New("forced error")
	}

	// Prepare the test file
	pathFileTest := filepath.Join("..", "testdata", "sorted_chunks", "input_shuffled.txt")

	sizeFile, _, err := datasize.File(pathFileTest)
	require.NoError(t, err, "it should get the file size")

	sizeChunk := datasize.InBytes(8)

	ptrFileIn, err := os.Open(pathFileTest)
	require.NoError(t, err, "failed to open the test file during test setup")

	defer ptrFileIn.Close()

	// Chunk the file now
	chunkList, err := Chunker(ptrFileIn, sizeFile, sizeChunk, nil)

	require.Error(t, err, "it should error if it fails to write the chunk")
	require.Nil(t, chunkList, "chunk list should be nil on error")
	assert.Contains(t, err.Error(), "failed to dump the chunk",
		"error message should contain the error reason")
	assert.Contains(t, err.Error(), "forced error",
		"error message should contain the wrapped error")
}
