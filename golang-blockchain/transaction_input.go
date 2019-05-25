package main

import "bytes"

// TXInput - sub type for tx
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey	  []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}