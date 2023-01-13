package chunk

import (
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileReader_scan_error(t *testing.T) {
	pathDir := t.TempDir()

	fReader, err := NewFileReader(pathDir)
	require.NoError(t, err, "failed to open the temp dir during test")

	defer fReader.Close()

	err = fReader.NextLine()

	require.Error(t, err,
		"reading a directory should fail to scan the next line")
	require.Contains(t, err.Error(), "failed to scan the next line",
		"the error should contain the error reason")
}

func TestFileReader_unknown_file(t *testing.T) {
	pathFile := filepath.Join("..", "testdata", "unknown.txt")

	fReader, err := NewFileReader(pathFile)

	require.Error(t, err,
		"non-existing file should return an error")
	require.Contains(t, err.Error(), "failed to open the file",
		"the error should contain the error reason")
	require.Nil(t, fReader,
		"the returned reader should be nil on error")
}

// FileReader shuld be thread safe.
func TestFileReader_opening_the_same_file(t *testing.T) {
	pathFile := filepath.Join("..", "testdata", "small_chunk.txt")

	// Open the first chunk file
	chunk1, err := NewFileReader(pathFile)
	require.NoError(t, err, "failed to open the 1st file")
	require.NotNil(t, chunk1, "successfully opened file should return a reader")

	defer chunk1.Close()

	// Open the same file as a second chunk file
	chunk2, err := NewFileReader(pathFile)
	require.NoError(t, err, "failed to open the 2nd file")
	require.NotNil(t, chunk2, "successfully opened file should return a reader")

	defer chunk2.Close()

	// Loop until both EOF
	for {
		err1 := chunk1.NextLine()
		if err1 != nil && !errors.Is(err1, io.EOF) {
			t.Fatal(err, "chunk1.NextLine failed with unexpected error")
		}

		err2 := chunk2.NextLine()
		if err2 != nil && !errors.Is(err2, io.EOF) {
			t.Fatal(err, "chunk2.NextLine failed with unexpected error")
		}

		if err1 != nil && err2 != nil {
			break
		}

		// Compare the current line
		require.Equal(t, chunk1.CurrentLine(), chunk2.CurrentLine())
	}

	require.True(t, chunk1.IsEOF(), "it should be true if the end of the file is reached")
	require.True(t, chunk2.IsEOF(), "it should be true if the end of the file is reached")

	require.ErrorAs(t, chunk1.NextLine(), &io.EOF,
		"once EOF is reached, it should always return io.EOF")
	require.ErrorAs(t, chunk2.NextLine(), &io.EOF,
		"once EOF is reached, it should always return io.EOF")
}
