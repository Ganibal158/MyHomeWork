package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHW2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "dff4a", expected: "dfffffa"},
		{input: "oka2b2c0", expected: "okaabb"},
		{input: "a\n2", expected: "a\n\n"},
	}

	for _, val := range tests {
		t.Run(val.input, func(t *testing.T) {
			result, err := Unpack(val.input)
			require.NoError(t, err)
			require.Equal(t, val.expected, result)
		})
	}
}

func TestHw2Err(t *testing.T) {
	errVal := []string{"1vgh:)", "578", "5erhreh", "gkjgkgg87"}
	for _, val := range errVal {
		t.Run(val, func(t *testing.T) {
			_, err := Unpack(val)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
