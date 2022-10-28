package bitvec

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiBit_String(t *testing.T) {
	tests := []struct {
		count  uint64
		data   []uint64
		output string
	}{
		{64, []uint64{50, 195}, "[64] [0000000000000000000000000000000000000000000000000000000000110010 0000000000000000000000000000000000000000000000000000000011000011]"},
		{33, []uint64{1059, 4611686018427387904}, "[33] [0000000000000000000000000000000000000000000000000000010000100011 0100000000000000000000000000000000000000000000000000000000000000]"},
		{32, []uint64{3027}, "[32] [0000000000000000000000000000000000000000000000000000101111010011]"},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)
		vec.Data = test.data

		assert.Equal(t, test.output, fmt.Sprint(vec))
	}
}

func TestDiBit_Set(t *testing.T) {
	tests := []struct {
		count                uint64
		index, state, output []uint64
	}{
		{32, []uint64{30, 31, 27, 28, 26, 29}, []uint64{0, 3, 3, 3, 2, 1}, []uint64{3027}},
		{64, []uint64{31, 29, 63, 60}, []uint64{2, 3, 3, 3}, []uint64{50, 195}},
		{33, []uint64{29, 26, 31, 32}, []uint64{2, 1, 3, 1}, []uint64{1059, 4611686018427387904}},
		{0, []uint64{}, []uint64{}, []uint64{}},
		{32, []uint64{30, 31, 27, 28, 26, 29, 37}, []uint64{0, 3, 3, 3, 2, 1, 0}, []uint64{3027}},
		{33, []uint64{29, 26, 31, 32, 0}, []uint64{2, 1, 3, 1, 7}, []uint64{1059, 4611686018427387904}},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)

		for i := 0; i < len(test.index); i++ {

			err := vec.Set(test.index[i], test.state[i])
			if err != nil {
				if test.index[i] >= test.count {
					assert.EqualError(t, err, fmt.Sprintf("index too large for DiBit Count (max: %v)", test.count))
				} else if test.state[i] > vec.MaxState() {
					assert.EqualError(t, err, fmt.Sprintf("state too large for DiBit state (maxL %v)", vec.MaxState()))
				} else {
					assert.Nil(t, err, "Unexpected Error")
				}
			}

		}

		assert.Equal(t, test.output, vec.Data)
	}

}

func TestDiBit_Unset(t *testing.T) {
	tests := []struct {
		count, index uint64
		data, output []uint64
	}{
		{32, 28, []uint64{3027}, []uint64{2835}},
		{64, 63, []uint64{50, 195}, []uint64{50, 192}},
		{64, 11, []uint64{50, 195}, []uint64{50, 195}},
		{33, 1, []uint64{1059, 4611686018427387904}, []uint64{1059, 4611686018427387904}},
		{33, 32, []uint64{1059, 4611686018427387904}, []uint64{1059, 0}},
		{64, 100, []uint64{50, 195}, []uint64{50, 195}},
		{32, 32, []uint64{3027}, []uint64{3027}},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)
		vec.Data = test.data
		err := vec.Unset(test.index)
		if err != nil {
			if test.index >= test.count {
				assert.EqualError(t, err, fmt.Sprintf("index too large for DiBit Count (max: %v)", test.count))
			} else {
				assert.Nil(t, err, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, vec.Data)
	}
}

func TestDiBit_Has(t *testing.T) {
	tests := []struct {
		count, index, state uint64
		data                []uint64
		output              bool
	}{
		{32, 27, 3, []uint64{3027}, true},
		{64, 63, 3, []uint64{50, 195}, true},
		{64, 11, 2, []uint64{50, 195}, false},
		{33, 1, 0, []uint64{1059, 4611686018427387904}, true},
		{33, 32, 1, []uint64{1059, 4611686018427387904}, true},
		{33, 33, 1, []uint64{1059, 4611686018427387904}, false},
		{64, 11, 5, []uint64{50, 195}, false},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)
		vec.Data = test.data
		check, err := vec.Has(test.index, test.state)
		if err != nil {
			if test.index >= test.count {
				assert.EqualError(t, err, fmt.Sprintf("index too large for DiBit Count (max: %v)", test.count))
			} else if test.state > vec.MaxState() {
				assert.EqualError(t, err, fmt.Sprintf("state too large for DiBit state (maxL %v)", vec.MaxState()))
			} else {
				assert.Nil(t, err, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, check)
	}
}

func TestDiBit_State(t *testing.T) {
	tests := []struct {
		count, index uint64
		data         []uint64
		output       uint64
	}{
		{32, 27, []uint64{3027}, 3},
		{64, 63, []uint64{50, 195}, 3},
		{64, 11, []uint64{50, 195}, 0},
		{33, 1, []uint64{1059, 4611686018427387904}, 0},
		{33, 32, []uint64{1059, 4611686018427387904}, 1},
		{64, 64, []uint64{50, 195}, 0},
		{33, 40, []uint64{1059, 4611686018427387904}, 0},
		{32, 1000, []uint64{3027}, 0},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)
		vec.Data = test.data
		value, err := vec.State(test.index)
		if err != nil {
			if test.index >= test.count {
				assert.EqualError(t, err, fmt.Sprintf("index too large for DiBit Count (max: %v)", test.count))
			} else {
				assert.Nil(t, err, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, value)
	}
}

func TestDiBit_GetIndexes(t *testing.T) {
	tests := []struct {
		count, state uint64
		data, output []uint64
	}{
		{32, 3, []uint64{3027}, []uint64{27, 28, 31}},
		{64, 3, []uint64{50, 195}, []uint64{29, 60, 63}},
		{33, 1, []uint64{1059, 4611686018427387904}, []uint64{26, 32}},
		{32, 0, []uint64{9223372036854775807}, []uint64{}},
	}

	for _, test := range tests {
		vec := NewDiBit(test.count)
		vec.Data = test.data
		indexes, err2 := vec.GetIndexes(test.state)
		if err2 != nil {
			if test.state > vec.MaxState() {
				assert.EqualError(t, err2, fmt.Sprintf("state too large for BitVec state (maxL %v)", vec.MaxState()))
			} else {
				assert.Nil(t, err2, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, indexes)
	}
}
