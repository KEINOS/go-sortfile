package sortfile

import (
	"bufio"
	"io"
	"strings"

	"github.com/KEINOS/go-sortfile/sortfile/inmemory"
	"github.com/pkg/errors"
)

// InMemory sorts the lines in-memory from the given io.Reader and writes the
// result to the given io.Writer.
// Note that the number of lines is required to be known in advance.
//
// Usually it is recommended to use the FromPath() function which detects
// whether to use the in-memory sort or the external merge sort.
func InMemory(numLines int, input io.Reader, output io.Writer, isLess func(string, string) bool) error {
	lines := make([]string, numLines)
	scanner := bufio.NewScanner(input)
	index := 0

	for scanner.Scan() {
		lines[index] = scanner.Text() + LF
		index++
	}

	if isLess == nil {
		inmemory.SortSlice(lines)
	} else {
		inmemory.SortSliceFunc(lines, isLess)
	}

	_, err := output.Write([]byte(strings.Join(lines, "")))

	return errors.Wrap(err, "failed to write to output")
}
