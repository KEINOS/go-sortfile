package sortfile

import (
	"io"
	"os"

	"github.com/KEINOS/go-sortfile/sortfile/chunk"
	"github.com/KEINOS/go-sortfile/sortfile/datasize"
	"github.com/pkg/errors"
)

// ExternalFile sorts the file using external merge sort (K-way merge sort).
//
// If the sizeFileIn is smaller than the sizeChunk, we recommend to use InMemory
// sort instead.
func ExternalFile(sizeFileIn, sizeChunk datasize.InBytes, ptrFileIn io.Reader, ptrFileOut io.Writer) error {
	// Avoid index out of range with length 0
	if sizeFileIn.IsSmallerThan(sizeChunk) {
		sizeChunk = sizeFileIn
	}

	// Split the file into sorted chunk files
	listChunkFiles, err := chunk.Chunker(ptrFileIn, sizeFileIn, sizeChunk)
	if err != nil {
		return errors.Wrap(err, "failed to split the file into chunks")
	}

	// Merge sort the chunk files
	chunks := make([]*chunk.FileReader, len(listChunkFiles))

	for index, pathFile := range listChunkFiles {
		// Create the chunk file
		reader, err := chunk.NewFileReader(pathFile)
		if err != nil {
			if FileExists(pathFile) {
				_ = os.Remove(pathFile)
			}

			return errors.Wrap(err, "failed to create reader for the chunk file: "+pathFile)
		}

		defer func(fReader *chunk.FileReader) error {
			return errors.Wrap(fReader.Close(), "failed to close the chunk file")
		}(reader)

		chunks[index] = reader
	}

	chunkWriter := chunk.NewFileWriterPtr(ptrFileOut, sizeChunk)
	mergeSorter := chunk.NewMergeSorter(chunks, chunkWriter)

	return errors.Wrap(mergeSorter.Sort(), "failed to merge sort the chunk files")
}
