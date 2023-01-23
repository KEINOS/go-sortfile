package chunk

import (
	"os"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// FileSplit is a file chunker. It returns a list of paths to the temporary
// chunked files. These files are sorted by lines using the given isLess
// function. If isLess is nil, the default is used.
//
// It is similar to Chunker() but takes a file path instead of file pointer. To
// chunk the file via file pointer, use Chunker().
func FileSplit(pathFileIn string, sizeChunk datasize.InBytes, isLess func(string, string) bool) (data []string, err error) {
	// Prepare the chunked files
	sizeFileIn, _, err := datasize.File(pathFileIn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file size")
	}

	// Open the file to read.
	// Not checking the error is intentional. The error	check is done in the
	// above datasize.File function. Thus, it never reaches here.
	ptrFileIn, _ := os.Open(pathFileIn)

	defer func() {
		err = ptrFileIn.Close()
	}()

	return Chunker(ptrFileIn, sizeFileIn, sizeChunk, isLess)
}
