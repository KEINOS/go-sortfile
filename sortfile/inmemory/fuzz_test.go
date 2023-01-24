package inmemory

import (
	"testing"
)

func FuzzSortSlice(f *testing.F) {
	for _, testCase := range listNumItems {
		f.Add(testCase)
	}

	f.Fuzz(func(t *testing.T, numItems int) {
		lines := randSlice(numItems)
		SortSlice(lines)
	})
}
