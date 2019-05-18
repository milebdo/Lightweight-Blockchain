import hashlib
import json
from time import time

class Blockchain(object):
	def __init__(self):
		self.chain = []
		self.current_transactions = []

		# Create the genesis block
		self.new_block(previous_hash = 1, proof = 100)


	def new_block(self, proof, previous_hash = Node):
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
		

	def new_transaction(self, sender, recipient, amount):
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
	def hash(block):
		"""
		Hashes a block (sha-256)
		:param block: <dict> Block
        :return: <str>
		"""
		block_string = json.dumps(block, sort_keys = True).encode()
		return hashlib.sha256(block_string).hexdigest()


	
	@property
	def last_block(self):
		# Return the last block in chain
		pass
