package bitvec

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// STATESIZE is the size for a BitVec state.
const (
	STATESIZE = 2
)

// DiBit is a struct that maintains some number of responses
type DiBit struct {
	// count is the number of responses
	count uint64
	// data stores the responses according to their indices
	data []uint64
}

// String implements the Stringer interface for BitVec
func (vec *DiBit) String() string {
	return fmt.Sprintf("[%v] %064b", vec.count, vec.data)
}

// NewDiBit is a constructor function for DiBit.
func NewDiBit(count uint64) *DiBit {
	return &DiBit{
		count: count,
		data:  make([]uint64, int(math.Ceil(float64(count*STATESIZE)/64))),
	}
}

// MaxState is a method of DiBit that returns the maximum value for the state.
// It is calculated as 2^StateBits-1.
func (vec *DiBit) MaxState() uint64 {
	return 1<<STATESIZE - 1
}

// Set is a method of DiBit that sets a given state at given index.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the DiBit.
func (vec *DiBit) Set(index, state uint64) error {
	// Check for out of bounds index
	if index >= vec.count {
		return errors.Errorf("index too large for DiBit count (max: %v)", vec.count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return errors.Errorf("state too large for DiBit state (maxL %v)", vec.MaxState())
	}

	// Get the start position for the response state in the data
	start := (index * STATESIZE) / 64

	// Calculate the start bit position
	startBit := (index * STATESIZE) % 64

	temp := state << (64 - STATESIZE - startBit)
	vec.data[start] |= temp

	return nil
}

// Unset is a method of DiBit that unsets the state for a given index.
// Returns an error index is out of bounds.
func (vec *DiBit) Unset(index uint64) error {
	// Check for out of bounds index
	if index >= vec.count {
		return errors.Errorf("index too large for DiBit count (max: %v)", vec.count)
	}

	// Get the start position for the response state in the data
	start := (index * STATESIZE) / 64

	// Calculate the start bit position
	startBit := (index * STATESIZE) % 64

	max := uint64(1<<64 - 1)
	temp := vec.MaxState() << (64 - STATESIZE - startBit)
	temp ^= max
	vec.data[start] &= temp

	return nil
}

// Has is a method of DiBit that checks whether the state at a given index matches the given state.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the DiBit.
func (vec *DiBit) Has(index, state uint64) (bool, error) {
	// Check for out of bounds index
	if index >= vec.count {
		return false, errors.Errorf("index too large for DiBit count (max: %v)", vec.count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return false, errors.Errorf("state too large for DiBit state (maxL %v)", vec.MaxState())
	}

	// Get the start position for the response state in the data
	start := (index * STATESIZE) / 64

	// Calculate the start bit position
	startBit := (index * STATESIZE) % 64

	var value uint64

	temp := vec.MaxState() << (64 - STATESIZE - startBit)
	value = temp & vec.data[start]
	value >>= 64 - STATESIZE - startBit

	return value == state, nil
}

// State is a method of DiBit that returns the state at a given index.
// Returns an error if the index is out of bounds.
func (vec *DiBit) State(index uint64) (uint64, error) {
	// Check for out of bounds index
	if index >= vec.count {
		return 0, errors.Errorf("index too large for DiBit count (max: %v)", vec.count)
	}

	// Get the start position for the response state in the data
	start := (index * STATESIZE) / 64

	// Calculate the start bit position
	startBit := (index * STATESIZE) % 64

	var value uint64

	temp := vec.MaxState() << (64 - STATESIZE - startBit)
	value = temp & vec.data[start]
	value >>= 64 - STATESIZE - startBit

	return value, nil
}
