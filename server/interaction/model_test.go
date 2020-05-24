package interaction

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeansSerializedInfo(t *testing.T) {
	testCases := []struct {
		inputMode Mode
		input     *Means
	}{{
		inputMode: UnsetMode,
		input: &Means{
			Mode: UnsetMode,
		},
	}, {
		inputMode: Unknown,
		input: &Means{
			Mode: Unknown,
		},
	}, {
		inputMode: Localhost,
		input: &Means{
			Mode: Localhost,
			Info: `8484`,
		},
	}, {
		inputMode: NPC,
		input: &Means{
			Mode: NPC,
		},
	}}

	for _, tc := range testCases {
		msg := fmt.Sprintf("running test for %v", tc.inputMode)
		ser, err := tc.input.GetSerializedInfo()
		require.NoError(t, err, msg)

		output := &Means{
			Mode: tc.inputMode,
		}
		err = output.AddSerializedInfo(ser)
		require.NoError(t, err, msg)
		assert.Equal(t, tc.input, output, msg)
	}
}
