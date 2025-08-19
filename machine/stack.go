package machine

import (
	"errors"
	"fmt"
	"math/big"
	"prevm/config"

	"github.com/charmbracelet/log"
)

type Stack struct {
	data     []*big.Int
	maxDepth int
}

var logger = config.Logger

// NewStack creates a new stack with a specified maximum depth.
// The EVM standard is a max depth of 1024.
func NewStack(maxDepth int) *Stack {
	return &Stack{
		data:     make([]*big.Int, 0),
		maxDepth: maxDepth,
	}
}

func (s *Stack) GetData() []*big.Int {
	return s.data
}

// Push adds an item to the top of the stack.
// Panics if the stack exceeds its maximum depth (stack overflow).
func (s *Stack) Push(val *big.Int) {
	if len(s.data) >= s.maxDepth {
		logger.Fatal("stack overflow")
	}
	s.data = append(s.data, val)
}

// Pop removes and returns the top item from the stack.
// Panics if the stack is empty (stack underflow).
func (s *Stack) Pop() *big.Int {

	if len(s.data) == 0 {
		logger.Fatal("stack underflow")
	}
	lastIndex := len(s.data) - 1
	val := s.data[lastIndex]
	s.data = s.data[:lastIndex]
	return val
}

func (s *Stack) Display() error {
	data := s.GetData()
	stackSize := len(data)

	fmt.Println("--- Stack ---")
	if stackSize == 0 {
		fmt.Println("[ empty ]")
		fmt.Println("-------------")
		log.Errorf("Stack empty.")
		return errors.New("")
	}

	// Print from top to bottom (last element to first)
	for i := stackSize - 1; i >= 0; i-- {
		// Format the big.Int as a 64-character hex string (32 bytes), left-padded with zeros.
		formattedValue := fmt.Sprintf("0x%064x", data[i])
		// Print the index from the top (0 is the top) and the value.
		fmt.Printf("[%d]: %s\n", stackSize-1-i, formattedValue)
	}
	fmt.Println("-------------")

	return nil
}

func (s *Stack) Dup(n int) {
	// Ensure the stack is deep enough for the operation.
	if len(s.data) < n {
		logger.Fatal("stack underflow on DUP operation")
	}

	indexToDup := len(s.data) - n
	val := s.data[indexToDup]

	s.Push(new(big.Int).Set(val))

	logger.Debug(s.Display())
}
