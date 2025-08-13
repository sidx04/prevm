package machine

import (
	"bytes"
	"math/big"
	"testing"
)

// TestMemorySetAndGet tests basic writing and reading from memory.
func TestMemorySetAndGet(t *testing.T) {
	mem := NewMemory()
	value := []byte{0x01, 0x02, 0x03}

	mem.Set(10, value)

	retrieved := mem.Get(10, 3)
	if !bytes.Equal(value, retrieved) {
		t.Errorf("Expected %v, got %v", value, retrieved)
	}
}

// TestMemoryExpansionOnSet tests that memory expands correctly when
// a value is set beyond its current size.
func TestMemoryExpansionOnSet(t *testing.T) {
	mem := NewMemory()
	// Initially, memory is size 0.
	if len(mem.data) != 0 {
		t.Errorf("Expected initial memory size to be 0, got %d", len(mem.data))
	}

	// Set data at offset 32. This should trigger an expansion.
	// The required size is 32 (offset) + 3 (len) = 35 bytes.
	// The memory should expand to the next word boundary, which is 64 bytes (2 * 32).
	mem.Set(32, []byte{0xaa, 0xbb, 0xcc})

	expectedSize := 64
	if len(mem.data) != expectedSize {
		t.Errorf("Expected memory size to be %d after expansion, got %d", expectedSize, len(mem.data))
	}

	// Verify the data was written correctly.
	retrieved := mem.Get(32, 3)
	if !bytes.Equal([]byte{0xaa, 0xbb, 0xcc}, retrieved) {
		t.Errorf("Data not set correctly after expansion")
	}
}

// TestMemoryExpansionOnGet tests that memory expands correctly when
// a value is read from beyond its current size.
func TestMemoryExpansionOnGet(t *testing.T) {
	mem := NewMemory()
	// Reading from offset 64 should expand the memory to 96 bytes (3 * 32).
	mem.Get(64, 1)

	expectedSize := 96
	if len(mem.data) != expectedSize {
		t.Errorf("Expected memory size to be %d after expansion on get, got %d", expectedSize, len(mem.data))
	}
}

// TestMemorySet32 tests the padding and setting of a 32-byte word.
func TestMemorySet32(t *testing.T) {
	mem := NewMemory()
	// A small number that doesn't fill 32 bytes.
	val := big.NewInt(12345) // 0x3039 in hex

	mem.Set32(0, val)

	// The value should be padded with leading zeros to 32 bytes.
	expected := make([]byte, 32)
	expected[30] = 0x30
	expected[31] = 0x39

	retrieved := mem.Get(0, 32)
	if !bytes.Equal(expected, retrieved) {
		t.Errorf("Expected padded value %x, got %x", expected, retrieved)
	}
}

// TestMemoryZeroSizeOperations tests that get/set with size 0 are handled correctly.
func TestMemoryZeroSizeOperations(t *testing.T) {
	mem := NewMemory()

	// Set with zero-length value should not change memory.
	mem.Set(0, []byte{})
	if len(mem.data) != 0 {
		t.Errorf("Set with zero-length value should not expand memory")
	}

	// Get with size 0 should return nil and not expand memory.
	retrieved := mem.Get(100, 0)
	if retrieved != nil {
		t.Errorf("Expected nil for zero-size get, got %v", retrieved)
	}
	if len(mem.data) != 0 {
		t.Errorf("Get with zero size should not expand memory")
	}
}
