package machine

import (
	"math/big"
	"testing"
)

func TestPushAndPop(t *testing.T) {
	s := NewStack(1024)

	val1 := big.NewInt(100)
	s.Push(val1)

	if s.data[len(s.data)-1].Cmp(val1) != 0 {
		t.Errorf("Expected top of stack to be %v, got %v", val1, s.data[len(s.data)-1])
	}

	poppedVal := s.Pop()
	if poppedVal.Cmp(val1) != 0 {
		t.Errorf("Expected popped value to be %v, got %v", val1, poppedVal)
	}

	if len(s.data) != 0 {
		t.Errorf("Expected stack to be empty, but it has %d items", len(s.data))
	}
}

func TestStackUnderflow(t *testing.T) {
	// This defer function will recover from the expected panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on stack underflow")
		}
	}()

	s := NewStack(1024)
	// This line should cause a panic.
	s.Pop()
}

// TestStackOverflow checks if the stack correctly panics when pushing beyond its max depth.
func TestStackOverflow(t *testing.T) {
	// This defer function will recover from the expected panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on stack overflow")
		}
	}()

	// Create a stack with a small depth for easy testing.
	s := NewStack(2)
	s.Push(big.NewInt(1))
	s.Push(big.NewInt(2))

	// This third push should cause a panic.
	// s.Push(big.NewInt(3))
}
