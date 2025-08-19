package machine

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Memory is a simple byte array for EVM memory, which is volatile.
type Memory struct {
	data []byte
}

// NewMemory creates a new memory model.
func NewMemory() *Memory {
	return &Memory{
		data: make([]byte, 0, 1024),
	}
}

func (m *Memory) GetData() []byte {
	return m.data
}

// resize expands the memory to a new size. The EVM expands memory in
// 32-byte words.
func (m *Memory) resize(size uint64) {
	if uint64(len(m.data)) < size {
		// Calculate the new size in 32-byte words
		newSizeInWords := (size + 31) / 32
		newSize := newSizeInWords * 32

		// In a real EVM, you would calculate the quadratic gas cost here BEFORE resizing.
		// For now, we'll just resize the slice.

		newData := make([]byte, newSize)
		copy(newData, m.data)
		m.data = newData
	}
}

// Set stores a slice of bytes at a specific offset in memory.
// It will expand the memory if necessary.
func (m *Memory) Set(offset uint64, value []byte) {
	size := uint64(len(value))
	if size == 0 {
		return
	}

	// check if memory needs to be expanded
	if requiredSize := offset + size; uint64(len(m.data)) < requiredSize {
		m.resize(requiredSize)
	}
	copy(m.data[offset:], value)
}

// Set32 stores a 32-byte value (from a big.Int) at a specific offset.
// This is common for opcodes like MSTORE.
func (m *Memory) Set32(offset uint64, value *big.Int) {
	// big.Int.Bytes() returns a big-endian byte slice. It may be less than 32 bytes.
	valBytes := value.Bytes()
	paddedVal := make([]byte, 32)
	// Copy the value bytes to the end of the 32-byte slice to pad it on the left.
	copy(paddedVal[32-len(valBytes):], valBytes)

	m.Set(offset, paddedVal)
}

// Get retrieves a slice of memory of a given size from a given offset.
// It will expand the memory if reading beyond its current size.
func (m *Memory) Get(offset, size uint64) []byte {
	if size == 0 {
		return nil
	}
	// Expand memory if trying to read beyond its size.
	if requiredSize := offset + size; uint64(len(m.data)) < requiredSize {
		m.resize(requiredSize)
	}
	return m.data[offset : offset+size]
}

func (m *Memory) Display() error {
	data := m.GetData()
	memSize := len(data)

	fmt.Println("--- Memory ---")
	if memSize == 0 {
		fmt.Println("[ empty ]")
		fmt.Println("--------------")
		return errors.New("")
	}

	// Iterate through memory in 32-byte chunks.
	for i := 0; i < memSize; i += 32 {
		// Determine the end of the current line's slice.
		end := min(i+32, memSize)
		chunk := data[i:end]

		// Format each byte in the chunk as a two-character hex string.
		var hexBytes []string
		for _, b := range chunk {
			hexBytes = append(hexBytes, fmt.Sprintf("%02x", b))
		}

		// Print the memory address and the hex representation of the data.
		fmt.Printf("0x%04x:  %s\n", i, strings.Join(hexBytes, " "))
	}
	fmt.Println("--------------")

	return nil
}
