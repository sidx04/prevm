package main

import (
	"math/big"
	"prevm/machine"
)

// ExecutionContext holds the state for the current execution scope.
type ExecutionContext struct {
	// Address is the The address of the contract whose code is currently running.
	// This value changes with every CALL or DELEGATECALL.
	Address [20]byte
	// Bytecode is the contract bytecode to be executed in this context.
	Bytecode []byte
	// Address of the account that called this context
	Caller [20]byte
	// Stack is the stack for this context.
	Stack *machine.Stack
	// Memory is the memory for this context.
	Memory *machine.Memory
	// PC is the program counter, pointing to the current instruction.
	PC uint64
	// Gas available for this context/frame
	Gas uint64

	// Value (in Wei) passed with this call
	CallValue *big.Int
	// Input data for this context/frame
	CallData []byte

	// Stopped indicates whether execution has been halted.
	Stopped bool
	// ReturnData is the data returned from this execution context.
	ReturnData []byte
	// True if this is a STATICCALL context
	IsStatic bool
}

// BlockContext holds information about the current block.
// This data is accessible to contracts via specific opcodes.
type BlockContext struct {
	// BaseFee (EIP-1559) is the base fee per gas. Accessible via BASEFEE opcode.
	BaseFee *big.Int
	// Coinbase is the address of the block producer (miner/validator). Accessible via COINBASE opcode.
	Coinbase [20]byte // Ethereum addresses are 20 bytes
	// Timestamp is the block's timestamp. Accessible via TIMESTAMP opcode.
	Timestamp *big.Int
	// Number is the current block number. Accessible via NUMBER opcode.
	Number *big.Int
	// Difficulty (or PrevRandao for post-Merge) is the difficulty of the block. Accessible via DIFFICULTY opcode.
	Difficulty *big.Int
	// GasLimit is the gas limit for the entire block. Accessible via GASLIMIT opcode.
	GasLimit *big.Int
	// ChainID identifies the specific chain. Accessible via CHAINID opcode.
	ChainID *big.Int
}

// This is the full transaction object, with the nonce.
type Transaction struct {
	Nonce    uint64
	GasLimit uint64
	GasPrice *big.Int
	To       *[20]byte
	Value    *big.Int
	Data     []byte
	// ... V, R, S for signature later
}

// TransactionContext holds information specific to the transaction being processed.
// This data is generally immutable during the execution of the transaction.
type TransactionContext struct {
	// Origin is the address of the Externally Owned Account (EOA) that
	// originally signed and sent the transaction. Accessible via the ORIGIN (0x32) opcode.
	Origin [20]byte

	// GasPrice is the price per unit of gas (in Wei) that the sender is
	// willing to pay for the transaction. Accessible via the GASPRICE (0x3a) opcode.
	GasPrice *big.Int

	// To is the recipient address of the transaction.
	// This will be 'nil' if the transaction is a contract creation.
	// To *[20]byte

	// Value is the amount of Ether (in Wei) transferred with this transaction.
	// Accessible via the CALLVALUE (0x34) opcode.
	Value *big.Int

	// Data is the input data of the transaction, also known as "calldata".
	// This is accessed by opcodes like CALLDATALOAD (0x35), CALLDATASIZE (0x36),
	// and CALLDATACOPY (0x37).
	Data []byte
}

// NewExecutionContext creates a new execution context.
func NewExecutionContext(caller, address [20]byte, bytecode []byte, calldata []byte, value *big.Int, gas uint64) *ExecutionContext {
	return &ExecutionContext{
		Address:   address,
		Bytecode:  bytecode,
		Stack:     machine.NewStack(1024),
		Memory:    machine.NewMemory(),
		PC:        0,
		Caller:    caller,
		CallValue: value,
		Gas:       gas,
		CallData:  calldata,
		Stopped:   false,
		IsStatic:  false,
	}
}

// Stop halts the execution of the current context.
func (ec *ExecutionContext) Stop() {
	ec.Stopped = true
}

// ReadCode reads a specified number of bytes from the code buffer
// and advances the program counter. It returns the bytes as a big.Int.
func (ec *ExecutionContext) ReadCode(numBytes uint64) *big.Int {
	// Ensure we don't read past the end of the code.
	if ec.PC+numBytes > uint64(len(ec.Bytecode)) {
		// In a real EVM, this might be handled differently, but for now,
		// we return 0 to avoid a panic.
		// Opcodes like PUSH expect bytes to be there. If they aren't,
		// the EVM specification says they should be treated as zeros.
		// We advance the PC to the end of the code.
		ec.PC = uint64(len(ec.Bytecode))
		return new(big.Int)
	}

	// Read the bytes from the code.
	data := ec.Bytecode[ec.PC : ec.PC+numBytes]
	value := new(big.Int).SetBytes(data)

	// Advance the program counter.
	ec.PC += numBytes

	return value
}

// GetOp reads the opcode at the current PC, advances the PC, and returns the opcode.
// It handles the edge case of reading past the end of the code.
func (ec *ExecutionContext) GetOp() Opcode {
	var opByte byte
	// If PC is out of bounds, return STOP.
	if ec.PC >= uint64(len(ec.Bytecode)) {
		return &Stop{} // OpCode for 0x00
	} else {
		opByte = ec.Bytecode[ec.PC]
	}

	// Advance the program counter.
	ec.PC++
	return InstructionSet[opByte]
}
