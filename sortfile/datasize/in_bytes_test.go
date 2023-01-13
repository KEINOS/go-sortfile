package datasize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInBytes_Round(t *testing.T) {
	for _, test := range []struct {
		name         string
		expectUnit   string
		expectString string
		input        InBytes
	}{
		{
			name:         "PiB is the nearest unit",
			input:        InBytes(1024*1024*1024*1024*1024*2 + 1),
			expectString: "2.00 PiB",
			expectUnit:   "PiB",
		},
		{
			name:         "TiB is the nearest unit",
			input:        InBytes(1024*1024*1024*1024*2 + 1),
			expectString: "2.00 TiB",
			expectUnit:   "TiB",
		},
		{
			name:         "GiB is the nearest unit",
			input:        InBytes(1024*1024*1024*2 + 1),
			expectString: "2.00 GiB",
			expectUnit:   "GiB",
		},
		{
			name:         "MiB is the nearest unit",
			input:        InBytes(1024*1024*2 + 1),
			expectString: "2.00 MiB",
			expectUnit:   "MiB",
		},
		{
			name:         "KiB is the nearest unit",
			input:        InBytes(1024*2 + 1),
			expectString: "2.00 KiB",
			expectUnit:   "KiB",
		},
		{
			name:         "Byte is the nearest unit",
			input:        InBytes(3),
			expectString: "3 Bytes",
			expectUnit:   "Bytes",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actualOut, actualUnit := test.input.Round()

			assert.Equal(t, test.expectString, actualOut.String(), "unexpected rounded value and print")
			assert.Equal(t, test.expectUnit, actualUnit, "unexpected unit name")
		})
	}
}
