package main

var InstructionSet [256]Opcode

var GasCosts [256]uint64

func init() {
	// --- 0x00: Stop and Arithmetic Operations ---
	InstructionSet[0x00] = &Stop{}
	InstructionSet[0x01] = &Add{}
	InstructionSet[0x02] = &Sub{}
	InstructionSet[0x03] = &Mul{}

	// --- PUSH Operations (Unified) ---
	// Register the single Push struct for all 32 PUSH opcodes.
	for i := 0x60; i <= 0x7F; i++ {
		InstructionSet[i] = &Push{}
	}

	// ...Rest of opcodes

	// Gas Costs
	GasCosts[0x00] = 0
	GasCosts[0x01] = 3
	GasCosts[0x02] = 3
	GasCosts[0x03] = 10

	for i := 0x60; i <= 0x7F; i++ {
		GasCosts[i] = uint64(i - 0x60 + 3)
	}

}
