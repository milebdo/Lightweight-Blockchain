# for basic blockchain data structure
import hashlib
import json
from time import time
from uuid import uuid4
from typing import Any, Dict, List, Optional

# for Flask web app
from flask import Flask, jsonify, request

# for consensus algorithms
from urllib.parse import urlparse
import requests


###################################### Blockchain Part ######################################
"""
basic class: Blockchain
the basic data structure for blockchain
"""
class Blockchain(object):
	# Constructor
	def __init__(self):
		self.chain = []
		self.current_transactions = []
		self.nodes = set()

		# Create the genesis block
		self.new_block(previous_hash = '1', proof = 100)


	# methods for block produce
	def new_block(self, proof: int, previous_hash: Optional[str]) -> Dict[str, Any]:
		"""
		Create a new block and add to the chain
		:param proof: <int> The proof given by the PoW algorithm
		:param previous_hash: (Optional) <str> Hash of previous Block
		:return: <dict> New Block
		"""
		# construcor a block
		block = {
			'index': len(self.chain) + 1,
			'timestamp': time(),
			'transactions': self.current_transactions,
			'proof': proof,
			'previous_hash': previous_hash or self.hash(self.chain[-1]),
		}
		
		# reset tx list to default
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
		# put tx into mempool
		self.current_transactions.append({
			'sender': sender,
			'recipient': recipient,
			'amount': amount,
		})
		
		return self.last_block['index'] + 1

	
	# an util method to return hash of a block
	@staticmethod
	def hash(block: Dict[str, Any]) -> str:
		"""
		Hashes a block (sha-256)
		:param block: <dict> Block
        :return: <str>
		"""
		block_string = json.dumps(block, sort_keys = True).encode()
		return hashlib.sha256(block_string).hexdigest()

	
	# Return the last block in chain
	@property
	def last_block(self) -> Dict[str, Any]:
		return self.chain[-1]

	
	# methods for PoW
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


	# methods for consensus
	def register_node(self, address: str) -> None:
		"""
		Add the new nodes to the node list
		:param address: <str> Address of node. Eg. 'http://192.168.0.5:5000'
		:return: Node
		"""
		parsed_url = urlparse(address)
		self.nodes.add(parsed_url.netloc)


	def valid_chain(self, chain: List[Dict[str, Any]]) -> bool:
		"""
		Determine if the given chain is valid
		:param chain: <List> a given blockchain
		:return: <bool> True if valid
		"""
		last_block = chain[0]
		current_index = 1
		while current_index < len(chain):
			block = chain[current_index]
			print(f'{last_block}')
			print(f'{block}')
			print("\n-----------------------\n")
			
			# Check that the hash of the block is correct
			if block['previous_hash'] != self.hash(last_block):
				return False
			
			# Check the PoW is correct
			if not self.valid_proof(last_block['proof'], block['proof']):
				return False
			
			last_block = block
			current_index += 1
		
		return True

	
	def resolve_conficts(self) -> bool:
		"""
		when conficts, use the longest chain
		:return: if replace the chain, return True
		"""
		neighbours = self.nodes
		new_chain = Node
		max_length = len(self.chain)

		# Grab and verify the chains from all the nodes in our network
		for node in neighbours:
			response = requests.get(f'http://{node}/chain')
			if response.status_code == 200:
				length = response.json()['length']
				chain = response.json()['chain']
				if length > max_length and self.valid_chain(chain):
					max_length = length
					new_chain = chain
		
		if new_chain:
			self.chain = new_chain
			return True
		return False



###################################  Server Part #####################################
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


@app.route('/nodes/register', methods = ['POST'])
def register_nodes():
	values = request.get_json()
	nodes = values.get('nodes')
	if nodes is None:
		return "Error: please supply a valid list of nodes", 400

	for node in nodes:
		blockchain.register_node(node)

	response = {
		'message': "New nodes have been added",
		'total_nodes': list(blockchain.nodes),
	}
	return jsonify(response), 201


@app.route('/node/resolve', methods = ['GET'])
def consensus():
	replaced = blockchain.register_node()

	if replaced:
		response = {
			'message': "Our chain is replaced",
			'new_chain': blockchain.chain,
		}
	else:
		response = {
			'message': "Our chain is OK",
			'chain': blockchain.chain
		}
	return jsonify(response), 200



######################################## Main #####################################

# run the node on localhost:5000 by default
if __name__ == '__main__':
	from argparse import ArgumentParser

	parser = ArgumentParser()
	parser.add_argument('-p', '--port', default = 5000, type = int, help = 'port to listen on')
	args = parser.parse_args()
	port = args.port

	app.run(host = '127.0.0.1', port = port)
