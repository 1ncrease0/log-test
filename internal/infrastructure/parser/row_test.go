package parser

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCol(t *testing.T) {
	t.Parallel()

	idx := map[string]int{"A": 0, "B": 1}
	row := []string{" x ", "y"}

	v, err := getCol(row, idx, "A")
	require.NoError(t, err)
	require.Equal(t, "x", v)

	_, err = getCol(row, idx, "Missing")
	require.Error(t, err)

	short := []string{"a"}
	_, err = getCol(short, idx, "B")
	require.Error(t, err)
}

func TestInt64Col_SignedAndUnsigned(t *testing.T) {
	t.Parallel()

	idx := map[string]int{"N": 0}
	v, err := int64Col([]string{"42"}, idx, "N")
	require.NoError(t, err)
	require.Equal(t, int64(42), v)

	v, err = int64Col([]string{"9223372036854775808"}, idx, "N")
	require.NoError(t, err)
	require.Equal(t, int64(math.MinInt64), v)
}

func TestOptIntCol_NA(t *testing.T) {
	t.Parallel()

	idx := map[string]int{"N": 0}
	v, err := optIntCol([]string{csvOptionalNA}, idx, "N")
	require.ErrorIs(t, err, errOptionalNA)
	require.Nil(t, v)

	v, err = optIntCol([]string{"7"}, idx, "N")
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 7, *v)
}

func TestNormCSVNodeGUID(t *testing.T) {
	t.Parallel()
	require.Equal(t, "0xab", normCSVNodeGUID(" 0XAB "))
}

func TestNormSharpNodeGUID(t *testing.T) {
	t.Parallel()
	require.Equal(t, "0xab", normSharpNodeGUID("ab"))
	require.Equal(t, "0xab", normSharpNodeGUID("0xAB"))
}
