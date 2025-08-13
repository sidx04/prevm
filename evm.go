package main

import (
	"fmt"
)

type EVM struct {
	State    *StateDB
	BlockCtx *BlockContext
	// TxCtx    *TransactionContext
}

type RunState struct {
	PC     uint
	Opcode Opcode
}

func NewEVM(
	state *StateDB,
	blockCtx *BlockContext,
	// txCtx *TransactionContext
) *EVM {
	return &EVM{
		State:    state,
		BlockCtx: blockCtx,
		// TxCtx:    txCtx,
	}
}

// ProcessTransaction is the main entry point for running a transaction.
func (evm *EVM) ProcessTransaction(tx *Transaction, sender [20]byte) ([]byte, uint64, error) {
	// 1. Pre-validation using the full 'tx' object
	// (Nonce check, sufficient balance for gas, etc.)
	senderAccount := evm.State.GetAccount(sender)
	if senderAccount.Nonce != tx.Nonce { // Simplified nonce check
		return nil, 0, fmt.Errorf("invalid nonce")
	}

	// 2. Calculate Intrinsic Gas
	// (Gas cost for the transaction data itself before any code execution)

	/*
		intrinsicGas := calculateIntrinsicGas(tx.Data)
		if tx.GasLimit < intrinsicGas {
			return nil, 0, fmt.Errorf("intrinsic gas too low")
		}
		gasRemaining := tx.GasLimit - intrinsicGas
	*/
	var gasRemaining uint64 = 100000

	// 3. Create the initial Execution Context (the first call frame)
	var code []byte
	if tx.To != nil {
		code = evm.State.GetAccount(*tx.To).Code
	}

	initialContext := NewExecutionContext(
		code,
		tx.Data,
		gasRemaining,
	)

	txCtx := &TransactionContext{
		Origin:   sender,
		GasPrice: tx.GasPrice,
		Value:    tx.Value,
		Data:     tx.Data,
	}

	// 4. Execute the code
	returnData, err := evm.Execute(initialContext, txCtx)
	if err != nil {
		return nil, 0, err
	}

	// 5. Post-execution logic (e.g., refund remaining gas)
	gasUsed := tx.GasLimit - initialContext.Gas

	return returnData, gasUsed, nil
}

// execute runs the bytecode for a given context and returns the output data.
func (evm *EVM) Execute(ec *ExecutionContext, tx *TransactionContext) ([]byte, error) {
	// logger := config.Logger

	for !ec.Stopped {
		// Get the opcode object from the instruction set.
		opcodeObj := ec.GetOp() // Assumes GetOp returns an Opcode interface object

		if opcodeObj == nil {
			return nil, fmt.Errorf("invalid or unimplemented opcode")
		}

		// --- Gas Calculation (Crucial Missing Piece) ---
		// A real implementation would have a complex gas calculation here.
		// For now, we'll assume a simple static cost.
		gasCost := GasCosts[ec.Bytecode[ec.PC-1]] // Get cost of the opcode we just read
		if ec.Gas < gasCost {
			return nil, fmt.Errorf("out of gas")
		}
		ec.Gas -= gasCost

		// --- Execute the Opcode ---
		if err := opcodeObj.Execute(ec, evm.BlockCtx, tx); err != nil {
			return nil, err
		}
	}

	// The loop has finished, return the data from the context.
	return ec.ReturnData, nil
}
