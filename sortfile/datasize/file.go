package datasize

import (
	"os"

	"github.com/KEINOS/go-countline/cl"
	"github.com/pkg/errors"
)

// OsOpen is a copy of os.Open to ease testing.
var OsOpen = os.Open

// File returns the data size of the given file and the number of lines.
func File(path string) (InBytes, int, error) {
	file, err := OsOpen(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to open file")
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get file stat")
	}

	numLines, err := cl.CountLines(file)

	return InBytes(stat.Size()), numLines, errors.Wrap(err, "failed to count lines in file")
}
