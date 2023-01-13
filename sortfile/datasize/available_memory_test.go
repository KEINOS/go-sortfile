package datasize

import (
	"testing"

	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMustAvailableMemory_goledn(t *testing.T) {
	require.NotPanics(t,
		func() {
			_ = MustAvailableMemory()
		},
		"it should not panic on non FreeBSD, OpenBSD platforms",
	)
}

func TestMustAvailableMemory_failed_to_get_memory_info(t *testing.T) {
	oldMemoryGet := MemoryGet
	defer func() {
		MemoryGet = oldMemoryGet
	}()

	MemoryGet = func() (*memory.Stats, error) {
		return nil, errors.New("forced error")
	}

	require.PanicsWithError(t,
		"failed to get memory information: forced error",
		func() {
			_ = MustAvailableMemory()
		},
		"it should contain the reason of the error on panic",
	)
}
