/*
Package sortfile provides functions to sort a file. Both in-memory and external merge sort.

## Usage

```go
import "github.com/KEINOS/go-sortfile/sortfile"
```
*/
package sortfile

import (
	"github.com/KEINOS/go-sortfile/sortfile/chunk"
)

const (
	LF   = chunk.LF   // LF is the line feed character
	CR   = chunk.CR   // CR is the carriage return character
	CRLF = chunk.CRLF // CRLF is the carriage return and line feed character
)

var GO_EOL = LF // GO_EOL is the end of line character for the current OS

func init() {
	// Set the end of line character for the current OS
	GO_EOL = chunk.GO_EOL
}
