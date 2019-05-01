package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block is the most basic type
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock construct the new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// NewGenesisBlock initialize the new genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/**
 * Serialization for BoltSB
 */

// Serialize the block for BoltDB
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock from DB
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
