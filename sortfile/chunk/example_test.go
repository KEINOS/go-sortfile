package chunk_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-sortfile/sortfile/chunk"
	"github.com/KEINOS/go-sortfile/sortfile/datasize"
)

// ----------------------------------------------------------------------------
//  chunk.FileReader()
// ----------------------------------------------------------------------------

func ExampleFileReader() {
	pathFileTarget := filepath.Join("..", "testdata", "small_chunk.txt")

	fReader, err := chunk.NewFileReader(pathFileTarget)
	if err != nil {
		log.Fatal(err)
	}

	defer fReader.Close()

	for {
		// Read the next line from the chunk file.
		// In K-way merge sort, this will be called if the line is used.
		if err := fReader.NextLine(); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("EOF reached")

				break
			}

			log.Fatal(err)
		}

		// Print the current line twice.
		// Note that the returned line is the same and does not have the traling
		// line break.
		fmt.Printf("1st call: %#v\n", fReader.CurrentLine())
		fmt.Printf("2nd call: %#v\n", fReader.CurrentLine())
	}
	// Output:
	// 1st call: "foo line"
	// 2nd call: "foo line"
	// 1st call: "bar line"
	// 2nd call: "bar line"
	// 1st call: "baz line"
	// 2nd call: "baz line"
	// EOF reached
}

// ----------------------------------------------------------------------------
//  chunk.FileSplit()
// ----------------------------------------------------------------------------

func ExampleFileSplit() {
	// Shuffled lines as a test file.
	pathFileIn := filepath.Join("..", "testdata", "sorted_chunks", "input_shuffled.txt")

	// Divide the lines of a file into 32 bytes and store in a temporary file.
	// Each file is also sorted by lines using the default isLess function (nil).
	chunks, err := chunk.FileSplit(pathFileIn, 32, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup the temporary chunk files after the test.
	defer func() {
		for _, chunk := range chunks {
			_ = os.Remove(chunk)
		}
	}()

	// Print the contents of each chunk file.
	for index, chunk := range chunks {
		fmt.Printf("Chunk file #%d\n", index+1)

		data, err := os.ReadFile(chunk)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))
	}
	// Output:
	// Chunk file #1
	// Bob
	// Carol
	// Eve
	// Peggy
	// Walter
	//
	// Chunk file #2
	// Mallet
	// Mallory
	// Matilda
	// Zoe
	//
	// Chunk file #3
	// Dave
	// Ellen
	// Frank
	// Ivan
	// Marvin
	//
	// Chunk file #4
	// Alice
	// Pat
	// Steve
	// Victor
}

// ----------------------------------------------------------------------------
//  chunk.Lines()
// ----------------------------------------------------------------------------

func ExampleLines() {
	const maxBytes = 32 // Max size of the chunk in bytes.

	// Instantiate a new chunk to be written to the output.
	chunk := chunk.Lines{}

	for _, line := range []string{
		"foo line",
		"bar line",
		"baz line",
		"hoge line", // this will overflow the maxSize so will not be appended.
	} {
		// Append only if the line will not overflow the maxSize by checking the
		// size before appending the line.
		if !chunk.WillOverSize(line, maxBytes) {
			// It will append the line with a line break at the end of the line.
			chunk.AppendLine(line)
			// Print current size of the chunk to be written to the output.
			fmt.Println("Size:", chunk.Size())
		}
	}

	// Recalculate the size of the chunk.
	fmt.Println("Size cached:", chunk.Size(), "Size caltulate:", chunk.SizeRaw())

	// Lines method returns the lines in the chunk (a slice of string) as is.
	fmt.Printf("%#v\n", chunk.Lines()[0])
	fmt.Printf("%#v\n", chunk.Lines()[1])
	fmt.Printf("%#v\n", chunk.Lines()[2])

	// Dump the chunk sorted by the lines to the io.Writer.
	var buff bytes.Buffer

	chunk.WriteSortedLines(&buff)
	fmt.Printf("Written: %#v\n", buff.String())

	// Output:
	// Size: 9
	// Size: 18
	// Size: 27
	// Size cached: 27 Size caltulate: 27
	// "foo line\n"
	// "bar line\n"
	// "baz line\n"
	// Written: "bar line\nbaz line\nfoo line\n"
}

