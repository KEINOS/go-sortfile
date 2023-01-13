package datasize

import (
	"fmt"
)

// ============================================================================
//  Constants of common sizes
// ============================================================================

const (
	KiB = InBytes(1024) // KiB = 1024 Bytes
	MiB = 1024 * KiB    // MiB = 1024 * KiB
	GiB = 1024 * MiB    // GiB = 1024 * MiB
	TiB = 1024 * GiB    // TiB = 1024 * GiB
	PiB = 1024 * TiB    // PiB = 1024 * TiB
)

// ============================================================================
//  Type: InBytes
// ============================================================================

// InBytes represents a size in bytes. It is a wrapper of uint64 to provide
// human readable output.
type InBytes uint64

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new InBytes instance from the given input. The input can be
// any of the following types:
//
//	int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
//	float32, float64, rune and InBytes
//
// Note that the input will be converted to uint64 type. Thus, this is not suitable
// for negative numbers.
// If the input is not one of the above types, it returns 0.
func New(input any) InBytes {
	switch value := input.(type) {
	case int:
		return InBytes(uint64(value))
	case int8:
		return InBytes(uint64(value))
	case int16:
		return InBytes(uint64(value))
	case int32:
		return InBytes(uint64(value))
	case int64:
		return InBytes(uint64(value))
	case uint:
		return InBytes(uint64(value))
	case uint8:
		return InBytes(uint64(value))
	case uint16:
		return InBytes(uint64(value))
	case uint32:
		return InBytes(uint64(value))
	case uint64:
		return InBytes(uint64(value))
	case float32:
		return InBytes(uint64(value))
	case float64:
		return InBytes(uint64(value))
	case InBytes:
		return input.(InBytes)
	default:
		return InBytes(0)
	}
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// InBytes returns the size in bytes unit without the unit name. The given argument
// is the number of digits after the decimal point.
func (size InBytes) InBytes() string {
	return fmt.Sprintf("%d", uint64(size))
}

// InGiB returns the size in GiB without the unit name. The given argument is the
// number of digits after the decimal point.
func (size InBytes) InGiB(DigDecPoint int) string {
	tpl := fmt.Sprintf("%%.%df", DigDecPoint)

	return fmt.Sprintf(tpl, float64(size)/float64(GiB))
}

// InKiB returns the size in KiB without the unit name. The given argument is the
// number of digits after the decimal point.
func (size InBytes) InKiB(DigDecPoint int) string {
	tpl := fmt.Sprintf("%%.%df", DigDecPoint)

	return fmt.Sprintf(tpl, float64(size)/float64(KiB))
}

// InMiB returns the size in MiB without the unit name. The given argument is the
// number of digits after the decimal point.
func (size InBytes) InMiB(DigDecPoint int) string {
	tpl := fmt.Sprintf("%%.%df", DigDecPoint)

	return fmt.Sprintf(tpl, float64(size)/float64(MiB))
}

// InPiB returns the size in PiB without the unit name. The given argument is the
// number of digits after the decimal point.
func (size InBytes) InPiB(DigDecPoint int) string {
	tpl := fmt.Sprintf("%%.%df", DigDecPoint)

	return fmt.Sprintf(tpl, float64(size)/float64(PiB))
}

// InTiB returns the size in TiB without the unit name. The given argument is the
// number of digits after the decimal point.
func (size InBytes) InTiB(DigDecPoint int) string {
	tpl := fmt.Sprintf("%%.%df", DigDecPoint)

	return fmt.Sprintf(tpl, float64(size)/float64(TiB))
}

// IsEqualTo returns true if the given size is equal to the current size.
func (size InBytes) IsEqualTo(other InBytes) bool {
	return size == other
}

// IsSmallerThan returns true if the given size is smaller than the current size.
func (size InBytes) IsSmallerThan(other InBytes) bool {
	return size < other
}

// IsGreaterThan returns true if the given size is greater than the current size.
func (size InBytes) IsGreaterThan(other InBytes) bool {
	return size > other
}

// Round returns the size rounded to the nearest unit in InBytes type.
// This method is useful to find the nice round number to chunk the data.
//
// E.g.
//
//	1025 will be rounded to 1 KiB(1024) thus the return value is `InBytes(1024)`.
func (s InBytes) Round() (InBytes, string) {
	size := InBytes(s)

	switch {
	case size >= PiB:
		return (size / PiB) * PiB, "PiB"
	case size >= TiB:
		return (size / TiB) * TiB, "TiB"
	case size >= GiB:
		return (size / GiB) * GiB, "GiB"
	case size >= MiB:
		return (size / MiB) * MiB, "MiB"
	case size >= KiB:
		return (size / KiB) * KiB, "KiB"
	default:
		return size, "Bytes"
	}
}

// String is the stringer implementation for InBytes. It will return the size
// in the nearest unit.
func (size InBytes) String() string {
	switch {
	case size >= PiB:
		return fmt.Sprintf("%.2f PiB", float64(size)/float64(PiB))
	case size >= TiB:
		return fmt.Sprintf("%.2f TiB", float64(size)/float64(TiB))
	case size >= GiB:
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(GiB))
	case size >= MiB:
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(MiB))
	case size >= KiB:
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(KiB))
	default:
		return fmt.Sprintf("%d Bytes", uint64(size))
	}
}
