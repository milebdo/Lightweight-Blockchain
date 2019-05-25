# Lightweight-Blockchain
A lightweight but complete blockchain.

Two version:

`./golang-blockchain` written purely in Golang with BoltDB

`./python-blockchain` written purely in Python 3.6 with Flask for Web server.

Golang version is more strict with all cryptography design and complete data structures, and use database with BoltDB. 

Python version is more lightweight and for demostration, without checking the validation of address in cryptography part, nor use a file system database. Python version directly use the HTTP APIs in the localhost (default port 5000).


## Golang Version Environment
Require golang 10 or more.

Require `BoltDB` as database module, and `ripemd160` algorithms library:
```
go get github.com/boltdb/bolt
go get golang.org/x/crypto/ripemd160
```

Build the code:
```
cd golang-blockchain
go build .
```

### APIs Documentation

- Create a blockchain and send genesis block reward to ADDRESS
```
./golang-blockchain createblockchain -address ADDRESS 
```
- Generates a new key-pair and saves it into the wallet file
```
./golang-blockchain createwallet 
```
- Get balance of ADDRESS
```
./golang-blockchain getbalance -address ADDRESS 
```
- Lists all addresses from the wallet file
```
./golang-blockchain listaddresses 
```
- Print all the blocks of the blockchain
```
./golang-blockchain printchain 
```
- Rebuilds the UTXO set
```
./golang-blockchain reindexutxo
```
- Send AMOUNT of coins from FROM address to TO
```
./golang-blockchain send -from FROM -to TO -amount AMOUNT 
```


## Python Version Environment and HTTP API

The Python version is more simple, mainly focus on web server request, without establish a full version of database.

Requre Python 3.6 or more. Requre `Flask` and `requests` modules

```
pip install Flask==0.12.2 requests==2.18.4
pip install pipenv          # for multi-nodes testing 
pipenv --python=python3.6
pipenv install
```

Run single node (default port is 5000 in local host):

```
cd python-blockchain
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

