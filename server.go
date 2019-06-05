package main

import (
	"bytes"
	"fmt"
	"net"
	"encoding/gob"
	"log"
)
var nodeAddress string
var knownNodes = []string{"localhost:3000"}

// version is the the version of chain in each node, for longest chain
type version struct {
	Version		int
	BestHeight	int
	AddrFrom	string
}


// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	logError(err)
	defer ln.Close()

	bc := NewBlockchain(nodeID)
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		go handleConnection(conn, bc)
	}
}

func sendVersion(addr string, bc *Blockchain) {
	BestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(address, request)
}