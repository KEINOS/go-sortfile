package chunk

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

func TestMergeSorter_different_size_of_chunks(t *testing.T) {
	chunks := make([]*FileReader, 3)

	chunks[0] = NewIOReader(strings.NewReader(heredoc.Doc(`
		Alice
		Bob
		Carol
		Charlie
		Dave
		Ellen
		Eve
		Frank
		Isaac
		Ivan
		Justin
		Mallet
		Mallory
		Marvin
		Matilda
		Oscar
		Pat
		Peggy
	`)))
	chunks[1] = NewIOReader(strings.NewReader(heredoc.Doc(`
		Trent
		Trudy
		Victor
		Walter
	`)))
	chunks[2] = NewIOReader(strings.NewReader(heredoc.Doc(`
		Steve
		Zoe
	`)))

	memFree, err := datasize.AvailableMemory()
	require.NoError(t, err, "failed to get the available memory during the test")

	var buf bytes.Buffer

	mergeSorter := NewMergeSorter(chunks, NewIOWriter(&buf, memFree))

	out := capturer.CaptureStdout(func() {
		err := mergeSorter.Sort()
		require.NoError(t, err, "failed to sort the chunks")
	})

	expectByte, err := os.ReadFile(filepath.Join("..", "testdata", "sorted_chunks", "expect_out.txt"))
	require.NoError(t, err, "failed to read the expected output file")

	require.Equal(t, string(expectByte), buf.String(), "the output is not as expected")
	require.Empty(t, out, "the output should be empty")
}

func TestMergeSorter_Sort_fail_to_create_done_group(t *testing.T) {
	mergeSorter := MergeSorter{}

	err := mergeSorter.Sort()

	require.Error(t, err, "it should error when the number of chunks is 0")
	require.Contains(t, err.Error(), "failed to create a new DoneGroup",
		"error message should contain the error reason")
}

func TestMergeSorter_Sort_fail_to_initialize(t *testing.T) {
	fReader, err := NewFileReader(t.TempDir())
	require.NoError(t, err, "failed to open the temp dir during test")

	mergeSorter := NewMergeSorter([]*FileReader{
		fReader,
	}, nil)

	err = mergeSorter.Sort()

	require.Error(t, err, "it should error when the number of chunks is 0")
	require.Contains(t, err.Error(), "failed to read the first line during initialization",
		"error message should contain the error reason")
}

type DummyWriter struct{}

func (dw DummyWriter) Write([]byte) (int, error) {
	return 0, errors.New("forced error")
}

func TestMergeSorter_Sort_fail_to_write(t *testing.T) {
	fReader, err := NewFileReader(
		filepath.Join("..", "testdata", "sorted_chunks", "input_shuffled.txt"),
	)
	require.NoError(t, err, "failed to open the temp dir during test")

	fWriter := NewIOWriter(DummyWriter{}, 16)

	mergeSorter := NewMergeSorter(
		[]*FileReader{fReader},
		fWriter,
	)

	err = mergeSorter.Sort()

	require.Error(t, err, "it should error when the number of chunks is 0")
	require.Contains(t, err.Error(), "failed to write the line",
		"error message should contain the error reason")
}

type DummyReader struct {
	CountMax int // count to return the error
	CountCur int // current count
}

func (dr DummyReader) Read(readme []byte) (int, error) {
	dr.CountCur++

	if dr.CountCur > dr.CountMax {
		return 0, errors.New("count exceeded")
	}

	// Do not set line break at the end of the line so that bufio.Scanner.Scan()
	// fails to read the next line because of impossible count.
	data := fmt.Sprintf("dummy line %d/%d ", dr.CountCur, dr.CountMax)
	copy(readme, data)

	return len(data), nil
}

func TestMergeSorter_Sort_fail_to_read(t *testing.T) {
	fReader := NewIOReader(DummyReader{CountMax: 3})
	fWriter := NewIOWriter(&bytes.Buffer{}, 16)

	mergeSorter := NewMergeSorter(
		[]*FileReader{fReader},
		fWriter,
	)

	err := mergeSorter.Sort()

	require.Error(t, err, "it should error if failed to read the next line")
	require.Contains(t, err.Error(), "failed to read the next line",
		"error message should contain the error reason")
}
