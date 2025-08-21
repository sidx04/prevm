package main

import (
	"fmt"
	"math/big"
	"prevm/config"
)

// Opcode represents a single executable EVM instruction.
type Opcode interface {
	Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error
}

var logger = config.Logger

// ====================================
// --- STOP AND ARITHMETIC OPCODES ---
// ====================================
// Stop implements the STOP opcode (0x00).
type Stop struct{}

func (o *Stop) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	ec.Stop()
	return nil
}

// Sub implements the ADD opcode (0x01).
type Add struct{}

func (o *Add) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {

	logger.Debug(fmt.Sprintln(ec.Stack.Display()))

	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Add(x, y)
	ec.Stack.Push(res)

	logger.Info("Result:", "ADD", res)

	return nil
}

// Sub implements the SUB opcode (0x02).
type Sub struct{}

func (o *Sub) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	logger.Debug("Stack", "Data", ec.Stack.GetData())

	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Sub(x, y)
	ec.Stack.Push(res)

	logger.Info("Result:", "SUB", res)

	return nil
}

// Mul implements the MUL opcode (0x03).
type Mul struct{}

func (o *Mul) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Mul(x, y)
	ec.Stack.Push(res)
	return nil
}

// Div implements the DIV opcode (0x04).
type Div struct{}

func (o *Div) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Div(x, y)
	ec.Stack.Push(res)
	return nil
}

// Sdiv implements the SDIV opcode (0x05).
type Sdiv struct{}

// S256 constants for signed conversion.
var (
	// 2^255, the boundary for a signed 256-bit integer
	s256Limit = new(big.Int).Exp(big.NewInt(2), big.NewInt(255), nil)
	// 2^256, used for two's complement conversion
	maxU256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
)

func (o *Sdiv) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// 1. Pop the numerator and denominator from the stack.
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()

	// 2. Handle division by zero, which results in 0.
	if y.Sign() == 0 {
		ec.Stack.Push(new(big.Int)) // Push 0
		return nil
	}

	// 3. Convert the numbers to their signed 256-bit representation.
	// If a number is >= 2^255, it's negative. Its value is num - 2^256.
	if x.Cmp(s256Limit) >= 0 {
		x.Sub(x, maxU256)
	}
	if y.Cmp(s256Limit) >= 0 {
		y.Sub(y, maxU256)
	}

	// 4. Handle the specific edge case: the most negative number (-2^255)
	// divided by -1 results in -2^255, not 2^255 (which would overflow).
	if x.Cmp(s256Limit) == 0 && y.Cmp(big.NewInt(-1)) == 0 {
		ec.Stack.Push(s256Limit)
		return nil
	}

	// 5. Perform the signed division.
	res := new(big.Int).Div(x, y)

	// 6. Convert the result back to its 256-bit two's complement
	// representation before pushing it to the stack.
	// This is done by taking the result modulo 2^256.
	if res.Sign() < 0 {
		res.Add(res, maxU256)
	}

	ec.Stack.Push(res)
	return nil
}

// Mod implements the MOD opcode (0x06).
type Mod struct{}

func (o *Mod) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Mod(x, y)
	ec.Stack.Push(res)
	return nil
}

// Smod implements the SMOD opcode (0x07).
type Smod struct{}

func (o *Smod) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	x := ec.Stack.Pop()
	y := ec.Stack.Pop()
	res := new(big.Int).Mod(x, y)
	ec.Stack.Push(res)
	return nil
}

// AddMod implements the ADDMOD opcode (0x08).
type AddMod struct{}

func (o *AddMod) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// The items are popped in reverse order: N, then y, then x.
	N := ec.Stack.Pop()
	y := ec.Stack.Pop()
	x := ec.Stack.Pop()

	// The result of (x + y) mod 0 is defined as 0 in the EVM.
	if N.Sign() == 0 {
		ec.Stack.Push(new(big.Int)) // Push 0
		return nil
	}

	// 3. Perform the modular addition using the methods from math/big.
	sum := new(big.Int).Add(x, y)
	res := new(big.Int).Mod(sum, N)

	ec.Stack.Push(res)
	return nil
}

// MulMod implements the MULMOD opcode (0x09).
type MulMod struct{}

func (o *MulMod) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// The items are popped in reverse order: N, then y, then x.
	N := ec.Stack.Pop()
	y := ec.Stack.Pop()
	x := ec.Stack.Pop()

	// The result of (x * y) mod 0 is defined as 0 in the EVM.
	if N.Sign() == 0 {
		ec.Stack.Push(new(big.Int)) // Push 0
		return nil
	}

	// Perform the modular multiplication.
	sum := new(big.Int).Mul(x, y)
	res := new(big.Int).Mod(sum, N)

	ec.Stack.Push(res)
	return nil
}

