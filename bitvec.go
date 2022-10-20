package bitvec

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// MAXVECSIZE is the maximum allowed size for a BitVec state.
// The size represents the number of bits consumed for a response state
const MAXVECSIZE = 64

// BitVec is a struct that maintains some number of responses
type BitVec struct {
	// count is the number of responses
	count uint64
	// size is the number of bits required for a response
	size uint64
	// data stores the responses according to their indices
	data []uint64
}

// String implements the Stringer interface for BitVec
func (vec *BitVec) String() string {
	return fmt.Sprintf("[%v|%v] %064b", vec.count, vec.size, vec.data)
}

// NewBitVec is a constructor function for BitVec.
// Returns an error if size is greater than MAXVECSIZE
func NewBitVec(count, size uint64) (*BitVec, error) {
	// Check if given size is under MAXVECSIZE
	if size > MAXVECSIZE {
		return nil, errors.New("state size greater 64 not allowed")
	}

	return &BitVec{
		count: count, size: size,
		data: make([]uint64, int(math.Ceil(float64(count*size)/64))),
	}, nil
}

// MaxState is a method of BitVec that returns the maximum value for a state for that BitVec.
// It is calculated as 2^StateBits-1.
func (vec *BitVec) MaxState() uint64 {
	return 1<<vec.size - 1
}

// Set is a method of BitVec that sets a given state at given index.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the BitVec.
func (vec *BitVec) Set(index, state uint64) error {
	// Check for out of bounds index
	if index >= vec.count {
		return errors.Errorf("index too large for BitVec count (max: %v)", vec.count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return errors.Errorf("state too large for BitVec state (maxL %v)", vec.MaxState())
	}

	// Get the start and end positions for the response state in the data
	start := (index * vec.size) / 64
	end := (((index + 1) * vec.size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.size) % 64
	endBit := (((index + 1) * vec.size) - 1) % 64

	// If the response is contained within a single uint64 in data
	if start == end {
		temp := state << (64 - vec.size - startBit)
		vec.data[start] |= temp

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		startTemp := state >> (vec.size - 64 + startBit)
		endTemp := state << (64 - endBit - 1)

		vec.data[start] |= startTemp
		vec.data[end] |= endTemp
	}

	return nil
}

// Unset is a method of BitVec that unsets the state for a given index.
// Returns an error index is out of bounds.
func (vec *BitVec) Unset(index uint64) error {
	// Check for out of bounds index
	if index >= vec.count {
		return errors.Errorf("index too large for BitVec count (max: %v)", vec.count)
	}

	// Get the start and end positions for the response state in the data
	start := (index * vec.size) / 64
	end := (((index + 1) * vec.size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.size) % 64
	endBit := (((index + 1) * vec.size) - 1) % 64

	// If the response is contained within a single uint64 in data
	if start == end {
		temp := vec.MaxState() << (64 - vec.size - startBit)
		temp ^= 1<<64 - 1

		vec.data[start] &= temp

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		max := uint64(1<<64 - 1)
		startTemp := max << (64 - startBit)
		endTemp := max >> (endBit + 1)

		vec.data[start] &= startTemp
		vec.data[end] &= endTemp
	}

	return nil
}

// Has is a method of BitVec that checks whether the state at a given index matches the given state.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the BitVec.
func (vec *BitVec) Has(index, state uint64) (bool, error) {
	// Check for out of bounds index
	if index >= vec.count {
		return false, errors.Errorf("index too large for BitVec count (max: %v)", vec.count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return false, errors.Errorf("state too large for BitVec state (maxL %v)", vec.MaxState())
	}

	// Get the start and end positions for the response state in the data
	start := (index * vec.size) / 64
	end := (((index + 1) * vec.size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.size) % 64
	endBit := (((index + 1) * vec.size) - 1) % 64

	var value uint64
	max := vec.MaxState()

	// If the response is contained within a single uint64 in data
	if start == end {
		temp := max << (64 - vec.size - startBit)
		value = temp & vec.data[start]
		value >>= 64 - vec.size - startBit

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		temp := max >> (vec.size - 64 + startBit)
		value = temp & vec.data[start]
		value <<= vec.size - 64 + startBit

		temp = max << (64 - endBit - 1)
		value2 := temp & vec.data[end]
		value2 >>= 64 - endBit - 1

		value |= value2
	}

	return value == state, nil
}

// State is a method of BitVec that returns the state at a given index.
// Returns an error if the index is out of bounds.
func (vec *BitVec) State(index uint64) (uint64, error) {
	// Check for out of bounds index
	if index >= vec.count {
		return 0, errors.Errorf("index too large for BitVec count (max: %v)", vec.count)
	}

	// Get the start and end positions for the response state in the data
	start := (index * vec.size) / 64
	end := (((index + 1) * vec.size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.size) % 64
	endBit := (((index + 1) * vec.size) - 1) % 64

	var value uint64
	max := vec.MaxState()

	// If the response is contained within a single uint64 in data
	if start == end {
		temp := max << (64 - vec.size - startBit)
		value = temp & vec.data[start]
		value >>= 64 - vec.size - startBit

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		temp := max >> (vec.size - 64 + startBit)
		value = temp & vec.data[start]
		value <<= vec.size - 64 + startBit

		temp = max << (64 - endBit - 1)
		value2 := temp & vec.data[end]
		value2 >>= 64 - endBit - 1

		value |= value2
	}

	return value, nil
}
