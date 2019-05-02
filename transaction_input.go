package main

// TXInput - sub type for tx
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}
