package chunk

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
)

// This example shows how to use MergeSorter to merge-sort the sorted chunk files.
//
// To generate the sorted chunk files, use Lines object.
func ExampleMergeSorter() {
	// Prepare sorted chunk of files
	chunks := make([]*FileReader, 3)

	for index, pathFile := range []string{
		filepath.Join("..", "testdata", "sorted_chunks", "chunk1_sorted.txt"),
		filepath.Join("..", "testdata", "sorted_chunks", "chunk2_sorted.txt"),
		filepath.Join("..", "testdata", "sorted_chunks", "chunk3_sorted.txt"),
	} {
		// Create the chunk file
		reader, err := NewFileReader(pathFile)
		if err != nil {
			log.Fatal(err)
		}

		defer func(fReader *FileReader) {
			if err := fReader.Close(); err != nil {
				log.Fatal(err)
			}
		}(reader)

		chunks[index] = reader
	}

	// Prepare the output file
	pathFileOut := filepath.Join(os.TempDir(), "example_merge_sorter.out.txt")

	defer func() {
		_ = os.Remove(pathFileOut) // clean up the output file after the test
	}()

	sizeBufMax := datasize.InBytes(128)

	fWriter, err := NewFileWriterPath(pathFileOut, sizeBufMax)
	if err != nil {
		log.Fatal(err)
	}

	// Create the merge sorter
	mergeSorter := NewMergeSorter(chunks, fWriter)

	// Sort and write to the output file
	if err := mergeSorter.Sort(); err != nil {
		log.Fatal(err)
	}

	// Close the output file
	if err := fWriter.Close(); err != nil {
		log.Fatal(err)
	}

	// Print the sorted output file
	outData, err := os.ReadFile(pathFileOut)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(outData))
	// Output:
	// Alice
	// Bob
	// Carol
	// Charlie
	// Dave
	// Ellen
	// Eve
	// Frank
	// Isaac
	// Ivan
	// Justin
	// Mallet
	// Mallory
	// Marvin
	// Matilda
	// Oscar
	// Pat
	// Peggy
	// Steve
	// Trent
	// Trudy
	// Victor
	// Walter
	// Zoe
}