// Exp implements the EXP opcode (0x0A).
type Exp struct{}

func (o *Exp) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// The items are popped in reverse order: N, then y, then x.
	y := ec.Stack.Pop()
	x := ec.Stack.Pop()

	res := new(big.Int).Exp(x, y, nil)

	ec.Stack.Push(res)
	logger.Debug(fmt.Sprintln(ec.Stack.Display()))

	return nil
}

// ==============
// --- SHA-3 ---
// ==============
type Keccak struct{}

func (o *Keccak) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	offset := ec.Stack.Pop().Uint64()
	size := ec.Stack.Pop().Uint64()

	data := ec.Memory.Get(offset, size)
	hash := config.Hash(data)

	ec.Stack.Push(new(big.Int).SetBytes(hash))

	logger.Debug("KECCAK", "data", fmt.Sprintf("0x%x", data), "hash", fmt.Sprintf("0x%x", hash))

	return nil
}

// =================================
// --- ENVIRONMENTAL OPERATIONS ---
// =================================
// Address (0x30)
type Address struct{}

func (o *Address) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	addr := new(big.Int).SetBytes(ec.Address[:])
	ec.Stack.Push(addr)

	logger.Debug("ADDRESS", "address", fmt.Sprintf("0x%x", addr))

	return nil
}

// Balance (0x31)
type Balance struct{}

func (o *Balance) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	var address [20]byte

	addressInt := ec.Stack.Pop()
	addressBytes := addressInt.Bytes()
	copy(address[20-len(addressBytes):], addressBytes)

	account := evm.State.GetAccount(address)
	bal := account.Balance

	logger.Debug("BALANCE", "address", fmt.Sprintf("0x%x", address), "balance", fmt.Sprintf("%d WEI", bal))

	ec.Stack.Push(bal)

	return nil
}

// Origin (0x32)
type Origin struct{}

func (o *Origin) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	addr := new(big.Int).SetBytes(tx.Origin[:])
	ec.Stack.Push(addr)

	return nil
}

// Caller (0x33)
type Caller struct{}

func (o *Caller) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	callAddr := ec.Caller

	logger.Debug("CALLER", "address", fmt.Sprintf("0x%x", callAddr))

	addrInt := new(big.Int).SetBytes(callAddr[:])

	ec.Stack.Push(addrInt)

	return nil
}

// CallValue(0x34)
type CallValue struct{}

func (o *CallValue) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	value := tx.Value

	logger.Debug("CALLVALUE", "value", fmt.Sprintf("%d WEI", value))

	ec.Stack.Push(value)

	return nil
}

// CallDataLoad (0x35)
type CallDataLoad struct{}

func (o *CallDataLoad) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	offset := ec.Stack.Pop().Uint64()

	data := new(big.Int).SetBytes(tx.Data[offset:])

	logger.Debug("CALLDATALOAD", "data", fmt.Sprintf("%X", data))

	ec.Stack.Push(data)

	return nil
}

// CallDataSize (0x36)
type CallDataSize struct{}

func (o *CallDataSize) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	offset := ec.Stack.Pop().Uint64()

	data := tx.Data[offset:]
	size := len(data)

	logger.Debug("CALLDATASIZE", "data", fmt.Sprintf("%X", data), "size", size)

	ec.Stack.Push(new(big.Int).SetInt64(int64(size)))

	return nil
}

// CallDataCopy (0x37)
type CallDataCopy struct{}

func (o *CallDataCopy) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	destOffset := ec.Stack.Pop().Uint64()
	offset := ec.Stack.Pop().Uint64()
	size := ec.Stack.Pop().Uint64()

	dataToCopy := make([]byte, size)

	calldataEnd := uint64(len(tx.Data))

	if offset < calldataEnd {
		copyEnd := min(offset+size,
			calldataEnd)
		copy(dataToCopy, tx.Data[offset:copyEnd])
	}
	// logger.Debug("CALLDATACOPY", "data", fmt.Sprintf("%x", dataToCopy))

	ec.Memory.Set(destOffset, dataToCopy)

	logger.Debug("CALLDATACOPY", "memory", ec.Memory.Display())

	return nil
}

// CodeSize (0x38)
type CodeSize struct{}

func (o *CodeSize) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	len := new(big.Int).SetUint64(uint64(len(ec.Bytecode)))
	logger.Debug("CODESIZE", "size", fmt.Sprint(len))

	ec.Stack.Push(len)

	return nil
}

