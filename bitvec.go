package bitvec

import (
	"fmt"
	"math"
	"sync"

	"github.com/pkg/errors"
)

// MAXVECSIZE is the maximum allowed Size for a BitVec state.
// The Size represents the number of bits consumed for a response state
const MAXVECSIZE = 64

// BitVec is a struct that maintains some number of responses
type BitVec struct {
	// mu is the thread safety mutex
	mu sync.Mutex

	// Count is the number of responses
	Count uint64
	// Size is the number of bits required for a response
	Size uint64
	// Data stores the responses according to their indices
	Data []uint64
}

// NewBitVec is a constructor function for BitVec.
// Returns an error if Size is greater than MAXVECSIZE
func NewBitVec(count, size uint64) (*BitVec, error) {
	// Check if given Size is under MAXVECSIZE
	if size > MAXVECSIZE {
		return nil, errors.New("state size greater 64 not allowed")
	}

	return &BitVec{
		mu: sync.Mutex{}, Count: count, Size: size,
		Data: make([]uint64, int(math.Ceil(float64(count*size)/64))),
	}, nil
}

// String implements the Stringer interface for BitVec
func (vec *BitVec) String() string {
	return fmt.Sprintf("[%v|%v] %064b", vec.Count, vec.Size, vec.Data)
}

// MaxState is a method of BitVec that returns the maximum value for a state for that BitVec.
// It is calculated as 2^StateBits-1.
func (vec *BitVec) MaxState() uint64 {
	return 1<<vec.Size - 1
}

// Set is a method of BitVec that sets a given state at given index.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the BitVec.
func (vec *BitVec) Set(index, state uint64) error {
	// Check for out of bounds index
	if index >= vec.Count {
		return errors.Errorf("index too large for bitvec count (max: %v)", vec.Count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return errors.Errorf("state too large for bitvec state (max: %v)", vec.MaxState())
	}

	// Acquire the mutex
	vec.mu.Lock()
	defer vec.mu.Unlock()

	// Get the start and end positions for the response state in the Data
	start := (index * vec.Size) / 64
	end := (((index + 1) * vec.Size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.Size) % 64
	endBit := (((index + 1) * vec.Size) - 1) % 64

	// If the response is contained within a single uint64 in Data
	if start == end {
		temp := state << (64 - vec.Size - startBit)
		vec.Data[start] |= temp

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		startTemp := state >> (vec.Size - 64 + startBit)
		endTemp := state << (64 - endBit - 1)

		vec.Data[start] |= startTemp
		vec.Data[end] |= endTemp
	}

	return nil
}

// Unset is a method of BitVec that unsets the state for a given index.
// Returns an error index is out of bounds.
func (vec *BitVec) Unset(index uint64) error {
	// Check for out of bounds index
	if index >= vec.Count {
		return errors.Errorf("index too large for bitvec count (max: %v)", vec.Count)
	}

	// Acquire the mutex
	vec.mu.Lock()
	defer vec.mu.Unlock()

	// Get the start and end positions for the response state in the Data
	start := (index * vec.Size) / 64
	end := (((index + 1) * vec.Size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.Size) % 64
	endBit := (((index + 1) * vec.Size) - 1) % 64

	// If the response is contained within a single uint64 in Data
	if start == end {
		temp := vec.MaxState() << (64 - vec.Size - startBit)
		temp ^= 1<<64 - 1

		vec.Data[start] &= temp

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		max := uint64(1<<64 - 1)
		startTemp := max << (64 - startBit)
		endTemp := max >> (endBit + 1)

		vec.Data[start] &= startTemp
		vec.Data[end] &= endTemp
	}

	return nil
}

// Has is a method of BitVec that checks whether the state at a given index matches the given state.
// Returns an error if the index is out of bounds or if the state value exceeds the maximum for the BitVec.
func (vec *BitVec) Has(index, state uint64) (bool, error) {
	// Check for out of bounds index
	if index >= vec.Count {
		return false, errors.Errorf("index too large for bitvec count (max: %v)", vec.Count)
	}

	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return false, errors.Errorf("state too large for bitvec state (max: %v)", vec.MaxState())
	}

	// Get the start and end positions for the response state in the Data
	start := (index * vec.Size) / 64
	end := (((index + 1) * vec.Size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.Size) % 64
	endBit := (((index + 1) * vec.Size) - 1) % 64

	var value uint64
	max := vec.MaxState()

	// If the response is contained within a single uint64 in Data
	if start == end {
		temp := max << (64 - vec.Size - startBit)
		value = temp & vec.Data[start]
		value >>= 64 - vec.Size - startBit

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		temp := max >> (vec.Size - 64 + startBit)
		value = temp & vec.Data[start]
		value <<= vec.Size - 64 + startBit

		temp = max << (64 - endBit - 1)
		value2 := temp & vec.Data[end]
		value2 >>= 64 - endBit - 1

		value |= value2
	}

	return value == state, nil
}

// State is a method of BitVec that returns the state at a given index.
// Returns an error if the index is out of bounds.
func (vec *BitVec) State(index uint64) (uint64, error) {
	// Check for out of bounds index
	if index >= vec.Count {
		return 0, errors.Errorf("index too large for bitvec count (max: %v)", vec.Count)
	}

	// Get the start and end positions for the response state in the Data
	start := (index * vec.Size) / 64
	end := (((index + 1) * vec.Size) - 1) / 64

	// Calculate the start and end bit positions
	startBit := (index * vec.Size) % 64
	endBit := (((index + 1) * vec.Size) - 1) % 64

	var value uint64
	max := vec.MaxState()

	// If the response is contained within a single uint64 in Data
	if start == end {
		temp := max << (64 - vec.Size - startBit)
		value = temp & vec.Data[start]
		value >>= 64 - vec.Size - startBit

	} else {
		// Calculate new values for both affected positions.
		// NOTE: This logic fails for a response that spans beyond 2 uint64.
		// This is regulated by MAXVECSIZE.

		temp := max >> (vec.Size - 64 + startBit)
		value = temp & vec.Data[start]
		value <<= vec.Size - 64 + startBit

		temp = max << (64 - endBit - 1)
		value2 := temp & vec.Data[end]
		value2 >>= 64 - endBit - 1

		value |= value2
	}

	return value, nil
}

// Indexes is a method of BitVec that returns the slice of indexes matching the given state.
// Returns an error if state value exceeds the maximum for the BitVec.
func (vec *BitVec) Indexes(state uint64) ([]uint64, error) {
	// Check for state value too large for BitVec
	if state > vec.MaxState() {
		return nil, errors.Errorf("state too large for bitvec state (max: %v)", vec.MaxState())
	}

	// Iterate over the BitVec and check each index for
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
