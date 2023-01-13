package inmemory_test

import (
	"fmt"

	"github.com/Code-Hex/dd"
	"github.com/KEINOS/go-sortfile/sortfile/inmemory"
)

// ============================================================================
//  SortSlice()
// ============================================================================

func ExampleSortSlice() {
	lines := []string{
		"foo",
		"bar",
		"baz",
	}

	// Sort slices of strings.
	// For benchmark between other algorithms, see the benchmark_test.go.
	inmemory.SortSlice(lines)

	fmt.Println(dd.Dump(lines))
	// Output:
	// []string{
	//   "bar",
	//   "baz",
	//   "foo",
	// }
}
