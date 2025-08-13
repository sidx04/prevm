package main

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
		Prefix:          "EVM",
	})
	logger.Info("Hello, Ethereum!")

	logger.Info("Starting EVM simulation...")

	// 2. --- EVM Component Setup ---
	state := NewStateDB()
	logger.Debug("Initialized StateDB")

	// --- Define Addresses ---
	accountA_Addr := [20]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}
	accountB_Addr := [20]byte{0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd, 0xdd}
	contractAddr := [20]byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}

	// --- Setup Accounts in StateDB ---
	// Account A starts with nonce 0 and 1 ETH.
	accountA := NewAccount()
	accountA.Balance = new(big.Int).SetUint64(1000000000000000000) // 1 ETH
	accountA.Nonce = 0
	state.accounts[accountA_Addr] = accountA
	logger.Debug("Created Account A", "address", fmt.Sprintf("0x%x", accountA_Addr))

	// Account B starts with nonce 0 and 1 ETH.
	accountB := NewAccount()
	accountB.Balance = new(big.Int).SetUint64(1000000000000000000) // 1 ETH
	accountB.Nonce = 0
	state.accounts[accountB_Addr] = accountB
	logger.Debug("Created Account B", "address", fmt.Sprintf("0x%x", accountB_Addr))

	// --- Contract Bytecode (calculates 10 + 5 = 15) ---
	bytecode := []byte{byte(0x60), 0x0A, byte(0x60), 0x05, byte(0x01), byte(0x00)}

	contractAccount := NewAccount()
	contractAccount.Code = bytecode
	state.accounts[contractAddr] = contractAccount
	logger.Debug("Created contract account", "address", fmt.Sprintf("0x%x", contractAddr))

	// --- Block Context and EVM Setup ---
	blockCtx := &BlockContext{Number: big.NewInt(1)}
	evm := NewEVM(state, blockCtx)
	logger.Info("EVM Initialized and all accounts are set up.")

	// ===================================================================
	// --- TRANSACTION 1: Account A calls the contract ---
	// ===================================================================
	logger.Info("--- Processing Tx 1: Account A calls contract ---")
	tx1 := &Transaction{
		Nonce:    0, // Account A's first transaction
		GasLimit: 100000,
		To:       &contractAddr,
	}

	_, gasUsed1, err1 := evm.ProcessTransaction(tx1, accountA_Addr)
	if err1 != nil {
		logger.Error("EVM execution failed for Tx 1", "error", err1)
	} else {
		logger.Info("Tx 1 successful!")
		logger.Info("Gas Used", "amount", gasUsed1)
		logger.Info("Account A Nonce after Tx 1", "nonce", state.GetAccount(accountA_Addr).Nonce)
	}

	fmt.Println() // Add a blank line for readability

	// ===================================================================
	// --- TRANSACTION 2: Account B calls the same contract ---
	// ===================================================================
	logger.Info("--- Processing Tx 2: Account B calls contract ---")
	tx2 := &Transaction{
		Nonce:    0, // Account B's first transaction
		GasLimit: 100000,
		To:       &contractAddr,
	}

	_, gasUsed2, err2 := evm.ProcessTransaction(tx2, accountB_Addr)
	if err2 != nil {
		logger.Error("EVM execution failed for Tx 2", "error", err2)
	} else {
		logger.Info("Tx 2 successful!")
		logger.Info("Gas Used", "amount", gasUsed2)
		logger.Info("Account B Nonce after Tx 2", "nonce", state.GetAccount(accountB_Addr).Nonce)
	}
}
