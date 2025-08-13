package main

import (
	"math/big"
)

// Account represents a single account in the Ethereum world state.
type Account struct {
	Nonce   uint64
	Balance *big.Int
	Code    []byte              // Contract bytecode
	Storage map[[32]byte][]byte // Contract's persistent storage
}

func NewAccount() *Account {
	return &Account{
		Nonce:   0,
		Balance: new(big.Int),
		Code:    make([]byte, 0),
		Storage: make(map[[32]byte][]byte),
	}
}
