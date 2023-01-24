/*
Package inmemory provides sorting algorithms for in-memory data.
*/
package inmemory

import "golang.org/x/exp/slices"

// SortSlice sorts the given slice of strings ascendant. The given input will be modified
// after its call. To sort a slice of strings descendant, use SortSliceFunc.
func SortSlice(input []string) {
	// Current implementation is using slices.Sort. See benchmark_test.go for benchmark
	// between other sorting algorithms.
	slices.Sort(input)
}

// SortSliceFunc is similar to SortSlice but it takes a function to compare two strings.
// The given input will be modified after its call.
func SortSliceFunc(input []string, less func(a, b string) bool) {
	slices.SortFunc(input, less)
}
