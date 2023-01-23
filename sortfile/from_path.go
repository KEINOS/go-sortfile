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
//
// It is similar to FromPathFunc() but it uses the default isLess() function.
func FromPath(pathFileIn, pathFileOut string, forceExternalSort bool) error {
	return FromPathFunc(pathFileIn, pathFileOut, forceExternalSort, nil)
}

// FromPath sorts the file by lines and stores the result in the given path.
//
// It will sort in-memory if the file size is smaller than the current free
// memory. Otherwise it will use the external merge sort.
//
// It is similar to FromPath() but it allows you to specify your own isLess()
// function. If isLess is nil, it will use the default isLess() function.
//
//	  // Default isLess function
//	  func isLess(a, b string) bool {
//		     return a < b // to reverse the sort, use a > b
//	  }
func FromPathFunc(pathFileIn, pathFileOut string, forceExternalSort bool, isLess func(string, string) bool) error {
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
		return errors.Wrap(sortInMemory(numLines, fileIn, fileOut, isLess),
			"FromPath failed")
	}

	// External merge sort with sizeMemoryFree as the chunk size
	return errors.Wrap(sortExternalFile(sizeFileIn, sizeMemoryFree, fileIn, fileOut, isLess),
		"FromPath failed")
}

func sortInMemory(numLines int, fileIn io.Reader, fileOut io.Writer, isLess func(string, string) bool) error {
	return errors.Wrap(InMemory(numLines, fileIn, fileOut, isLess),
		"failed to sort in-memory")
}

func sortExternalFile(sizeFileIn, sizeChunkFile datasize.InBytes, fileIn io.Reader, fileOut io.Writer, isLess func(string, string) bool) error {
	return errors.Wrap(ExternalFile(sizeFileIn, sizeChunkFile, fileIn, fileOut, isLess),
		"failed to sort by external merge sort")
}
