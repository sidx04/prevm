package config

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func Hash(data []byte) []byte {
	hash := crypto.Keccak256(data)
	return hash
}
