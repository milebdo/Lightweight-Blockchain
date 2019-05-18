import hashlib
import json
from time import time
from uuid import uuid4
from typing import Any, Dict, List, Optional

from textwrap import dedent
from flask import Flask, jsonify, request


"""
basic class: Blockchain
the basic data structure for blockchain
"""
class Blockchain(object):
	def __init__(self):
		self.chain = []
		self.current_transactions = []

		# Create the genesis block
		self.new_block(previous_hash = '1', proof = 100)


	def new_block(self, proof: int, previous_hash: Optional[str]) -> Dict[str, Any]:
		"""
		Create a new block and add to the chain
		:param proof: <int> The proof given by the PoW algorithm
		:param previous_hash: (Optional) <str> Hash of previous Block
		:return: <dict> New Block
		"""
		block = {
			'index': len(self.chain) + 1,
			'timestamp': time(),
			'transactions': self.current_transactions,
			'proof': proof,
			'previous_hash': previous_hash or self.hash(self.chain[-1]),
		}
		
		# reset tx list
		self.current_transactions = []
		self.chain.append(block)
		return block
		

	def new_transaction(self, sender: str, recipient: str, amount :int) -> int:
		"""
		Add new tx to the txs list for next block
		:param sender: <str> Address of the Sender
		:param recipient: <str> Address of the Recipient
		:param amount: <int> Amount
		:return: <int> The index of the Block that will hold this transaction
		"""
		self.current_transactions.append({
			'sender': sender,
			'recipient': recipient,
			'amount': amount,
		})
		
		return self.last_block['index'] + 1


	@staticmethod
	def hash(block: Dict[str, Any]) -> str:
		"""
		Hashes a block (sha-256)
		:param block: <dict> Block
        :return: <str>
		"""
		block_string = json.dumps(block, sort_keys = True).encode()
		return hashlib.sha256(block_string).hexdigest()

	
	@property
	def last_block(self) -> Dict[str, Any]:
		# Return the last block in chain
		return self.chain[-1]

	
	def proof_of_work(self, last_proof: int) -> int:
		"""
		PoW Algorithm:
		 - find p' s.t. hash(pp') begin with four 0s
		:param last_proof: <int>
		:return: <int>
		"""
		proof = 0
		while self.valid_proof(last_proof, proof) is False:
			proof += 1

		return proof


	@staticmethod
	def valid_proof(last_proof: int, proof: int) -> bool:
		"""
		Verify if hash has four 0s at the beginning
		:param last_proof: <int> Previous Proof
		:param proof: <int> Current Proof
		:return: <bool> True if correct, False if not.
		"""

		guess = f'{last_proof}{proof}'.encode()
		guess_hash = hashlib.sha256(guess).hexdigest()
		return guess_hash[:4] == "0000"



"""
Flask web API
"""
# Instantiate our Node
app = Flask(__name__)

# Generate a globally unique address for this node
node_indentifier = str(uuid4()).replace('-', '')

# Instantiate the Blockchain
blockchain = Blockchain()


@app.route('/mine', methods = ['GET'])
def mine():
	# mine a new block
	last_block = blockchain.last_block
	last_proof = last_block['proof']
	proof = blockchain.proof_of_work(last_block)

	# reward:
	blockchain.new_transaction(
		sender = "0", 	# 0 - coinbase
		recipient = node_indentifier,
		amount = 1,
	)

	block = blockchain.new_block(proof, None)

	response = {
		'message': "New block mined",
		'index': block['index'],
		'transactions': block['transactions'],
		'proof': block['proof'],
		'previous_hash': block['previous_hash'],
	}
	return jsonify(response), 200


@app.route('/transactions/new', methods = ['POST'])
def new_transaction():
	values = request.get_json()
	
	# Check that the required fields are in the POST'ed data
	required = ['sender', 'recipient', 'amount']
	if not all(k in values for k in required):
		return 'Missing values', 400

	# Create a new Transaction
	index = blockchain.new_transaction(values['sender'], values['recipient'], values['amount'])
	
	response = {'message': f'Transaction will be added to Block {index}'}
	return jsonify(response), 201



@app.route('/chain', methods = ['GET'])
def full_chain():
	response = {
		'chain': blockchain.chain,
		'length': len(blockchain.chain),
	}
	return jsonify(response), 200



# run the node on localhost:5000
if __name__ == '__main__':
	app.run(host = '0.0.0.0', port = 5000)
