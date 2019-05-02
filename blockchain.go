package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const blocksBucket = "blocks"

const dbFile = "blockchain_.db"

// Blockchain - the basic type
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock add a block to the existing chain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	// open a read-only transactions
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// mining a new Block
	newBlock := NewBlock(data, lastHash)
	// update l-ley with new block
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}

// NewBlockchain initialize the new chain
func NewBlockchain() *Blockchain {
	var tip []byte
	// standard way of opening a BoltDB file
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// open a read-write(.Update) transactions
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		// if not exist, generate genesis block
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l")) // l-key: last block
		}

		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}