// CodeCopy (0x39)
type CodeCopy struct{}

func (o *CodeCopy) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	code := ec.Bytecode

	ec.Memory.Set(0, code)

	return nil
}

// GasPrice (0x3A)
type GasPrice struct{}

func (o *GasPrice) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	gas := tx.GasPrice

	ec.Stack.Push(gas)

	return nil
}

// =========================
// --- BLOCK OPERATIONS ---
// =========================
// BlockHash (0x40)
type BlockHash struct{}

func (o *BlockHash) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	panic("unimplemented")

	return nil
}

// CoinBase (0x41)
type CoinBase struct{}

func (o *CoinBase) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	cb := block.Coinbase

	ec.Stack.Push(new(big.Int).SetBytes(cb[:]))

	return nil
}

// TimeStamp (0x42)
type TimeStamp struct{}

func (o *TimeStamp) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	time := block.Timestamp

	ec.Stack.Push(time)

	logger.Debug("TIMESTAMP", "unix", time)

	return nil
}

// Number (0x43)
type BlockNumber struct{}

func (o *BlockNumber) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	blockNumber := block.Number

	ec.Stack.Push(blockNumber)

	return nil
}

// PrevRandao (0x43)
type PrevRandao struct{}

func (o *PrevRandao) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	randao := block.Difficulty // PrevRANDAO after the Merge

	ec.Stack.Push(randao)

	return nil
}

// ChainID (0x43)
type ChainId struct{}

func (o *ChainId) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	chainId := block.ChainID

	ec.Stack.Push(chainId)

	return nil
}

// GasLimit (0x45)
type GasLimit struct{}

func (o *GasLimit) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	ec.ReturnData = block.GasLimit.Bytes()
	return nil
}

// =================================================
// --- STACK MEMORY STORAGE AND FLOW OPERATIONS ---
// =================================================
type Pop struct{}

func (o *Pop) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	ec.Stack.Pop()
	return nil
}

type Mload struct{}

func (o *Mload) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	offset := ec.Stack.Pop().Uint64()

	data := ec.Memory.Get(offset, 32)
	value := new(big.Int).SetBytes(data)

	ec.Stack.Push(value)

	return nil
}

type Mstore struct{}

func (o *Mstore) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	offset := ec.Stack.Pop().Uint64()
	value := ec.Stack.Pop()

	ec.Memory.Set32(offset, value)

	logger.Debug("Memory", ec.Memory.Display())

	return nil
}

type Pc struct{}

func (o *Pc) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	ec.Stack.Push(new(big.Int).SetUint64(ec.PC))
	return nil
}

// =====================
// --- PUSH OPCODES ---
// =====================
// Push handles all PUSH opcodes from PUSH1 to PUSH32
type Push struct{}

func (o *Push) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// The PC was already advanced by GetOp(). We look back one byte
	// to see which PUSH opcode it was.
	opValue := ec.Bytecode[ec.PC-1]
	numToRead := uint64(opValue - 0x5F) // e.g., 0x60 (PUSH1) - 0x5F = 1

	value := ec.ReadCode(numToRead)
	ec.Stack.Push(value)
	return nil
}

type Dup struct{}

func (o *Dup) Execute(evm *EVM, ec *ExecutionContext, block *BlockContext, tx *TransactionContext) error {
	// The PC was already advanced by GetOp(). We look back one byte
	// to see which DUP opcode it was.
	opValue := ec.Bytecode[ec.PC-1]

	depth := int(opValue - DUP1 + 1)

	// Call the stack's Dup method.
	ec.Stack.Dup(depth)

	return nil
}

// ==========================
// --- SYSTEM OPERATIONS ---
// ==========================
// Call implements the CALL opcode (0xf1).
type Call struct{}

