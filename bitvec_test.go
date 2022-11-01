package bitvec

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBitVec_MaxState(t *testing.T) {
	tests := []struct {
		size, max uint64
	}{
		{1, 1},
		{2, 3},
		{4, 15},
		{8, 255},
		{64, math.MaxUint64},
	}

	for _, test := range tests {
		vec, err := NewBitVec(1, test.size)
		require.Nil(t, err, "Unexpected Error")

		max := vec.MaxState()
		assert.Equal(t, test.max, max)
	}
}

func TestNewBitVec(t *testing.T) {
	tests := []struct {
		count, size uint64
		data        []uint64
		err         string
	}{
		{20, 8, []uint64{0, 0, 0}, ""},
		{10, 10, []uint64{0, 0}, ""},
		{12, 24, []uint64{0, 0, 0, 0, 0}, ""},
		{32, 2, []uint64{0}, ""},
		{20, 70, nil, "state size greater 64 not allowed"},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
			assert.Equal(t, test.data, vec.Data)
		} else {
			assert.EqualError(t, err, test.err)
			assert.Nil(t, vec)
		}
	}
}

func TestBitVec_String(t *testing.T) {
	tests := []struct {
		bitvec *BitVec
		output string
	}{
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{3027}},
			"[32|2] [0000000000000000000000000000000000000000000000000000101111010011]",
		},
		{
			&BitVec{Count: 16, Size: 4, Data: []uint64{11297799}},
			"[16|4] [0000000000000000000000000000000000000000101011000110010000000111]",
		},
		{
			&BitVec{Count: 8, Size: 10, Data: []uint64{369099440, 4785074604081152}},
			"[8|10] [0000000000000000000000000000000000010110000000000000001010110000 0000000000010001000000000000000000000000000000000000000000000000]",
		},
		{
			&BitVec{Count: 42, Size: 3, Data: []uint64{81, 9223372036854777604}},
			"[42|3] [0000000000000000000000000000000000000000000000000000000001010001 1000000000000000000000000000000000000000000000000000011100000100]",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, test.bitvec.String())
	}
}

func TestBitVec_Set(t *testing.T) {
	tests := []struct {
		bitvec   *BitVec
		idx, val uint64
		output   []uint64
		err      string
	}{
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{0}},
			30, 3, []uint64{12}, "",
		},
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{12}},
			10, 2, []uint64{8796093022220}, "",
		},
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{0}},
			14, 0, []uint64{0}, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{0, 0}},
			3, 3, []uint64{12884901888, 0}, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 0}},
			10, 2, []uint64{12884901888, 2199023255552}, "",
		},
		{
			&BitVec{Count: 42, Size: 3, Data: []uint64{0, 0}},
			21, 5, []uint64{1, 4611686018427387904}, "",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			30, 12, []uint64{0}, "index too large for bitvec count (max: 10)",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			9, 18, []uint64{0}, "state too large for bitvec state (max: 15)",
		},
	}

	for _, test := range tests {
		err := test.bitvec.Set(test.idx, test.val)
		assert.Equal(t, test.bitvec.Data, test.output)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestBitVec_Unset(t *testing.T) {
	tests := []struct {
		bitvec *BitVec
		idx    uint64
		output []uint64
		err    string
	}{
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{12}},
			30, []uint64{0}, "",
		},
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{8796093022220}},
			10, []uint64{12}, "",
		},
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{0}},
			14, []uint64{0}, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 0}},
			3, []uint64{0, 0}, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 2199023255552}},
			10, []uint64{12884901888, 0}, "",
		},
		{
			&BitVec{Count: 42, Size: 3, Data: []uint64{1, 4611686018427387904}},
			21, []uint64{0, 0}, "",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			30, []uint64{0}, "index too large for bitvec count (max: 10)",
		},
	}

	for _, test := range tests {
		err := test.bitvec.Unset(test.idx)
		assert.Equal(t, test.bitvec.Data, test.output)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestBitVec_Has(t *testing.T) {
	tests := []struct {
		bitvec   *BitVec
		idx, val uint64
		exists   bool
		err      string
	}{
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{12}},
			30, 3, true, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 0}},
			3, 0, false, "",
		},
		{
			&BitVec{Count: 42, Size: 3, Data: []uint64{1, 4611686018427387904}},
			21, 0, false, "",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			30, 12, false, "index too large for bitvec count (max: 10)",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			9, 18, false, "state too large for bitvec state (max: 15)",
		},
	}

	for _, test := range tests {
		exists, err := test.bitvec.Has(test.idx, test.val)
		assert.Equal(t, test.exists, exists)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestBitVec_State(t *testing.T) {
	tests := []struct {
		bitvec      *BitVec
		idx, output uint64
		err         string
	}{
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{12}},
			30, 3, "",
		},
		{
			&BitVec{Count: 32, Size: 2, Data: []uint64{12}},
			10, 0, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 0}},
			3, 3, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 0}},
			5, 0, "",
		},
		{
			&BitVec{Count: 42, Size: 3, Data: []uint64{1, 4611686018427387904}},
			21, 5, "",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			30, 0, "index too large for bitvec count (max: 10)",
		},
	}

	for _, test := range tests {
		state, err := test.bitvec.State(test.idx)
		assert.Equal(t, test.output, state)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
		} else {
			assert.EqualError(t, err, test.err)
		}
	}
}

func TestBitVec_Indexes(t *testing.T) {
	tests := []struct {
		bitvec  *BitVec
		query   uint64
		indexes []uint64
		err     string
	}{
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 12884901888}},
			5, []uint64{}, "",
		},
		{
			&BitVec{Count: 16, Size: 8, Data: []uint64{12884901888, 12884901888}},
			3, []uint64{3, 11}, "",
		},
		{
			&BitVec{Count: 10, Size: 4, Data: []uint64{0}},
			18, []uint64{}, "state too large for bitvec state (max: 15)",
		},
	}

	for _, test := range tests {
		indexes, err := test.bitvec.Indexes(test.query)

		if test.err == "" {
			assert.Nil(t, err, "Unexpected Error")
			assert.Equal(t, test.indexes, indexes)
		} else {
			assert.EqualError(t, err, test.err)
			assert.Nil(t, indexes)
		}
	}
}
