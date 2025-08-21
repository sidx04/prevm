package main

var InstructionSet [256]Opcode

var GasCosts [256]uint64

func init() {
	// --- 0x00: Stop and Arithmetic Operations ---
	InstructionSet[STOP] = &Stop{}
	InstructionSet[ADD] = &Add{}
	InstructionSet[MUL] = &Mul{}
	InstructionSet[SUB] = &Sub{}
	InstructionSet[DIV] = &Div{}
	InstructionSet[SDIV] = &Sdiv{}
	InstructionSet[MOD] = &Mod{}
	InstructionSet[SMOD] = &Smod{}
	InstructionSet[ADDMOD] = &AddMod{}
	InstructionSet[MULMOD] = &MulMod{}
	InstructionSet[EXP] = &Exp{}
	// InstructionSet[SIGNEXTEND] = &SignExtend{}

	// --- 0x10: Comparison & Bitwise Logic Operations ---
	// InstructionSet[LT] = &Lt{}
	// InstructionSet[GT] = &Gt{}
	// InstructionSet[SLT] = &Slt{}
	// InstructionSet[SGT] = &Sgt{}
	// InstructionSet[EQ] = &Eq{}
	// InstructionSet[ISZERO] = &IsZero{}
	// InstructionSet[AND] = &And{}
	// InstructionSet[OR] = &Or{}
	// InstructionSet[XOR] = &Xor{}
	// InstructionSet[NOT] = &Not{}
	// InstructionSet[BYTE] = &Byte{}
	// InstructionSet[SHL] = &Shl{}
	// InstructionSet[SHR] = &Shr{}
	// InstructionSet[SAR] = &Sar{}

	// --- 0x20: Cryptographic ---
	InstructionSet[KECCAK256] = &Keccak{}

	// --- 0x30: Environmental Information ---
	InstructionSet[ADDRESS] = &Address{}
	InstructionSet[BALANCE] = &Balance{}
	InstructionSet[ORIGIN] = &Origin{}
	InstructionSet[CALLER] = &Caller{}
	InstructionSet[CALLVALUE] = &CallValue{}
	InstructionSet[CALLDATALOAD] = &CallDataLoad{}
	InstructionSet[CALLDATASIZE] = &CallDataSize{}
	InstructionSet[CALLDATACOPY] = &CallDataCopy{}
	InstructionSet[CODESIZE] = &CodeSize{}
	InstructionSet[CODECOPY] = &CodeCopy{}
	InstructionSet[GASPRICE] = &GasPrice{}
	// InstructionSet[EXTCODESIZE] = &ExtCodeSize{}
	// InstructionSet[EXTCODECOPY] = &ExtCodeCopy{}
	// InstructionSet[RETURNDATASIZE] = &ReturnDataSize{}
	// InstructionSet[RETURNDATACOPY] = &ReturnDataCopy{}
	// InstructionSet[EXTCODEHASH] = &ExtCodeHash{}

	// --- 0x40: Block Information ---
	InstructionSet[BLOCKHASH] = &BlockHash{}
	InstructionSet[COINBASE] = &CoinBase{}
	InstructionSet[TIMESTAMP] = &TimeStamp{}
	InstructionSet[NUMBER] = &BlockNumber{}
	InstructionSet[DIFFICULTY] = &PrevRandao{}
	InstructionSet[GASLIMIT] = &GasLimit{}
	InstructionSet[CHAINID] = &ChainId{}
	// InstructionSet[SELFBALANCE] = &SelfBalance{}
	// InstructionSet[BASEFEE] = &BaseFee{}

	// --- 0x50: Stack, Memory, Storage and Flow Operations ---
	InstructionSet[POP] = &Pop{}
	InstructionSet[MLOAD] = &Mload{}
	InstructionSet[MSTORE] = &Mstore{}
	// InstructionSet[MSTORE8] = &Mstore8{}
	// InstructionSet[SLOAD] = &Sload{}
	// InstructionSet[SSTORE] = &Sstore{}
	// InstructionSet[JUMP] = &Jump{}
	// InstructionSet[JUMPI] = &Jumpi{}
	// InstructionSet[PC] = &Pc{}
	// InstructionSet[MSIZE] = &Msize{}
	// InstructionSet[GAS] = &Gas{}
	// InstructionSet[JUMPDEST] = &JumpDest{}

	// --- 0x60 & 0x70: Push Operations (Unified) ---
	for i := 0x60; i <= 0x7F; i++ {
		InstructionSet[i] = &Push{}
	}

	// --- 0x80: Duplication Operations (Unified) ---
	for i := 0x80; i <= 0x8F; i++ {
		InstructionSet[i] = &Dup{}
	}

	// --- 0x90: Swap Operations (Unified) ---
	// for i := 0x90; i <= 0x9F; i++ {
	// 	InstructionSet[i] = &Swap{}
	// }

	// --- 0xa0: Logging Operations (Unified) ---
	// for i := 0xa0; i <= 0xa4; i++ {
	// 	InstructionSet[i] = &Log{}
	// }

	// --- 0xf0: System Operations ---
	// InstructionSet[CREATE] = &Create{}
	// InstructionSet[CALL] = &Call{}
	// InstructionSet[CALLCODE] = &CallCode{}
	// InstructionSet[RETURN] = &Return{}
	// InstructionSet[DELEGATECALL] = &DelegateCall{}
	// InstructionSet[CREATE2] = &Create2{}
	// InstructionSet[STATICCALL] = &StaticCall{}
	// InstructionSet[REVERT] = &Revert{}
	// InstructionSet[SELFDESTRUCT] = &SelfDestruct{}

	// ===================================================================
	// --- Gas Costs (Static Minimums) ---
	// ===================================================================

	GasCosts[STOP] = 0
	GasCosts[ADD] = 3
	GasCosts[MUL] = 5
	GasCosts[SUB] = 3
	GasCosts[DIV] = 5
	GasCosts[SDIV] = 5
	GasCosts[MOD] = 5
	GasCosts[SMOD] = 5
	GasCosts[ADDMOD] = 8
	GasCosts[MULMOD] = 8
	GasCosts[EXP] = 10
	GasCosts[SIGNEXTEND] = 5
	GasCosts[LT] = 3
	GasCosts[GT] = 3
	GasCosts[SLT] = 3
	GasCosts[SGT] = 3
	GasCosts[EQ] = 3
	GasCosts[ISZERO] = 3
	GasCosts[AND] = 3
	GasCosts[OR] = 3
	GasCosts[XOR] = 3
	GasCosts[NOT] = 3
	GasCosts[BYTE] = 3
	GasCosts[SHL] = 3
	GasCosts[SHR] = 3
	GasCosts[SAR] = 3
	GasCosts[KECCAK256] = 30
	GasCosts[ADDRESS] = 2
	GasCosts[BALANCE] = 100
	GasCosts[ORIGIN] = 2
	GasCosts[CALLER] = 2
	GasCosts[CALLVALUE] = 2
	GasCosts[CALLDATALOAD] = 3
	GasCosts[CALLDATASIZE] = 2
	GasCosts[CALLDATACOPY] = 3
	GasCosts[CODESIZE] = 2
	GasCosts[CODECOPY] = 3
	GasCosts[GASPRICE] = 2
	GasCosts[EXTCODESIZE] = 100
	GasCosts[EXTCODECOPY] = 100
	GasCosts[RETURNDATASIZE] = 2
	GasCosts[RETURNDATACOPY] = 3
	GasCosts[EXTCODEHASH] = 100
	GasCosts[BLOCKHASH] = 20
	GasCosts[COINBASE] = 2
	GasCosts[TIMESTAMP] = 2
	GasCosts[NUMBER] = 2
	GasCosts[DIFFICULTY] = 2
	GasCosts[GASLIMIT] = 2
	GasCosts[CHAINID] = 2
	GasCosts[SELFBALANCE] = 5
	GasCosts[BASEFEE] = 2
	GasCosts[POP] = 2
	GasCosts[MLOAD] = 3
	GasCosts[MSTORE] = 3
	GasCosts[MSTORE8] = 3
	GasCosts[SLOAD] = 100
	GasCosts[SSTORE] = 100
	GasCosts[JUMP] = 8
	GasCosts[JUMPI] = 10
	GasCosts[PC] = 2
	GasCosts[MSIZE] = 2
	GasCosts[GAS] = 2
	GasCosts[JUMPDEST] = 1
	GasCosts[CREATE] = 32000
	GasCosts[CALL] = 100
	GasCosts[CALLCODE] = 100
	GasCosts[RETURN] = 0
	GasCosts[DELEGATECALL] = 100
	GasCosts[CREATE2] = 32000
	GasCosts[STATICCALL] = 100
	GasCosts[REVERT] = 0
	GasCosts[SELFDESTRUCT] = 5000

	// Gas costs for PUSH, DUP, and SWAP are all 3
	for i := 0x60; i <= 0x7F; i++ {
		GasCosts[i] = 3
	}
	for i := 0x80; i <= 0x8F; i++ {
		GasCosts[i] = 3
	}
	for i := 0x90; i <= 0x9F; i++ {
		GasCosts[i] = 3
	}

	// Gas costs for LOG opcodes
	GasCosts[LOG0] = 375
	GasCosts[LOG1] = 375 * 2
	GasCosts[LOG2] = 375 * 3
	GasCosts[LOG3] = 375 * 4
	GasCosts[LOG4] = 375 * 5
}
