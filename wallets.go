package main

import (
	"bytes"
)

// Wallets stores a collection of wallets
type Wallets struct {
	wallets map[string]*Wallet
}

