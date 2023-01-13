package inmemory

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yourbasic/radix"
	"golang.org/x/exp/slices"
)

// Various sorting algorithms to compare. We use map to pick the algorithm randomly.
var algorithmsSort = map[string]func([]string){
	"radix.SortSlice": func(input []string) {
		radix.SortSlice(input, func(i int) string { return input[i] })
	},
	"sort.Slice": func(input []string) {
		sort.Slice(input, func(i, j int) bool { return input[i] < input[j] })
	},
	"slices.Sort": func(input []string) {
		slices.Sort(input)
	},
	"slices.SortStableFunc": func(input []string) {
		slices.SortStableFunc(input, func(i, j string) bool { return i < j })
	},
}

var listNumItems = []int{1000, 10000, 100000}

func Benchmark_various_sort_algorithm(b *testing.B) {
	for _, numItems := range listNumItems {
		for nameTest, fnTest := range algorithmsSort {
			require.True(b, preTestSortFunction(fnTest))

			nameTest := fmt.Sprintf("%s %d", nameTest, numItems)
			b.Run(nameTest, func(b *testing.B) {
				lines := randSlice(numItems)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					fnTest(lines)
				}
			})
		}
	}
}

// ----------------------------------------------------------------------------
//  Helper functions
// ----------------------------------------------------------------------------

func randSlice(numLines int) []string {
	const numChars = 100 // number of characters per line

	lines := make([]string, numLines)

	for index := range lines {
		lines[index] = randString(numChars) + "\n"
	}

	return lines
}

func randString(numChars int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lenLetters = len(letters)

	line := make([]byte, numChars)
	for i := range line {
		line[i] = letters[rand.Intn(lenLetters)]
	}

	return string(line)
}
