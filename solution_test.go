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

		{
			name: "single in the middle",
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 5},
			},
			expected: []image.Point{
				{X: 0, Y: 0},
				{X: 2, Y: 0},
				{X: 2, Y: 5},
				{X: 4, Y: 5},
				{X: 4, Y: 0},
			},
		},

		{
			name: "twin towers",
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 10},
				{Sides: [2]int{6, 8}, Height: 10},
			},
			expected: []image.Point{
				{X: 0, Y: 0},

				{X: 2, Y: 0},
				{X: 2, Y: 10},
				{X: 4, Y: 10},
				{X: 4, Y: 0},

				{X: 6, Y: 0},
				{X: 6, Y: 10},
				{X: 8, Y: 10},
				{X: 8, Y: 0},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			points, err := Solve(append([]internal.Building(nil), tc.data...))
			require.NoError(t, err, "expected Solve() to not fail")
			assert.Equal(t, tc.expected, points, "expected output points")
		})
	}
}
