package main

// StateDB represents the world state.
type StateDB struct {
	accounts map[[20]byte]*Account
}

func NewStateDB() *StateDB {
	return &StateDB{
		accounts: make(map[[20]byte]*Account),
	}
}

// Helper functions to interact with the state.
func (s *StateDB) GetAccount(addr [20]byte) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
	return NewAccount() // Return a new, empty account if it doesn't exist.
}

func (s *StateDB) SetCode(addr [20]byte, code []byte) {
	s.GetAccount(addr).Code = code
}

func (s *StateDB) SetStorage(addr [20]byte, key [32]byte, value []byte) {
	s.GetAccount(addr).Storage[key] = value
}
