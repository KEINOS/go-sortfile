//go:generate go run ./testdata/gen_huge_file.go
package sortfile

import (
	"io"
	"os"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// FromPath sorts the file by lines and stores the result in the given path.
//
// It will sort in-memory if the file size is smaller than the current free
// memory. Otherwise it will use the external merge sort.
func FromPath(pathFileIn, pathFileOut string, forceExternalSort bool) error {
	// Get file and memory information
	sizeFileIn, numLines, err := datasize.File(pathFileIn)
	if err != nil {
		return errors.Wrap(err, "failed to get file size")
	}

	sizeMemoryFree, err := datasize.AvailableMemory()
	if err != nil {
		return errors.Wrap(err, "failed to get free memory size")
	}

	isInMemory := true
	if sizeMemoryFree.IsSmallerThan(sizeFileIn) || forceExternalSort {
		isInMemory = false
	}

	// Open the file to read. Error is not checked since the previous functions
	// already checked the file existence.
	fileIn, _ := os.Open(pathFileIn)

	defer fileIn.Close()

	// Open/create the file to write
	fileOut, err := os.Create(pathFileOut)
	if err != nil {
		return errors.Wrap(err, "failed to create the output file")
	}

	defer fileOut.Close()

	// Sort file in-memory
	if isInMemory {
		// Sort file by external merge sort
		return errors.Wrap(sortInMemory(numLines, fileIn, fileOut),
			"FromPath failed")
	}

	// External merge sort with sizeMemoryFree as the chunk size
	return errors.Wrap(sortExternalFile(sizeFileIn, sizeMemoryFree, fileIn, fileOut),
		"FromPath failed")
}

func sortInMemory(numLines int, fileIn io.Reader, fileOut io.Writer) error {
	return errors.Wrap(InMemory(numLines, fileIn, fileOut),
		"failed to sort in-memory")
}

func sortExternalFile(sizeFileIn, sizeChunkFile datasize.InBytes, fileIn io.Reader, fileOut io.Writer) error {
	return errors.Wrap(ExternalFile(sizeFileIn, sizeChunkFile, fileIn, fileOut),
		"failed to sort by external merge sort")
}
