package bitvec

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBitVec_MaxState(t *testing.T) {
	tests := []struct {
		size, max uint64
	}{
		{1, 1},
		{2, 3},
		{4, 15},
		{64, math.MaxUint64},
	}

	for tno, test := range tests {
		vec, _ := NewBitVec(1, test.size)

		max := vec.MaxState()
		assert.Equal(t, test.max, max, "Test No. %v", tno)
	}

}

func TestNewBitVec(t *testing.T) {
	tests := []struct {
		count, size uint64
		output      []uint64
	}{
		{20, 70, nil},
		{20, 8, []uint64{0, 0, 0}},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		if test.size > MAXVECSIZE {
			assert.EqualError(t, err, "state Size greater 64 not allowed")
			assert.Nil(t, vec)
		} else {
			assert.Nil(t, err, "Unexpected Error")
			assert.Equal(t, test.output, vec.Data)
		}
	}

}

func TestBitVec_String(t *testing.T) {
	tests := []struct {
		count, size uint64
		data        []uint64
		output      string
	}{
		{32, 2, []uint64{3027}, "[32|2] [0000000000000000000000000000000000000000000000000000101111010011]"},
		{16, 4, []uint64{11297799}, "[16|4] [0000000000000000000000000000000000000000101011000110010000000111]"},
		{8, 10, []uint64{369099440, 4785074604081152}, "[8|10] [0000000000000000000000000000000000010110000000000000001010110000 0000000000010001000000000000000000000000000000000000000000000000]"},
		{42, 3, []uint64{81, 9223372036854777604}, "[42|3] [0000000000000000000000000000000000000000000000000000000001010001 1000000000000000000000000000000000000000000000000000011100000100]"},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

		vec.Data = test.data
		assert.Equal(t, test.output, fmt.Sprint(vec))
	}
}

func TestBitVec_Set(t *testing.T) {
	tests := []struct {
		count, size          uint64
		index, state, output []uint64
	}{
		{32, 2, []uint64{30, 31, 27, 28, 26, 29}, []uint64{0, 3, 3, 3, 2, 1}, []uint64{3027}},
		{16, 4, []uint64{10, 11, 15, 12, 13, 14}, []uint64{10, 12, 7, 6, 4, 0}, []uint64{11297799}},
		{0, 0, []uint64{}, []uint64{}, []uint64{}},
		{20, 64, []uint64{6, 4, 3, 1, 0, 10, 12, 16}, []uint64{111, 1000, 36510, 12000, 4, 786, 1001, 9999}, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0, 786, 0, 1001, 0, 0, 0, 9999, 0, 0, 0}},
		{16, 32, []uint64{1, 9, 5, 13, 7}, []uint64{33, 5676, 12101, 333, 1001}, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}},
		{8, 10, []uint64{5, 7, 3}, []uint64{43, 17, 22}, []uint64{369099440, 4785074604081152}},
		{42, 3, []uint64{19, 21, 41, 39}, []uint64{5, 6, 1, 7}, []uint64{81, 9223372036854777604}},
		{30, 3, []uint64{31}, []uint64{5}, []uint64{0, 0}},
		{42, 3, []uint64{19, 21, 41, 39, 0}, []uint64{5, 6, 1, 7, 8}, []uint64{81, 9223372036854777604}},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

		for i := 0; i < len(test.index); i++ {

			err := vec.Set(test.index[i], test.state[i])
			if err != nil {
				if test.index[i] >= test.count {
					assert.EqualError(t, err, fmt.Sprintf("index too large for BitVec Count (max: %v)", test.count))
				} else if test.state[i] > vec.MaxState() {
					assert.EqualError(t, err, fmt.Sprintf("state too large for BitVec state (maxL %v)", vec.MaxState()))
				} else {
					assert.Nil(t, err, "Unexpected Error")
				}
			}

		}

		assert.Equal(t, test.output, vec.Data)
	}

}

