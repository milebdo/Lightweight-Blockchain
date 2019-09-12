# Lightweight-Blockchain
A lightweight but complete blockchain.

Two version:

`./src` contains source files purely in Golang with BoltDB

`./demo` contains a simple demo written purely in Python 3.6 with Flask for Web server.

Golang version is more strict with all cryptography design and complete data structures, and use database with BoltDB. 

Python version is more lightweight and just for demostration, without checking the validation of address in cryptography part, nor use a file system database. Python version directly use the HTTP APIs in the localhost (default port 5000).


## Golang Version Environment and APIs
Require golang 10 or more.

Require `BoltDB` as database module, and `ripemd160` algorithms library:
```
go get github.com/boltdb/bolt
go get golang.org/x/crypto/ripemd160
```

Build the code:
```
go build -o ../build/lightchain ./src
```

### APIs Documentation

- Create a blockchain and send genesis block reward to ADDRESS
```
./lightchain createblockchain -address ADDRESS 
```
- Generates a new key-pair and saves it into the wallet file
```
./lightchain createwallet 
```
- Get balance of ADDRESS
```
./lightchain getbalance -address ADDRESS 
```
- Lists all addresses from the wallet file
```
./lightchain listaddresses 
```
- Print all the blocks of the blockchain
```
./lightchain printchain 
```
- Rebuilds the UTXO set
```
./lightchain reindexutxo
```
- Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set
```
./lightchain send -from FROM -to TO -amount AMOUNT -mine
```
- Start a node with ID specified in NODE_ID env. var. -miner enables mining
```
./lightchain startnode -miner ADDRESS 
```


## Demo (Python) Version Environment and HTTP APIs

The Python version is more simple, mainly focus on web server request, without establishing a full version of database.

Require Python 3.6 or more. Require `Flask` and `requests` modules.

```
pip install Flask==0.12.2 requests==2.18.4
pip install pipenv          # for multi-nodes testing 
pipenv --python=python3.6
pipenv install
```

Run single node (default port is 5000 in local host):

```
cd demo
Python3.6 blockchain.py
```

Run several nodes in virtual environment:
```
pipenv run python blockchain.py
pipenv run python blockchain.py -p 5001
pipenv run python blockchain.py -p 5002
```

### APIs Document
To get full chain, send a `GET` request:

```
curl http://localhost:5000/chain
```

To publish a transaction with a `POST` request:

```
curl -X POST -H "Content-Type: application/json" \
    -d '{"sender": "1senderAdd", "recipient": "3recipientAdd", "amount": 4}' \
    "http://localhost:5000/transactions/new"
```

To mine a block, send a `GET` request:
```
curl http://localhost:5000/mine
```

