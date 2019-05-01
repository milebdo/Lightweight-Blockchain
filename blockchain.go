package main

// Blockchain - the basic type
type Blockchain struct {
	blocks []*Block
}

// AddBlock add a block to the existing chain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain initialize the new chain
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
