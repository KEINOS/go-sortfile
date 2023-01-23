package chunk

import (
	"io"
	"os"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: FileWriter
// ----------------------------------------------------------------------------

// FileWriter is wrapper of file writer but with a buffer.
//
// You can specify a line to be written to a file and buffer it up to a certain
// amount before writing. See the WriteLine method for more information.
//
// The data is added in the specified order, so if it is necessary to sort the
// data or perform other processing, use a Line object.
type FileWriter struct {
	file       io.Writer
	closer     func() error
	lineBreak  string
	buf        []byte
	sizeBufMax datasize.InBytes
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewFileWriter returns a new FileWriter object.
//
// The maxSizeBuf is the max size of the buffer before it flushes the buffer to
// the file.
// It is similar to NewIOWriter but it takes a path instead of a file pointer.
func NewFileWriter(pathFileOut string, maxSizeBuf datasize.InBytes) (*FileWriter, error) {
	file, err := os.Create(pathFileOut)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the file")
	}

	return &FileWriter{
		file:       file,
		buf:        make([]byte, 0, int(maxSizeBuf)),
		sizeBufMax: maxSizeBuf,
		closer: func() error {
			return errors.Wrap(file.Close(), "internal closer failed to close the file")
		},
		lineBreak: "\n",
	}, nil
}

// NewIOWriter returns a new FileWriter object.
//
// The maxSizeBuf is the max size of the buffer before it flushes the buffer to
// the file.
// It is similar to NewFileWriterPath but it takes a file pointer instead of a path.
func NewIOWriter(ptrFileOut io.Writer, maxSizeBuf datasize.InBytes) *FileWriter {
	return &FileWriter{
		file:       ptrFileOut,
		buf:        make([]byte, 0, int(maxSizeBuf)),
		sizeBufMax: maxSizeBuf,
		closer: func() error {
			return nil
		},
		lineBreak: "\n",
	}
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

func (fw *FileWriter) Close() error {
	if len(fw.buf) > 0 {
		return errors.New("buffer is not empty. Call Done() before Close()")
	}

	return errors.Wrap(fw.closer(), "Close() failed")
}

// Done flushes the remaining buffer to the file.
func (fw *FileWriter) Done() error {
	_, err := fw.file.Write(fw.buf)
	if err == nil {
		fw.buf = make([]byte, 0)
	}

	return errors.Wrap(err, "failed to flush the buffer")
}

func (fw *FileWriter) flushBuffer() (int, error) {
	written, err := fw.file.Write(fw.buf)
	fw.buf = make([]byte, 0, int(fw.sizeBufMax))

	return written, errors.Wrap(err, "failed to flush the buffer")
}

// WriteLine writes the line to the file adding a line break at the end.
//
// It will buffer the line until it reaches the max size of the buffer, then flushes
// the buffer to the file.
func (fw *FileWriter) WriteLine(line string) (int, error) {
	line += fw.lineBreak
	written := 0

	// If the line exceeds the max size of the buffer, flush it to the file and
	// clear before appending the new line
	if len(fw.buf)+len(line) > int(fw.sizeBufMax) {
		writtenBuf, err := fw.flushBuffer()

		if err != nil {
			return 0, errors.Wrap(err, "buffer exceeded the max size but failed to flush")
		} else {
			written = writtenBuf
		}
	}

	fw.buf = append(fw.buf, []byte(line)...)

	return written, nil
}
