package main

var InstructionSet [256]Opcode

var GasCosts [256]uint64

func init() {
	// --- 0x00: Stop and Arithmetic Operations ---
	InstructionSet[0x00] = &Stop{}
	InstructionSet[0x01] = &Add{}
	InstructionSet[0x02] = &Sub{}
	InstructionSet[0x03] = &Mul{}
	InstructionSet[0x08] = &AddMod{}
	InstructionSet[0x09] = &MulMod{}
	InstructionSet[0x0A] = &Exp{}
	InstructionSet[0x20] = &Keccak{}
	InstructionSet[0x30] = &Address{}
	InstructionSet[0x31] = &Balance{}
	InstructionSet[0x32] = &Origin{}

	InstructionSet[0x51] = &Mload{}
	InstructionSet[0x52] = &Mstore{}

	// --- PUSH Operations (Unified) ---
	// Register the single Push struct for all 32 PUSH opcodes.
	for i := 0x60; i <= 0x7F; i++ {
		InstructionSet[i] = &Push{}
	}

	// --- 0x80: Duplication Operations ---
	// Register the single Dup struct for all 16 DUP opcodes.
	for i := 0x80; i <= 0x8F; i++ {
		InstructionSet[i] = &Dup{}
	}

	// ...Rest of opcodes

	// Gas Costs
	GasCosts[STOP] = 0
	GasCosts[ADD] = 3
	GasCosts[SUB] = 3
	GasCosts[MUL] = 5
	GasCosts[EXP] = 10
	GasCosts[ADDRESS] = 2
	GasCosts[BALANCE] = 4
	GasCosts[MLOAD] = 3
	GasCosts[MSTORE] = 8

	for i := 0x60; i <= 0x7F; i++ {
		GasCosts[i] = uint64(i - 0x60 + 3)
	}

}
