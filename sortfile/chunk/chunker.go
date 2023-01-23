package chunk

import (
	"bufio"
	"io"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// Chunker is a file chunker. It returns a list of paths to the temporary chunked
// files. These files are sorted by lines using the given isLess function. If
// isLess is nil, the default is used.
//
// It is similar to FileSplit() but takes a file pointer instead of file path.
// To chunk the file via file path, use FileSplit().
func Chunker(inFile io.Reader, sizeFileIn datasize.InBytes, sizeChunk datasize.InBytes, isLess func(string, string) bool) ([]string, error) {
	if inFile == nil {
		return nil, errors.New("input file is nil")
	}

	// Chunk the file
	listFileChunk := []string{}
	buf := bufio.NewScanner(inFile)
	index := 0
	line := ""

	for {
		lines := NewLines() // a chunk
		lines.IsLess = isLess

		for buf.Scan() {
			line = buf.Text()

			// Loop until the chunk size is reached
			if !lines.WillOverSize(line, int(sizeChunk)) {
				lines.AppendLine(line)

				continue
			}

			break
		}

		pathFile, err := lines.Dump()
		if err != nil {
			return nil, errors.Wrap(err, "failed to dump the chunk")
		}

		listFileChunk = append(listFileChunk, pathFile)
		index++

		// EOF
		if !buf.Scan() {
			break
		}
	}

	return listFileChunk, nil
}
