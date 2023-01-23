package chunk

import (
	"io"
	"os"
	"strings"

	"github.com/KEINOS/go-sortfile/sortfile/inmemory"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: Lines
// ----------------------------------------------------------------------------

// Lines holds a slice of string as a chunk of data with sort and save functionality.
//
// This is a helper object to create a temporary file with sorted lines for K-way
// merge sort process.
// To simply write huge number of lines to a file, use FileWriter object instead.
type Lines struct {
	// IsLess is the function to compare two strings during chunk file creation.
	// This function must be the same as the one to be used for merge-sorting.
	IsLess   func(a, b string) bool
	lines    []string
	sizeCurr uint64
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewLines returns a new object of Lines.
//
// By default it uses slice.Sort to sort the lines. Set IsLess to use a custom
// function to compare two strings while sorting.
func NewLines() Lines {
	return Lines{
		IsLess:   nil,
		lines:    []string{},
		sizeCurr: 0,
	}
}

// osCreateTemp is a copy of os.CreateTemp to ease testing.
var osCreateTemp = os.CreateTemp

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// AppendLine appends the given line to the chunk.
func (l *Lines) AppendLine(line string) {
	line = l.UniformLineBreak(line)

	l.lines = append(l.lines, line)
	l.sizeCurr += uint64(len(line))
}

// Dump sorts and writes the lines in the chunk to a temporary file and returns
// the path to the file.
func (l *Lines) Dump() (string, error) {
	file, err := osCreateTemp(os.TempDir(), "sortfile-*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create a temporary file")
	}

	defer file.Close()

	if err := l.WriteSortedLines(file); err != nil {
		return "", errors.Wrap(err, "failed to write sorted lines")
	}

	return file.Name(), nil
}

// Lines returns the lines in the chunk (a slice of string) as is.
func (l *Lines) Lines() []string {
	return l.lines
}

// Size returns the byte size of the chunk to be written to the output.
func (l *Lines) Size() int {
	return int(l.sizeCurr)
}

// SizeRaw is similar to the Size method but it calculates the size from scratch
// an update the cached value. Thus, it is slower than Size().
//
// It must return the same value as Size() method by design. It is for debugging
// purpose only.
func (l *Lines) SizeRaw() int {
	size := 0

	for _, line := range l.lines {
		size += len(line)
	}

	l.sizeCurr = uint64(size)

	return size
}

// UniformLineBreak returns the given line with the uniformed line break at the end.
func (l Lines) UniformLineBreak(line string) string {
	const cutset = "\r\n"

	return strings.TrimRight(line, cutset) + "\n"
}

// WillOverSize returns true if the given line will make the chunk over the
// sizeMax, the size limit.
func (l *Lines) WillOverSize(line string, sizeMax int) bool {
	return l.Size()+len(l.UniformLineBreak(line)) > sizeMax
}

// WriteSortedLines writes the sorted lines in the chunk to the given output.
func (l *Lines) WriteSortedLines(output io.Writer) error {
	if l.IsLess != nil {
		inmemory.SortSliceFunc(l.lines, l.IsLess)
	} else {
		inmemory.SortSlice(l.lines)
	}

	result := strings.Join(l.lines, "")

	_, err := output.Write([]byte(result))

	return errors.Wrap(err, "failed to dump the final output")
}
