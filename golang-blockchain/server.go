package main

import (
	"bytes"
	"fmt"
	"net"
)

var nodeAddress string
var knownNodes = []string{"localhost:3000"}

type version {
	Version    int
	BestHeight int
	AddrFrom   string
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
		logError(err)
		go handleConnection(conn, bc)
	}
}


//
func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}