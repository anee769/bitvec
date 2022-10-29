package bitvec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiBit_String(t *testing.T) {
	tests := []struct {
		dibit  *DiBit
		output string
	}{
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			"[64] [0000000000000000000000000000000000000000000000000000000000110010 0000000000000000000000000000000000000000000000000000000011000011]",
		},
		{
			&DiBit{Count: 33, Data: []uint64{1059, 4611686018427387904}},
			"[33] [0000000000000000000000000000000000000000000000000000010000100011 0100000000000000000000000000000000000000000000000000000000000000]",
		},
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			"[32] [0000000000000000000000000000000000000000000000000000101111010011]",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, test.dibit.String())
	}
}

func TestDiBit_Set(t *testing.T) {
	tests := []struct {
		dibit    *DiBit
		idx, val uint64
		output   []uint64
		err      string
	}{
		{
			&DiBit{Count: 32, Data: []uint64{0}},
			30, 3, []uint64{12}, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{12}},
			10, 2, []uint64{8796093022220}, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{0}},
			14, 0, []uint64{0}, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{0, 0}},
			63, 3, []uint64{0, 3}, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			37, 0, []uint64{3027}, "index too large for dibit count (max: 32)",
		},
		{
			&DiBit{Count: 64, Data: []uint64{3, 3}},
			0, 7, []uint64{3, 3}, "state too large for dibit state (max: 3)",
		},
	}

	for _, test := range tests {
		err := test.dibit.Set(test.idx, test.val)
		assert.Equal(t, test.dibit.Data, test.output)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}

}

func TestDiBit_Unset(t *testing.T) {
	tests := []struct {
		dibit  *DiBit
		idx    uint64
		output []uint64
		err    string
	}{
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			28, []uint64{2835}, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			63, []uint64{50, 192}, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			11, []uint64{50, 195}, "",
		},
		{
			&DiBit{Count: 33, Data: []uint64{1059, 4611686018427387904}},
			1, []uint64{1059, 4611686018427387904}, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			100, []uint64{50, 195}, "index too large for dibit count (max: 64)",
		},
	}

	for _, test := range tests {
		err := test.dibit.Unset(test.idx)
		assert.Equal(t, test.dibit.Data, test.output)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestDiBit_Has(t *testing.T) {
	tests := []struct {
		dibit    *DiBit
		idx, val uint64
		exists   bool
		err      string
	}{
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			27, 3, true, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			63, 3, true, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			11, 2, false, "",
		},
		{
			&DiBit{Count: 33, Data: []uint64{1059, 4611686018427387904}},
			33, 1, false, "index too large for dibit count (max: 33)",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			11, 5, false, "state too large for dibit state (max: 3)",
		},
	}

	for _, test := range tests {
		exists, err := test.dibit.Has(test.idx, test.val)
		assert.Equal(t, test.exists, exists)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestDiBit_State(t *testing.T) {
	tests := []struct {
		dibit       *DiBit
		idx, output uint64
		err         string
	}{
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			27, 3, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			63, 3, "",
		},
		{
			&DiBit{Count: 33, Data: []uint64{1059, 4611686018427387904}},
			1, 0, "",
		},
		{
			&DiBit{Count: 33, Data: []uint64{1059, 4611686018427387904}},
			32, 1, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			1000, 0, "index too large for dibit count (max: 32)",
		},
	}

	for _, test := range tests {
		state, err := test.dibit.State(test.idx)
		assert.Equal(t, test.output, state)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestDiBit_Indexes(t *testing.T) {
	tests := []struct {
		dibit   *DiBit
		query   uint64
		indexes []uint64
		err     string
	}{
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			3, []uint64{27, 28, 31}, "",
		},
		{
			&DiBit{Count: 64, Data: []uint64{50, 195}},
			3, []uint64{29, 60, 63}, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{9223372036854775807}},
			0, []uint64{}, "",
		},
		{
			&DiBit{Count: 32, Data: []uint64{3027}},
			18, []uint64{}, "state too large for dibit state (max: 3)",
		},
	}

	for _, test := range tests {
		indexes, err := test.dibit.Indexes(test.query)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
			assert.Equal(t, test.indexes, indexes)
		} else {
			assert.EqualError(t, err, test.err)
			assert.Nil(t, indexes)
		}
	}
}
