package main

import (
	"github.com/boltdb/bolt"
)

// BlockchainIterator stores current hash and attached to a blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator over blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// Next return the next block from tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	logError(err)

	i.currentHash = block.PrevBlockHash
	return block
}
