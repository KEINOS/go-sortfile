package chunk

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: FileReader
// ----------------------------------------------------------------------------

// FileReader is a line reader for the chunk file.
//
// It aims to provide a simple interface for K-way merge sort usage to read the
// chunk file line by line.
type FileReader struct {
	file    io.Reader
	scanner *bufio.Scanner
	closer  func() error
	line    string
	isEOF   bool
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// OsOpen is a copy of os.Open() as a dependency injection to ease testing.
var OsOpen = os.Open

// NewFileReader returns a new FileReader object.
//
// It will open the file of the given path but it will not read the initial first
// line. The caller should call NextLine() to read the first line right after
// the FileReader.Close() is deferred.
// See the example_test.go for the actual use case.
func NewFileReader(path string) (*FileReader, error) {
	file, err := OsOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the file")
	}

	reader := NewIOReader(file)

	return reader, nil
}

// NewIOReader retruns a new FileReader object.
//
// It is similar to NewFileReader() but it takes io.Reader instead of file path.
func NewIOReader(reader io.Reader) *FileReader {
	return &FileReader{
		line:    "",
		file:    reader,
		scanner: bufio.NewScanner(reader),
		closer: func() error {
			return nil
		},
		isEOF: false,
	}
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// Close closes the current chunk file.
//
// It is the callers responsibility to close the file. Use defer to close the file.
func (f *FileReader) Close() error {
	return errors.Wrap(f.closer(), "failed to close the file")
}

// CurrentLine returns the line currently read from the file.
//
// It will return the same line until NextLine() is called. If the line is used
// or selected for merge sort, the caller should call NextLine() to move to the
// next line.
func (f *FileReader) CurrentLine() string {
	return f.line
}

// IsEOF returns true if the end of the file is reached.
func (f *FileReader) IsEOF() bool {
	return f.isEOF
}

// NextLine reads the next line from the file and sets it to the CurrentLine().
//
// Once it reaches the end of the file, it will return io.EOF error.
func (f *FileReader) NextLine() error {
	if f.isEOF {
		return io.EOF
	}

	if f.scanner.Scan() {
		f.line = f.scanner.Text()

		return nil
	}

	if f.scanner.Err() != nil {
		return errors.Wrap(f.scanner.Err(), "failed to scan the next line")
	}

	f.isEOF = true

	return io.EOF
}
