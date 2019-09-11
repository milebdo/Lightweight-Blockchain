package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
)

const commandLength = 12

var nodeAddress string
var knownNodes = []string{"localhost:3000"}

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

func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	logError(err)
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	logError(err)

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payloadBestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

//
func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

type getblocks struct {
	AddrFrom string
}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff byte.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	logError(err)

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}
