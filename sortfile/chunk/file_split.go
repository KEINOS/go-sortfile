package chunk

import (
	"os"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// FileSplit is a file chunker. It returns a list of paths to the temporary
// chunked files. These files are sorted by lines.
//
// It is similar to Chunker() but takes a file path instead of file pointer. To
// chunk the file via file pointer, use Chunker().
func FileSplit(pathFileIn string, sizeChunk datasize.InBytes) ([]string, error) {
	// Prepare the chunked files
	sizeFileIn, _, err := datasize.File(pathFileIn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file size")
	}

	// Open the file to read
	ptrFileIn, err := os.Open(pathFileIn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the file")
	}

	defer ptrFileIn.Close()

	return Chunker(ptrFileIn, sizeFileIn, sizeChunk)
}
