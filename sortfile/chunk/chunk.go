package chunk

import "runtime"

const (
	LF   = "\n"    // LF is the line feed character
	CR   = "\r"    // CR is the carriage return character
	CRLF = CR + LF // CRLF is the carriage return and line feed character
)

var GO_EOL = LF // GO_EOL is the end of line character for the current OS

func init() {
	// Set the end of line character for the current OS
	GO_EOL = CRLF

	if runtime.GOOS != "windows" {
		GO_EOL = LF
	}
}
