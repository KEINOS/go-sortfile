package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/KEINOS/go-sortfile/sortfile/chunk"
	"github.com/KEINOS/go-sortfile/sortfile/datasize"
)

var (
	sizeDataMax = 1000000000 // 1GB
	sizeLine    = 30         // 30 bytes per line
)

func panicOnError(err error) {
	if err != nil {
		log.Panic(err) // panic to call defer
	}
}

// It will generate a huge file for benchmarking.
func main() {
	nameFileOut := "shuffled_huge.txt"
	pathFileOut := filepath.Join("testdata", nameFileOut)

	// Use current free memory to buffer before writing to disk
	sizeMemFree, err := datasize.AvailableMemory()
	panicOnError(err)

	fmt.Println("Generating a huge file for benchmarking with meomry size:", sizeMemFree.String())

	fWriter, err := chunk.NewFileWriter(pathFileOut, sizeMemFree)
	panicOnError(err)

	defer func() {
		panicOnError(fWriter.Close())
		fmt.Println("Done.")
	}()

	sizeWritten := 0

	for {
		line := RandString(sizeLine)
		lenLine := len(line)

		if sizeWritten > sizeDataMax {
			panicOnError(fWriter.Done())

			break
		}

		_, err := fWriter.WriteLine(line)
		panicOnError(err)

		sizeWritten += lenLine
		if sizeWritten%1000000 == 0 {
			fmt.Printf("\rWritten: %.02f%%", float64(sizeWritten)/float64(sizeDataMax)*100)
		}
	}

	fmt.Printf("\rWritten: %.02f%%\n", float64(sizeWritten)/float64(sizeWritten)*100)
}

// Letters to be used in the random string
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString generates a random string of the given byte length.
// Ref: https://stackoverflow.com/a/31832326/18152508
func RandString(lenOut int) string {
	byteOut := make([]byte, lenOut)
	src := rand.NewSource(time.Now().UnixNano())

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for index, cache, remain := lenOut-1, src.Int63(), letterIdxMax; index >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			byteOut[index] = letterBytes[idx]
			index--
		}

		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&byteOut))
}
