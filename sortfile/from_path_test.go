package sortfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestFromPath_empty_input_path(t *testing.T) {
	t.Parallel()

	err := FromPath("", filepath.Join(os.TempDir(), "test.txt"), false)

	require.Error(t, err, "empty input path should return error")
	require.Contains(t, err.Error(), "failed to get file size",
		"error message should contain the reason")
}

func TestFromPath_empty_output_path(t *testing.T) {
	t.Parallel()

	err := FromPath(filepath.Join("testdata", "size67byte.txt"), "", false)

	require.Error(t, err, "empty input path should return error")
	require.Contains(t, err.Error(), "failed to create the output file",
		"error message should contain the reason")
}

func TestFromPath_fail_get_memory_size(t *testing.T) {
	// Mock/monkey patch the function variable to force error
	oldMemoryGet := datasize.MemoryGet
	defer func() {
		datasize.MemoryGet = oldMemoryGet
	}()

	datasize.MemoryGet = func() (*memory.Stats, error) {
		return nil, errors.New("forced error")
	}

	// Test
	err := FromPath(
		filepath.Join("testdata", "size67byte.txt"),
		filepath.Join(os.TempDir(), "test.txt"),
		false,
	)

	require.Error(t, err, "empty input path should return error")
	require.Contains(t, err.Error(), "failed to get free memory size",
		"error message should contain the reason")
}

func TestFromPath_force_use_external_sort(t *testing.T) {
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")
	pathFileOut := filepath.Join(os.TempDir(), t.Name()+".txt")
	pathFileExpect := filepath.Join("testdata", "sorted_chunks", "expect_out.txt")

	forceExternalSort := true

	// Test
	err := FromPath(
		pathFileIn,
		pathFileOut,
		forceExternalSort,
	)

	require.NoError(t, err, "empty input path should return error")

	expectByte, err := os.ReadFile(pathFileExpect)
	require.NoError(t, err, "failed to read expected output file")

	actualByte, err := os.ReadFile(pathFileOut)
	require.NoError(t, err, "failed to read actual output file")

	require.Equal(t, string(expectByte), string(actualByte), "output file should be sorted")
}
