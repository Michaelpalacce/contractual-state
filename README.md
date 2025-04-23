# Contractual State

This is an idea of having a state object that can be updated only by specifying you can update it. Check the main.go for examples

## Proposal

- ContractualState 
	- Uses obligations to fill the data bellow
		- Each obligation will have several properties and will be able to point to another obligation, or if no child obligation is given, it is assumed to be the last one in the chain and everything under it is workable 
	- Passes in a contract, getting a new child state 
		- Provides
			- Data that the contract will provide. The data may exist and in that case it's the contract holder to decide if it should be overwritten or not. 
		- Consumes
			- Data that the contract holder will consume, what is given.
		- Locks
			- Will throw an error is another resource tries to Provide the same data later on
	- Contract is fulfilled.
		- By doing so, we commit changes to parent state
	- Contact can be denied
		- Error 
- Finalization is done by the callee, a defer on the fulfillment should be enough. 
- Contracts are atomic. They are either submitted at the end or discarded. 
	- When submitted all values are applied or none in case of error
- Pessimistic Locking is to be done where a second contract that modifies the same obligation will fail
- Locks are never persisted 