func TestBitVec_Unset(t *testing.T) {
	tests := []struct {
		count, size, index uint64
		data, output       []uint64
	}{
		{32, 2, 28, []uint64{3027}, []uint64{2835}},
		{16, 4, 11, []uint64{11297799}, []uint64{10511367}},
		{20, 64, 16, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0, 786, 0, 1001, 0, 0, 0, 9999, 0, 0, 0}, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0, 786, 0, 1001, 0, 0, 0, 0, 0, 0, 0}},
		{16, 32, 8, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}},
		{16, 32, 13, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, []uint64{33, 0, 12101, 1001, 5676, 0, 0, 0}},
		{8, 10, 7, []uint64{369099440, 4785074604081152}, []uint64{369099440, 0}},
		{16, 4, 7, []uint64{11297799}, []uint64{11297799}},
		{42, 3, 21, []uint64{81, 9223372036854777604}, []uint64{80, 1796}},
		{42, 3, 49, []uint64{81, 9223372036854777604}, []uint64{81, 9223372036854777604}},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

		vec.Data = test.data
		err = vec.Unset(test.index)
		if err != nil {
			if test.index >= test.count {
				assert.EqualError(t, err, fmt.Sprintf("index too large for BitVec Count (max: %v)", test.count))
			} else {
				assert.Nil(t, err, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, vec.Data)
	}
}

func TestBitVec_Has(t *testing.T) {
	tests := []struct {
		count, size, index, state uint64
		data                      []uint64
		output                    bool
	}{
		{32, 2, 27, 3, []uint64{3027}, true},
		{16, 4, 15, 7, []uint64{11297799}, true},
		{16, 4, 12, 8, []uint64{11297799}, false},
		{20, 64, 11, 0, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0, 786, 0, 1001, 0, 0, 0, 9999, 0, 0, 0}, true},
		{16, 32, 13, 333, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, true},
		{16, 32, 0, 137, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, false},
		{8, 10, 4, 199, []uint64{369099440, 4785074604081152}, false},
		{8, 10, 7, 17, []uint64{369099440, 4785074604081152}, true},
		{42, 3, 21, 6, []uint64{81, 9223372036854777604}, true},
		{42, 3, 16, 3, []uint64{81, 9223372036854777604}, false},
		{42, 3, 50, 3, []uint64{81, 9223372036854777604}, false},
		{8, 10, 7, 10000, []uint64{369099440, 4785074604081152}, false},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

		vec.Data = test.data
		check, err2 := vec.Has(test.index, test.state)
		if err2 != nil {
			if test.index >= test.count {
				assert.EqualError(t, err2, fmt.Sprintf("index too large for BitVec Count (max: %v)", test.count))
			} else if test.state > vec.MaxState() {
				assert.EqualError(t, err2, fmt.Sprintf("state too large for BitVec state (maxL %v)", vec.MaxState()))
			} else {
				assert.Nil(t, err2, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, check)
	}
}

func TestBitVec_State(t *testing.T) {
	tests := []struct {
		count, size, index uint64
		data               []uint64
		output             uint64
	}{
		{32, 2, 27, []uint64{3027}, 3},
		{16, 4, 15, []uint64{11297799}, 7},
		{16, 4, 12, []uint64{11297799}, 6},
		{20, 64, 11, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0, 786, 0, 1001, 0, 0, 0, 9999, 0, 0, 0}, 0},
		{16, 32, 13, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, 333},
		{16, 32, 0, []uint64{33, 0, 12101, 1001, 5676, 0, 333, 0}, 0},
		{8, 10, 4, []uint64{369099440, 4785074604081152}, 0},
		{8, 10, 7, []uint64{369099440, 4785074604081152}, 17},
		{42, 3, 21, []uint64{81, 9223372036854777604}, 6},
		{42, 3, 10, []uint64{81, 9223372036854777604}, 0},
		{42, 3, 50, []uint64{81, 9223372036854777604}, 0},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

		vec.Data = test.data
		value, err2 := vec.State(test.index)
		if err2 != nil {
			if test.index >= test.count {
				assert.EqualError(t, err2, fmt.Sprintf("index too large for BitVec Count (max: %v)", test.count))
			} else {
				assert.Nil(t, err2, "Unexpected Error")
			}
		}
		assert.Equal(t, test.output, value)
	}
}

func TestBitVec_GetIndexes(t *testing.T) {
	tests := []struct {
		count, size, state uint64
		data, output       []uint64
	}{
		{32, 2, 3, []uint64{3027}, []uint64{27, 28, 31}},
		{16, 4, 1, []uint64{11297799}, []uint64{}},
		{10, 64, 0, []uint64{4, 12000, 0, 36510, 1000, 0, 111, 0, 0, 0}, []uint64{2, 5, 7, 8, 9}},
		{8, 10, 17, []uint64{369099440, 4785074604081152}, []uint64{7}},
		{42, 3, 21, []uint64{81, 9223372036854777604}, nil},
	}

	for _, test := range tests {
		vec, err := NewBitVec(test.count, test.size)
		assert.Nil(t, err, "Unexpected Error")

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
