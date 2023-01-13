package inmemory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortSlice(t *testing.T) {
	require.True(t, preTestSortFunction(SortSlice))
}

func preTestSortFunction(fnTest func([]string)) bool {
	input := []string{
		"foo",
		"bar",
		"baz",
	}

	fnTest(input)

	for index, value := range []string{
		"bar",
		"baz",
		"foo",
	} {
		if input[index] != value {
			return false
		}
	}

	return true
}
