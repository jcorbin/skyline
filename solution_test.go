package main_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/jcorbin/skyline"
	"github.com/jcorbin/skyline/internal"
)

func TestSolve(t *testing.T) {
	for _, tc := range []struct {
		name     string
		data     []internal.Building
		expected []image.Point
	}{
		{
			name:     "empty data",
			data:     nil,
			expected: nil,
		},
		// TODO actual test cases
	} {
		t.Run(tc.name, func(t *testing.T) {
			points, err := Solve(tc.data)
			require.NoError(t, err, "expected Solve() to not fail")
			assert.Equal(t, tc.expected, points, "expected output points")
		})
	}
}
