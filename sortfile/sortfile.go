/*
Package sortfile provides functions to sort a file. Both in-memory and external merge sort.

## Usage

```go
import "github.com/KEINOS/go-sortfile/sortfile"
```
*/
package sortfile

const (
	LF   = "\n"    // LF is the line feed character
	CR   = "\r"    // CR is the carriage return character
	CRLF = CR + LF // CRLF is the carriage return and line feed character
)