// EVM Opcodes as constants
const (
	// --- 0x00: Stop and Arithmetic Operations ---
	STOP       = 0x00
	ADD        = 0x01
	MUL        = 0x02
	SUB        = 0x03
	DIV        = 0x04
	SDIV       = 0x05
	MOD        = 0x06
	SMOD       = 0x07
	ADDMOD     = 0x08
	MULMOD     = 0x09
	EXP        = 0x0a
	SIGNEXTEND = 0x0b

	// --- 0x10: Comparison & Bitwise Logic Operations ---
	LT     = 0x10
	GT     = 0x11
	SLT    = 0x12
	SGT    = 0x13
	EQ     = 0x14
	ISZERO = 0x15
	AND    = 0x16
	OR     = 0x17
	XOR    = 0x18
	NOT    = 0x19
	BYTE   = 0x1a
	SHL    = 0x1b
	SHR    = 0x1c
	SAR    = 0x1d

	// --- 0x20: Cryptographic ---
	KECCAK256 = 0x20

	// --- 0x30: Environmental Information ---
	ADDRESS        = 0x30
	BALANCE        = 0x31
	ORIGIN         = 0x32
	CALLER         = 0x33
	CALLVALUE      = 0x34
	CALLDATALOAD   = 0x35
	CALLDATASIZE   = 0x36
	CALLDATACOPY   = 0x37
	CODESIZE       = 0x38
	CODECOPY       = 0x39
	GASPRICE       = 0x3a
	EXTCODESIZE    = 0x3b
	EXTCODECOPY    = 0x3c
	RETURNDATASIZE = 0x3d
	RETURNDATACOPY = 0x3e
	EXTCODEHASH    = 0x3f

	// --- 0x40: Block Information ---
	BLOCKHASH   = 0x40
	COINBASE    = 0x41
	TIMESTAMP   = 0x42
	NUMBER      = 0x43
	DIFFICULTY  = 0x44
	GASLIMIT    = 0x45
	CHAINID     = 0x46
	SELFBALANCE = 0x47
	BASEFEE     = 0x48

	// --- 0x50: Stack, Memory, Storage and Flow Operations ---
	POP      = 0x50
	MLOAD    = 0x51
	MSTORE   = 0x52
	MSTORE8  = 0x53
	SLOAD    = 0x54
	SSTORE   = 0x55
	JUMP     = 0x56
	JUMPI    = 0x57
	PC       = 0x58
	MSIZE    = 0x59
	GAS      = 0x5a
	JUMPDEST = 0x5b

	// --- 0x60 & 0x70: Push Operations ---
	PUSH1  = 0x60
	PUSH2  = 0x61
	PUSH3  = 0x62
	PUSH4  = 0x63
	PUSH5  = 0x64
	PUSH6  = 0x65
	PUSH7  = 0x66
	PUSH8  = 0x67
	PUSH9  = 0x68
	PUSH10 = 0x69
	PUSH11 = 0x6a
	PUSH12 = 0x6b
	PUSH13 = 0x6c
	PUSH14 = 0x6d
	PUSH15 = 0x6e
	PUSH16 = 0x6f
	PUSH17 = 0x70
	PUSH18 = 0x71
	PUSH19 = 0x72
	PUSH20 = 0x73
	PUSH21 = 0x74
	PUSH22 = 0x75
	PUSH23 = 0x76
	PUSH24 = 0x77
	PUSH25 = 0x78
	PUSH26 = 0x79
	PUSH27 = 0x7a
	PUSH28 = 0x7b
	PUSH29 = 0x7c
	PUSH30 = 0x7d
	PUSH31 = 0x7e
	PUSH32 = 0x7f

	// --- 0x80: Duplication Operations ---
	DUP1  = 0x80
	DUP2  = 0x81
	DUP3  = 0x82
	DUP4  = 0x83
	DUP5  = 0x84
	DUP6  = 0x85
	DUP7  = 0x86
	DUP8  = 0x87
	DUP9  = 0x88
	DUP10 = 0x89
	DUP11 = 0x8a
	DUP12 = 0x8b
	DUP13 = 0x8c
	DUP14 = 0x8d
	DUP15 = 0x8e
	DUP16 = 0x8f

	// --- 0x90: Swap Operations ---
	SWAP1  = 0x90
	SWAP2  = 0x91
	SWAP3  = 0x92
	SWAP4  = 0x93
	SWAP5  = 0x94
	SWAP6  = 0x95
	SWAP7  = 0x96
	SWAP8  = 0x97
	SWAP9  = 0x98
	SWAP10 = 0x99
	SWAP11 = 0x9a
	SWAP12 = 0x9b
	SWAP13 = 0x9c
	SWAP14 = 0x9d
	SWAP15 = 0x9e
	SWAP16 = 0x9f

	// --- 0xa0: Logging Operations ---
	LOG0 = 0xa0
	LOG1 = 0xa1
	LOG2 = 0xa2
	LOG3 = 0xa3
	LOG4 = 0xa4

	// --- 0xf0: System Operations ---
	CREATE       = 0xf0
	CALL         = 0xf1
	CALLCODE     = 0xf2
	RETURN       = 0xf3
	DELEGATECALL = 0xf4
	CREATE2      = 0xf5
	STATICCALL   = 0xfa
	REVERT       = 0xfd
	SELFDESTRUCT = 0xff
)
