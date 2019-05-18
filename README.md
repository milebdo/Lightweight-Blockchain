# Lightweight-Blockchain
A lightweight but complete blockchain.

Two version:

`./golang-blockchain` written purely in Golang with BoltDB

`./python-blockchain` written purely in Python 3.6


## Python Version Environment and HTTP API

Requre Python 3.6 or more. Requre `Flask` and `requests` modules

```
pip install Flask==0.12.2 requests==2.18.4
pip install pipenv          # for multi-nodes testing 
pipenv --python=python3.6
pipenv install
```

Run single node (default port is 5000 in local host):

```
Python3.6 blockchain.py
```

Run several nodes in virtual environment:
```
pipenv run python blockchain.py
pipenv run python blockchain.py -p 5001
pipenv run python blockchain.py -p 5002
```

To get full chain, send a `GET` request:

```
curl http://localhost:5000/chain
```

To publish a transaction with a `POST` request:

```
curl -X POST -H "Content-Type: application/json" \
    -d '{"sender": "1senderAdd", \
         "recipient": "3recipientAdd", \
         "amount": 4 \
        }' \
    "http://localhost:5000/transactions/new"
```

To mine a block, send a `GET` request:
```
curl http://localhost:5000/mine
```
