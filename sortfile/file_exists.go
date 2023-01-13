package sortfile

import "os"

// FileExists returns true if the path exists and is a file.
func FileExists(pathFile string) bool {
	info, err := os.Stat(pathFile)
	if err == nil {
		return !info.IsDir()
	}

	return false
}
