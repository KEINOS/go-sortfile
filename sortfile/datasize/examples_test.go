package datasize_test

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/KEINOS/go-sortfile/sortfile/datasize"
)

// ============================================================================
//  File
// ============================================================================

func ExampleFile() {
	pathFile := filepath.Join("..", "testdata", "example.txt")

	sizeFile, numLines, err := datasize.File(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File size:", sizeFile.String())
	fmt.Println("Num lines:", numLines)
	// Output:
	// File size: 27 Bytes
	// Num lines: 3
}

// ============================================================================
//  InBytes Type
// ============================================================================

func ExampleNew() {
	// Supported types
	for index, input := range []interface{}{
		int(65),
		int8(65),
		int16(65),
		int32(65),
		int64(65),
		uint(65),
		uint8(65),
		uint16(65),
		uint32(65),
		uint64(65),
		float32(65.0),
		float64(65.0),
		rune('A'), // ASCII 65
		datasize.InBytes(65),
	} {
		size := datasize.New(input)

		expect := "datasize.InBytes: 65"
		actual := fmt.Sprintf("%T: %d", size, size)

		if expect != actual {
			log.Fatalf("Test #%d failed. Input: %#v, Expected: %s, Actual: %s",
				index+1, input, expect, actual)
		}
	}
	fmt.Println("OK")

	// Unsupported type returns 0
	size := datasize.New("65") // string
	fmt.Printf("%T: %d", size, size)
	// Output:
	// OK
	// datasize.InBytes: 0
}

func ExampleInBytes_print_variations() {
	var size datasize.InBytes = 1234567890 // automatic cast to uint64

	// InBytes has a Stringer implementation and prints the size with the nearest
	// unit including the unit name. The digits after the decimal point are 2.
	fmt.Println(size)
	// Raw value in decimal.
	fmt.Printf("%d\n", size)
	// Print the value in defined units without the unit name. The given argument
	// is the number of digits after the decimal point.
	fmt.Println(size.InPiB(6), "PiB")
	fmt.Println(size.InTiB(5), "TiB")
	fmt.Println(size.InGiB(4), "GiB")
	fmt.Println(size.InMiB(3), "MiB")
	fmt.Println(size.InKiB(2), "KiB")
	fmt.Println(size.InBytes(), "Bytes")
	// Output:
	// 1.15 GiB
	// 1234567890
	// 0.000001 PiB
	// 0.00112 TiB
	// 1.1498 GiB
	// 1177.376 MiB
	// 1205632.71 KiB
	// 1234567890 Bytes
}

// ----------------------------------------------------------------------------
//  InBytes.IsEqualTo()
// ----------------------------------------------------------------------------

func ExampleInBytes_IsEqualTo() {
	var data datasize.InBytes = 100

	// negative values are not allowed
	fmt.Println(data.IsEqualTo(10))   // 100 == 10 = false
	fmt.Println(data.IsEqualTo(100))  // 100 == 100 = true
	fmt.Println(data.IsEqualTo(1000)) // 100 == 1000 = false
	// Output:
	// false
	// true
	// false
}

// ----------------------------------------------------------------------------
//  InBytes.IsGreaterThan()
// ----------------------------------------------------------------------------

func ExampleInBytes_IsGreaterThan() {
	var data datasize.InBytes = 100

	// negative values are not allowed
	fmt.Println(data.IsGreaterThan(10))   // 100 > 10 = true
	fmt.Println(data.IsGreaterThan(100))  // 100 > 100 = false
	fmt.Println(data.IsGreaterThan(1000)) // 100 > 1000 = false
	// Output:
	// true
	// false
	// false
}

// ----------------------------------------------------------------------------
//  InBytes.IsSmallerThan()
// ----------------------------------------------------------------------------

func ExampleInBytes_IsSmallerThan() {
	var data datasize.InBytes = 100

	// negative values are not allowed
	fmt.Println(data.IsSmallerThan(10))   // 100 < 10 = false
	fmt.Println(data.IsSmallerThan(100))  // 100 < 100 = false
	fmt.Println(data.IsSmallerThan(1000)) // 100 < 1000 = true
	// Output:
	// false
	// false
	// true
}

// ----------------------------------------------------------------------------
//  InBytes.Round()
// ----------------------------------------------------------------------------

func ExampleInBytes_Round() {
	// Round method rounds the value to the nearest unit.
	{
		// 1024 + 1 ∴ 1 KiB(1024) * 1 is the nearest unit
		size := datasize.InBytes(1025)

		rounded, unit := size.Round()
		fmt.Printf("origin: %d --> rounded: %d unit: %s\n", size, rounded, unit)
	}
	{
		// 1024 *2 + 2 ∴ 1 KiB(1024) * 2 is the nearest unit
		size := datasize.InBytes(2050)

		rounded, unit := size.Round()
		fmt.Printf("origin: %d --> rounded: %d unit: %s\n", size, rounded, unit)
	}
	{
		// 1024*1024 + 1 ∴ 1 MiB(1024*1024) * 1 is the nearest unit
		size := datasize.InBytes(1048577)

		// Round the value to the nearest unit.
		rounded, unit := size.Round()
		fmt.Printf("origin: %d --> rounded: %d unit: %s\n", size, rounded, unit)
	}
	// Output:
	// origin: 1025 --> rounded: 1024 unit: KiB
	// origin: 2050 --> rounded: 2048 unit: KiB
	// origin: 1048577 --> rounded: 1048576 unit: MiB
}
