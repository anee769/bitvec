package bitvec

import (
	"fmt"
	"math"
	"sync"

	"github.com/pkg/errors"
)

// DIBITSIZE is the size for a DiBit state.
const DIBITSIZE = 2

// DiBit is a struct that maintains some number of responses
type DiBit struct {
	// mu is the thread safety mutex
	mu sync.Mutex

	// Count is the number of responses
	Count uint64
	// Data stores the responses according to their indices
	Data []uint64
}

// NewDiBit is a constructor function for DiBit.
func NewDiBit(count uint64) *DiBit {
	return &DiBit{
		mu: sync.Mutex{}, Count: count,
		Data: make([]uint64, int(math.Ceil(float64(count*DIBITSIZE)/64))),
	}
}

// String implements the Stringer interface for DiBit
func (vec *DiBit) String() string {
	return fmt.Sprintf("[%v] %064b", vec.Count, vec.Data)
}

// MaxState is a method of DiBit that returns the maximum value for the state.
// It is calculated as 2^StateBits-1.
func (vec *DiBit) MaxState() uint64 {
	return 1<<DIBITSIZE - 1
}

// Set is a method of DiBit that sets a given state at given index.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the DiBit.
func (vec *DiBit) Set(index, state uint64) error {
	// Check for out of bounds index
	if index >= vec.Count {
		return errors.Errorf("index too large for dibit count (max: %v)", vec.Count)
	}

	// Check for state value too large for DiBit
	if state > vec.MaxState() {
		return errors.Errorf("state too large for dibit state (max: %v)", vec.MaxState())
	}

	// Acquire the mutex
	vec.mu.Lock()
	defer vec.mu.Unlock()

	// Get the start position for the response state in the Data
	start := (index * DIBITSIZE) / 64

	// Calculate the start bit position
	startBit := (index * DIBITSIZE) % 64

	temp := state << (64 - DIBITSIZE - startBit)
	vec.Data[start] |= temp

	return nil
}

// Unset is a method of DiBit that unsets the state for a given index.
// Returns an error index is out of bounds.
func (vec *DiBit) Unset(index uint64) error {
	// Check for out of bounds index
	if index >= vec.Count {
		return errors.Errorf("index too large for dibit count (max: %v)", vec.Count)
	}

	// Acquire the mutex
	vec.mu.Lock()
	defer vec.mu.Unlock()

	// Get the start position for the response state in the Data
	start := (index * DIBITSIZE) / 64

	// Calculate the start bit position
	startBit := (index * DIBITSIZE) % 64

	max := uint64(1<<64 - 1)
	temp := vec.MaxState() << (64 - DIBITSIZE - startBit)
	temp ^= max
	vec.Data[start] &= temp

	return nil
}

// Has is a method of DiBit that checks whether the state at a given index matches the given state.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the DiBit.
func (vec *DiBit) Has(index, state uint64) (bool, error) {
	// Check for out of bounds index
	if index >= vec.Count {
		return false, errors.Errorf("index too large for dibit count (max: %v)", vec.Count)
	}

	// Check for state value too large for DiBit
	if state > vec.MaxState() {
		return false, errors.Errorf("state too large for dibit state (max: %v)", vec.MaxState())
	}

	// Get the start position for the response state in the Data
	start := (index * DIBITSIZE) / 64

	// Calculate the start bit position
	startBit := (index * DIBITSIZE) % 64

	var value uint64

	temp := vec.MaxState() << (64 - DIBITSIZE - startBit)
	value = temp & vec.Data[start]
	value >>= 64 - DIBITSIZE - startBit

	return value == state, nil
}

// State is a method of DiBit that returns the state at a given index.
// Returns an error if the index is out of bounds.
func (vec *DiBit) State(index uint64) (uint64, error) {
	// Check for out of bounds index
	if index >= vec.Count {
		return 0, errors.Errorf("index too large for dibit count (max: %v)", vec.Count)
	}

	// Get the start position for the response state in the Data
	start := (index * DIBITSIZE) / 64

	// Calculate the start bit position
	startBit := (index * DIBITSIZE) % 64

	var value uint64

	temp := vec.MaxState() << (64 - DIBITSIZE - startBit)
	value = temp & vec.Data[start]
	value >>= 64 - DIBITSIZE - startBit

	return value, nil
}

// Indexes is a method of DiBit that returns the slice of indexes matching the given state.
// Returns an error if state value exceeds the maximum for the DiBit.
func (vec *DiBit) Indexes(state uint64) ([]uint64, error) {
	// Check for state value too large for DiBit
	if state > vec.MaxState() {
		return nil, errors.Errorf("state too large for dibit state (max: %v)", vec.MaxState())
	}

	// Iterate over the DiBit and check each index for
	// equality with state and append index if equal
	indexes := make([]uint64, 0)
	for i := uint64(0); i < vec.Count; i++ {
		// Error from Has() can be ignored because max state has already been checked
		// and the count will never overflow because it is bounded by the loop condition.
		if exists, _ := vec.Has(i, state); exists {
			indexes = append(indexes, i)
		}
	}

	return indexes, nil
}
