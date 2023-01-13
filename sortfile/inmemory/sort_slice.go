/*
Package inmemory provides sorting algorithms for in-memory data.
*/
package inmemory

import "golang.org/x/exp/slices"

// SortSlice sorts the given slice of strings. The given input will be modified after
// its call.
func SortSlice(input []string) {
	// Current implementation is using slices.Sort. See benchmark_test.go for benchmark
	// between other sorting algorithms.
	slices.Sort(input)
}