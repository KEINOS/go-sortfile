package datasize

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
)

// memoryGet is a copy of memory.Get function to ease testing.
var MemoryGet = memory.Get

// AvailableMemory returns the amount of current available free memory.
//
// It will error if it fails to get the memory information. Mostly on platforms
// such as NetBSD and OpenBSD.
func AvailableMemory() (InBytes, error) {
	mem, err := MemoryGet()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get memory information")
	}

	return InBytes(mem.Free), nil
}

// MustAvailableMemory is the same as AvailableMemory but panics if it fails to
// get the memory information.
func MustAvailableMemory() InBytes {
	mem, err := AvailableMemory()
	if err != nil {
		panic(err)
	}

	return mem
}
