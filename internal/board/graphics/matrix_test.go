package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/utils"
)

const (
	testKey = "testKey"
	tag1    = "tag1"
	tag2    = "tag2"
)

func TestMatrix_RayCastVertical(t *testing.T) {
	a := require.New(t)
	x := 4
	matrix := NewMatrix(10, 10)

	tag1Y := 9
	tag2Y := 7

	matrix.Cell(utils.Vec(x, tag1Y)).Tag(testKey, tag1)
	matrix.Cell(utils.Vec(x, tag2Y)).Tag(testKey, tag2)

	hit, err := matrix.RayCastVertical(
		utils.Vec(x, 3),
		map[string]func(string) bool{
			testKey: func(value string) bool {
				return value == tag1
			},
		},
		tag1Y-3,
	)
	a.NoError(err)
	a.Equal(true, hit)

	hit, err = matrix.RayCastVertical(
		utils.Vec(x, 3),
		map[string]func(string) bool{
			testKey: func(value string) bool {
				return value == tag1
			},
		},
		tag1Y-3-1,
	)
	a.NoError(err)
	a.Equal(false, hit)

	hit, err = matrix.RayCastVertical(
		utils.Vec(x, 3),
		map[string]func(string) bool{
			testKey: func(value string) bool {
				return value == tag2
			},
		},
		tag1Y-3-1,
	)
	a.NoError(err)
	a.Equal(true, hit)
}

func TestMatrix_RayCastVertical_reverse(t *testing.T) {
	a := require.New(t)
	x := 4
	matrix := NewMatrix(10, 10)

	matrix.Cell(utils.Vec(x, 4)).Tag(testKey, tag1)

	hit, err := matrix.RayCastVertical(
		utils.Vec(x, 8),
		map[string]func(string) bool{
			testKey: func(value string) bool {
				return value == tag1
			},
		},
		-4,
	)
	a.NoError(err)
	a.Equal(true, hit)

	hit, err = matrix.RayCastVertical(
		utils.Vec(x, 8),
		map[string]func(string) bool{
			testKey: func(value string) bool {
				return value == tag1
			},
		},
		-3,
	)
	a.NoError(err)
	a.Equal(false, hit)
}

func TestMatrix_RayCastVertical_Fail(t *testing.T) {
	a := require.New(t)
	matrix := NewMatrix(10, 10)

	_, err := matrix.RayCastVertical(utils.Vec(0, 10), map[string]func(string) bool{}, 2)
	a.Error(err)
}
