package main

import "math/big"

// Opcode represents a single executable EVM instruction.
type Opcode interface {
	Execute(ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error
}

// Stop implements the STOP opcode (0x00).
type Stop struct{}

func (o *Stop) Execute(ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	ec.Stop()
	return nil
}

// OpAdd implements the ADD opcode (0x01).
type Add struct{}

func (o *Add) Execute(ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	sum := new(big.Int).Add(x, y)
	ec.Stack.Push(sum)
	return nil
}

// OpAdd implements the SUB opcode (0x02).
type Sub struct{}

// Push1 implements the PUSH1 opcode (0x60).
type Push1 struct{}

func (o *Push1) Execute(ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	value := ec.ReadCode(1)
	ec.Stack.Push(value)
	return nil
}

// Instruction Set

var InstructionSet [256]Opcode

var GasCosts [256]uint64

func init() {
	// --- 0x00: Stop and Arithmetic Operations ---
	InstructionSet[0x00] = &Stop{}
	InstructionSet[0x01] = &Add{}
	// --- 0x60: PUSH1 Operation ---
	InstructionSet[0x60] = &Push1{}
	// ...Rest of opcodes

	// Gas Costs
	GasCosts[0x00] = 0
	GasCosts[0x60] = 3
	GasCosts[0x60] = 3
}
