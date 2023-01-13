package sortfile

import (
	"path/filepath"
	"testing"
)

func BenchmarkFromPath(b *testing.B) {
	pathFileIn := filepath.Join("testdata", "shuffled_huge.txt")
	pathFileOut := filepath.Join(b.TempDir(), b.Name()+".txt")
	forceExternalSort := true

	if !FileExists(pathFileIn) {
		b.Skip("the huge file for benchmarking does not exist. Run 'go generate ./...' to generate it first.")
	}

	b.ResetTimer()

	err := FromPath(pathFileIn, pathFileOut, forceExternalSort)
	if err != nil {
		b.Fatal(err)
	}
}
