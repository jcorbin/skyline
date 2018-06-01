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
			/* 0 2 4 6
			 *   |-|
			 *   | |
			 * __| |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 4, Y: 3},
				{X: 4, Y: 0},
			},
		},

		{
			name: "twin towers",
			/* 0 2 4 6 8 a
			 *   |-| |-|
			 *   | | | |
			 * __| |_| |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 3},
				{Sides: [2]int{6, 8}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 4, Y: 3},
				{X: 4, Y: 0},

				{X: 6, Y: 0},
				{X: 6, Y: 3},
				{X: 8, Y: 3},
				{X: 8, Y: 0},
			},
		},

		{
			name: "joined towers",
			/* 0 2 4 6 8 a c e
			 *   |---| |---|
			 *   |   | |   |
			 *   |   | |   |
			 *   | ..|-|.. |
			 * __| .     . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 5},
				{Sides: [2]int{8, 12}, Height: 5},
				{Sides: [2]int{4, 10}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 5},
				{X: 6, Y: 5},

				{X: 6, Y: 3},
				{X: 8, Y: 3},

				{X: 8, Y: 5},
				{X: 12, Y: 5},
				{X: 12, Y: 0},
			},
		},

		{
			name: "L",
			/*
			 * 0 2 4 6 8 a c e
			 *   |---|
			 *   | ..|__
			 * __| . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 3},
				{Sides: [2]int{4, 8}, Height: 1},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 6, Y: 3},

				{X: 6, Y: 1},

				{X: 8, Y: 1},
				{X: 8, Y: 0},
			},
		},

		{
			name: "stair",
			/*
			 * 0 2 4 6 8 a c e
			 *   |---|
			 *   |   |
			 *   | ..|---|
			 *   | . . ..|__
			 * __| . . . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 5},
				{Sides: [2]int{4, 10}, Height: 3},
				{Sides: [2]int{8, 12}, Height: 1},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 5},
				{X: 6, Y: 5},

				{X: 6, Y: 3},
				{X: 10, Y: 3},
				{X: 10, Y: 1},

				{X: 12, Y: 1},
				{X: 12, Y: 0},
			},
		},

		{
			name: "mirror stair",
			/*
			 * 0 2 4 6 8 a c e
			 *         |---|
			 *         |   |
			 *     |---|.. |
			 *   __|.. . . |
			 * __| . . . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 1},
				{Sides: [2]int{4, 10}, Height: 3},
				{Sides: [2]int{8, 12}, Height: 5},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 4, Y: 1},

				{X: 4, Y: 3},

				{X: 8, Y: 3},
				{X: 8, Y: 5},

				{X: 12, Y: 5},
				{X: 12, Y: 0},
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
