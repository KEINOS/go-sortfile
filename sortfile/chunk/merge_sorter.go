package chunk

import (
	"io"
	"strings"

	"github.com/KEINOS/go-donegroup/donegroup"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: MergeSorter
// ----------------------------------------------------------------------------

// MergeSorter merge-sorts the sorted chunk files.
//
// Note that each chunk file must be sorted.
type MergeSorter struct {
	outFile *FileWriter
	// IsLess is the function to compare two strings during merge-sorting. This
	// function must be the same as the one used to sort the chunk files.
	IsLess func(a, b string) bool
	chunks []*FileReader
	lenK   int
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewMergeSorter returns a new MergeSorter object. The inFiles must be a slice
// of FileReader objects and each file must be sorted.
func NewMergeSorter(inFiles []*FileReader, outFile *FileWriter) *MergeSorter {
	return &MergeSorter{
		lenK:    len(inFiles),
		outFile: outFile,
		chunks:  inFiles,
		IsLess:  IsLess,
	}
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// Sort merge-sorts the chunk files and writes the result to the output file.
func (ms *MergeSorter) Sort() error {
	// Initialize the first line of each chunk
	for indexK := 0; indexK < ms.lenK; indexK++ {
		if err := ms.chunks[indexK].NextLine(); err != nil {
			return errors.Wrap(err, "failed to read the first line during initialization")
		}
	}

	doneList, err := donegroup.New(ms.lenK)
	if err != nil {
		return errors.Wrap(err, "failed to create a new DoneGroup")
	}

	leastLine := ms.chunks[0].CurrentLine() // use the 1st line as the least line
	lastUsedIndex := 0

	for {
		if doneList.IsDoneAll() {
			break
		}

		// Find the least line in K.
		for indexK := 0; indexK < ms.lenK; indexK++ {
			if ms.chunks[indexK].IsEOF() {
				doneList.Done(indexK + 1)

				continue
			}

			if leastLine == "" {
				leastLine = ms.chunks[indexK].CurrentLine()
				lastUsedIndex = indexK

				continue
			}

			// Is current line less than the least line?
			if ms.IsLess(ms.chunks[indexK].CurrentLine(), leastLine) {
				// Update
				leastLine = ms.chunks[indexK].CurrentLine()

				lastUsedIndex = indexK
			}
		}

		// Append the least line to the output file if not empty
		if strings.TrimSpace(leastLine) != "" {
			if _, err := ms.outFile.WriteLine(leastLine); err != nil {
				return errors.Wrap(err, "failed to write the line")
			}
		}

		// Forward to the next line of the chunk used
		err := ms.chunks[lastUsedIndex].NextLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return errors.Wrap(err, "failed to read the next line")
		}

		leastLine = "" // reset
	}

	return errors.Wrap(ms.outFile.Done(), "failed to dump the remaining buffer")
}
