package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// an example function that contains a bug (fails for negative num)
func AddInt(x, y int) int {
	for i := 0; i < x; i++ {
		y += 1
	}
	return y
}

func FuzzAddInt(f *testing.F) {
	//optional, this is a starting point for fuzzing process (first pair is x=0, y=1)
	f.Add(0, 1)
	f.Fuzz(func(t *testing.T, x int, y int) {
		result := AddInt(x, y)
		require.Equal(t, x+y, result)
	})
}
