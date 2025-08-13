package machine

import (
	"math/big"

	"github.com/charmbracelet/log"
)

type Stack struct {
	data     []*big.Int
	maxDepth int
}

// NewStack creates a new stack with a specified maximum depth.
// The EVM standard is a max depth of 1024.
func NewStack(maxDepth int) *Stack {
	return &Stack{
		data:     make([]*big.Int, 0),
		maxDepth: maxDepth,
	}
}

// Push adds an item to the top of the stack.
// Panics if the stack exceeds its maximum depth (stack overflow).
func (s *Stack) Push(val *big.Int) {
	if len(s.data) >= s.maxDepth {
		log.Fatal("stack overflow")
	}
	s.data = append(s.data, val)
}

// Pop removes and returns the top item from the stack.
// Panics if the stack is empty (stack underflow).
func (s *Stack) Pop() *big.Int {
	if len(s.data) == 0 {
		log.Fatal("stack underflow")
	}
	lastIndex := len(s.data) - 1
	val := s.data[lastIndex]
	s.data = s.data[:lastIndex]
	return val
}
