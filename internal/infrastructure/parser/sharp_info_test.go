package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSharpInfo_Minimal(t *testing.T) {
	t.Parallel()

	in := `SW_GUID=AB
endianness=1
enable_endianness_per_job=0
reproducibility_disable=2
`
	out, err := parseSharpInfo([]byte(in))
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "0xab", out[0].NodeGUID)
	require.Equal(t, 1, out[0].Endianness)
	require.Equal(t, 0, out[0].EnableEndiannessPerJob)
	require.Equal(t, 2, out[0].ReproducibilityDisable)
}

func TestParseSharpInfo_MultipleBlocks(t *testing.T) {
	t.Parallel()

	in := `SW_GUID=1
endianness=1
---
SW_GUID=2
endianness=3
`
	out, err := parseSharpInfo([]byte(in))
	require.NoError(t, err)
	require.Len(t, out, 2)
	require.Equal(t, "0x1", out[0].NodeGUID)
	require.Equal(t, "0x2", out[1].NodeGUID)
	require.Equal(t, 3, out[1].Endianness)
}

func TestParseSharpInfo_InvalidKeyValue(t *testing.T) {
	t.Parallel()

	_, err := parseSharpInfo([]byte(`SW_GUID=1
notkeyvalue
`))
	require.Error(t, err)
}

func TestParseSharpInfo_InvalidInt(t *testing.T) {
	t.Parallel()

	_, err := parseSharpInfo([]byte(`SW_GUID=1
endianness=nan
`))
	require.Error(t, err)
}

func TestParseSharpInfo_IgnoresLeadingNoiseUntilFirstSW_GUID(t *testing.T) {
	t.Parallel()

	in := `---
garbage line
SW_GUID=5
endianness=1
`
	out, err := parseSharpInfo([]byte(in))
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "0x5", out[0].NodeGUID)
	require.Equal(t, 1, out[0].Endianness)
}