func ExampleLines_custom_sort() {
	const maxBytes = 32 // Max size of the chunk in bytes.

	// User defined sort function. Which is the reverse of the default sort.
	isLess := func(a, b string) bool {
		return a > b
	}

	// Instantiate a new chunk to be written to the output.
	chunk := chunk.Lines{}

	// Assign the custom sort function.
	chunk.IsLess = isLess

	// The below will use "alice", "bob", and "david" lines for the chunk.
	for _, line := range []string{
		"alice line",
		"bob line",
		"charlie line", // this will overflow the maxSize so will not be appended.
		"david line",
	} {
		// Append only if the line will not overflow the maxSize by checking the
		// size before appending the line.
		if !chunk.WillOverSize(line, maxBytes) {
			// It will append the line with a line break at the end of the line.
			chunk.AppendLine(line)
			// Print current size of the chunk to be written to the output.
			fmt.Printf("'%s' appended. Current chunk size: %v\n", line, chunk.Size())
		} else {
			fmt.Printf("Skip: '%s' will overflow. Current chunk size: %v\n", line, chunk.Size())
		}
	}

	// Recalculate the size of the chunk.
	fmt.Println("Size cached:", chunk.Size(), "Size caltulate:", chunk.SizeRaw())

	// Lines method returns the lines in the chunk (a slice of string) as is.
	fmt.Printf("%#v\n", chunk.Lines()[0])
	fmt.Printf("%#v\n", chunk.Lines()[1])
	fmt.Printf("%#v\n", chunk.Lines()[2])

	// Dump the chunk sorted by the lines to the io.Writer.
	var buff bytes.Buffer

	chunk.WriteSortedLines(&buff)
	fmt.Printf("Written: %#v\n", buff.String())

	// Output:
	// 'alice line' appended. Current chunk size: 11
	// 'bob line' appended. Current chunk size: 20
	// Skip: 'charlie line' will overflow. Current chunk size: 20
	// 'david line' appended. Current chunk size: 31
	// Size cached: 31 Size caltulate: 31
	// "alice line\n"
	// "bob line\n"
	// "david line\n"
	// Written: "david line\nbob line\nalice line\n"
}

// ----------------------------------------------------------------------------
//  chunk.NewFileWriter()
// ----------------------------------------------------------------------------

func ExampleNewFileWriter() {
	// Preapare temp file to write to.
	pathFileTemp := filepath.Join(os.TempDir(), "ExampleFileWriter.txt")
	// Cleanup temp file after the test.
	defer func() {
		if err := os.Remove(pathFileTemp); err != nil {
			log.Fatal(err)
		}
	}()

	// The size of the buffer. It will write to the file when the buffer reaches
	// 128 bytes each time.
	sizeBufMax := datasize.InBytes(128)

	// Instantiate a new file writer.
	fWriter, err := chunk.NewFileWriter(pathFileTemp, sizeBufMax)
	if err != nil {
		log.Fatal(err)
	}

	defer fWriter.Close()

	for _, line := range []string{
		"智恵子は東京に空が無いといふ、",
		"ほんとの空が見たいといふ。",
		"私は驚いて空を見る。",
		"桜若葉の間に在るのは、",
		"切つても切れない",
		"むかしなじみのきれいな空だ。",
		"どんよりけむる地平のぼかしは",
		"うすもも色の朝のしめりだ。",
		"智恵子は遠くを見ながら言ふ。",
		"阿多多羅山あたたらやまの山の上に",
		"毎日出てゐる青い空が",
		"智恵子のほんとの空だといふ。",
		"あどけない空の話である。",
	} {
		// Let the writer write the line whether it buffers or not.
		if _, err := fWriter.WriteLine(line); err != nil {
			log.Fatal(err)
		}
	}

	// Flash the remaining buffer to the file.
	if err := fWriter.Done(); err != nil {
		log.Fatal(err)
	}

	dataWritten, err := os.ReadFile(pathFileTemp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(dataWritten))
	// Output:
	// 智恵子は東京に空が無いといふ、
	// ほんとの空が見たいといふ。
	// 私は驚いて空を見る。
	// 桜若葉の間に在るのは、
	// 切つても切れない
	// むかしなじみのきれいな空だ。
	// どんよりけむる地平のぼかしは
	// うすもも色の朝のしめりだ。
	// 智恵子は遠くを見ながら言ふ。
	// 阿多多羅山あたたらやまの山の上に
	// 毎日出てゐる青い空が
	// 智恵子のほんとの空だといふ。
	// あどけない空の話である。
}

// ----------------------------------------------------------------------------
//  chunk.NewMergeSorter()
// ----------------------------------------------------------------------------

// This example shows how to use MergeSorter to merge-sort the sorted chunk files.
//
// To generate the sorted chunk files, use Lines object.
func ExampleMergeSorter() {
	// Prepare sorted chunk of files
	chunks := make([]*chunk.FileReader, 3)

	for index, pathFile := range []string{
		filepath.Join("..", "testdata", "sorted_chunks", "chunk1_sorted.txt"),
		filepath.Join("..", "testdata", "sorted_chunks", "chunk2_sorted.txt"),
		filepath.Join("..", "testdata", "sorted_chunks", "chunk3_sorted.txt"),
	} {
		// Create the chunk file
		reader, err := chunk.NewFileReader(pathFile)
		if err != nil {
			log.Fatal(err)
		}

		defer func(fReader *chunk.FileReader) {
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

	fWriter, err := chunk.NewFileWriter(pathFileOut, sizeBufMax)
	if err != nil {
		log.Fatal(err)
	}

	// Create the merge sorter
	mergeSorter := chunk.NewMergeSorter(chunks, fWriter)

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
