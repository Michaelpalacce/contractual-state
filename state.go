package main

import (
	"fmt"
)

type StateHolder interface {
	// GetState will return the map that holds the data.
	GetState() map[string]interface{}
	// GetLocks returns a slice of lock strings
	GetLocks() []string
	SetLocks([]string)
}

// WithContract creates a new ContractualState, initializing the data
// state must be a pointer
func WithContract(state StateHolder, contract Contract) (*ContractualState, error) {
	cs := &ContractualState{
		parent:   state,
		contract: contract,
		Provides: make(map[string]interface{}),
		Consumes: make(map[string]interface{}),
	}

	if err := cs.consume(state); err != nil {
		return nil, err
	}

	return cs, nil
}

type Obligation struct {
	Required bool
	Key      string

	// Used for Provides only. Will lock the key from ever being changed by other Contracts
	Lock bool
}

type Contract struct {
	// Provides is data that will be provided to the State after fulfillment
	Provides []Obligation
	// Consumes is data that is given to the ContractualState from the State on creation
	Consumes []Obligation
}

type ContractualState struct {
	// Provides is data that will be provided to the State after fulfillment
	Provides map[string]interface{}
	// Consumes is data that is given to the ContractualState from the State
	Consumes map[string]interface{}

	// parent should be a pointer
	parent   StateHolder
	contract Contract
}

// consume will fetch data from the `state` and set it in the ContractualState
// state Must be a pointer
func (s *ContractualState) consume(state StateHolder) error {
	locks := state.GetLocks()
	stateData := state.GetState()

	for _, obligation := range s.contract.Consumes {
		key := obligation.Key

		fmt.Println(locks)
		for _, lockName := range locks {
			if lockName == key {
				return fmt.Errorf("Key %s is locked", key)
			}
		}

		value, ok := stateData[key]

		if !ok && obligation.Required {
			return fmt.Errorf("Obligation %s does not exist and cannot be consumed", key)
		}

		if ok {
			s.Consumes[key] = value
		}
	}
	return nil
}

// Fulfill will atomically set the state to the parent object or fail
func (s *ContractualState) Fulfill() error {
	staging := make(map[string]interface{})
	locks := make([]string, 0)

	for _, obligation := range s.contract.Provides {
		key := obligation.Key
		value, ok := s.Provides[key]

		if obligation.Lock {
			locks = append(locks, key)
		}

		if !ok && obligation.Required {
			return fmt.Errorf("Contract was not fulfilled. %s was not set", key)
		}

		if ok {
			staging[key] = value
		}
	}

	stateData := s.parent.GetState()
	locksData := s.parent.GetLocks()

	for key, value := range staging {
		stateData[key] = value
	}

	locksData = append(locksData, locks...)

	return nil
}
