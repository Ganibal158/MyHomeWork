package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.

var str1 = `    aaa aaa bbb aaa  bbb aaa ddd ddd ddd xxx nnn  ppp zzzz  aaa ccc `

func TestHW2(t *testing.T) {
	t.Run("empty string test", func(t *testing.T) {
		require.Len(t, Top10("      "), 0)
	})

	t.Run("successful test", func(t *testing.T) {
		expected := []string{
			"aaa",  // 4
			"ddd",  // 3
			"bbb",  // 2
			"ccc",  // 1
			"nnn",  // 1
			"ppp",  // 1
			"xxx",  // 1
			"zzzz", // 1
		}
		require.Equal(t, expected, Top10(str1))
	})
}
