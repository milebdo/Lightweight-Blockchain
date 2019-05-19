package main

import (
	"fmt"
	"log"
	"strconv"
)

//================================== CLI APIs =================================//

func (cli *CLI) createBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Create!")
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. Hash: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")
		
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func(cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) && !ValidateAddress(to) {
		log.Panic("ERROR: address is no valid")
	}
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := newCoinbaseTX(from, "")
	txs := []*Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	fmt.Println("Success!")
}