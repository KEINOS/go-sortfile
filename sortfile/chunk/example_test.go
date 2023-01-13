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

func ExampleFileSplit() {
	// Shuffled lines as a test file.
	pathFileIn := filepath.Join("..", "testdata", "sorted_chunks", "input_shuffled.txt")

	// Divide the lines of a file into 32 bytes and store in a temporary file.
	// Each file is also sorted by lines.
	chunks, err := chunk.FileSplit(pathFileIn, 32)
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

func ExampleNewFileWriterPath() {
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
	fWriter, err := chunk.NewFileWriterPath(pathFileTemp, sizeBufMax)
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
