package sortfile_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-sortfile/sortfile"
	"github.com/KEINOS/go-sortfile/sortfile/datasize"
)

// ----------------------------------------------------------------------------
//  ExternalFile
// ----------------------------------------------------------------------------

func ExampleExternalFile() {
	exitOnError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Input and output file paths
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")

	// Get file and memory information
	sizeFileIn, _, err := datasize.File(pathFileIn)
	exitOnError(err)

	sizeMemoryFree, err := datasize.AvailableMemory()
	exitOnError(err)

	// Open the file to read
	fileIn, err := os.Open(pathFileIn)
	exitOnError(err)

	defer fileIn.Close()

	fileOut := os.Stdout

	// External merge sort with sizeMemoryFree as the chunk size. Use the default
	// sort function (by nil).
	err = sortfile.ExternalFile(sizeFileIn, sizeMemoryFree, fileIn, fileOut, nil)
	exitOnError(err)
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

// ----------------------------------------------------------------------------
//  FileExists
// ----------------------------------------------------------------------------

func ExampleFileExists() {
	for _, pathTarget := range []string{
		filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt"), // Existing file
		os.TempDir(),                // Exists but not a file
		"unknown-non-existing-file", // Not exists
	} {
		exists := sortfile.FileExists(pathTarget)

		fmt.Println("Is file:", exists, ":", pathTarget)
	}
	// Output:
	// Is file: true : testdata/sorted_chunks/input_shuffled.txt
	// Is file: false : /var/folders/8c/lmckjks95fj4h_jqzw4v3k_w0000gn/T/
	// Is file: false : unknown-non-existing-file
}

// ----------------------------------------------------------------------------
//  FromPath
// ----------------------------------------------------------------------------

func ExampleFromPath() {
	exitOnError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Input and output file paths
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")
	pathFileOut := filepath.Join(os.TempDir(), "pkg-sortfile_example_from_path.txt")

	// Clean up the output file after the test
	defer func() {
		exitOnError(os.Remove(pathFileOut))
	}()

	// Sort file in-memory since the file size is small
	forceExternalSort := false // auto detect

	err := sortfile.FromPath(pathFileIn, pathFileOut, forceExternalSort)
	exitOnError(err)

	// Print the result
	data, err := os.ReadFile(pathFileOut)
	exitOnError(err)

	fmt.Println(string(data))
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

func ExampleFromPathFunc() {
	exitOnError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Input and output file paths
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")
	pathFileOut := filepath.Join(os.TempDir(), "pkg-sortfile_example_from_path.txt")

	// Clean up the output file after the test
	defer func() {
		exitOnError(os.Remove(pathFileOut))
	}()

	// Sort file in-memory since the file size is small
	forceExternalSort := false // auto detect

	// User defined sort function (reverse sort)
	isLess := func(a, b string) bool {
		return a > b
	}

	err := sortfile.FromPathFunc(pathFileIn, pathFileOut, forceExternalSort, isLess)
	exitOnError(err)

	// Print the result
	data, err := os.ReadFile(pathFileOut)
	exitOnError(err)

	fmt.Println(string(data))
	// Output:
	// Zoe
	// Walter
	// Victor
	// Trudy
	// Trent
	// Steve
	// Peggy
	// Pat
	// Oscar
	// Matilda
	// Marvin
	// Mallory
	// Mallet
	// Justin
	// Ivan
	// Isaac
	// Frank
	// Eve
	// Ellen
	// Dave
	// Charlie
	// Carol
	// Bob
	// Alice
}

// ----------------------------------------------------------------------------
//  InMemory
// ----------------------------------------------------------------------------

// Example of using the in-memory sort.
//
// Note that the number of lines is required to be known in advance. Usually it
// is recommended to use the FromPath() function which detects whether to use
// the in-memory sort or the external merge sort.
func ExampleInMemory() {
	pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")

	// Get the number of lines in the file
	sizeFile, numLines, err := datasize.File(pathFileIn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File size:", sizeFile)

	// Open the input file
	ptrFileIn, err := os.Open(pathFileIn)
	if err != nil {
		log.Fatal(err)
	}

	defer ptrFileIn.Close()

	// Output to stdout
	ptrFileOut := os.Stdout

	// Custom sort function as reverse alphabetical order
	isLess := func(a, b string) bool {
		return a > b
	}

	// Sort the file in-memory. Use default isLess function for sorting (by nil).
	if err := sortfile.InMemory(numLines, ptrFileIn, ptrFileOut, isLess); err != nil {
		log.Fatal(err)
	}
	// Output:
	// File size: 145 Bytes
	// Zoe
	// Walter
	// Victor
	// Trudy
	// Trent
	// Steve
	// Peggy
	// Pat
	// Oscar
	// Matilda
	// Marvin
	// Mallory
	// Mallet
	// Justin
	// Ivan
	// Isaac
	// Frank
	// Eve
	// Ellen
	// Dave
	// Charlie
	// Carol
	// Bob
	// Alice
}
