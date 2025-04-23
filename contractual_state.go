package main

import (
	"fmt"
	"maps"
	"slices"
)

type StateHolder interface {
	// GetState will return the map that holds the data.
	GetState() map[string]any
	// GetLocks returns a slice of lock strings
	GetLocks() []string
	// AddLocks will add locks to the StateHolder's locks
	AddLocks(locks []string)
}

// WithContract creates a new ContractualState, initializing the data
// state must be a pointer
func WithContract(state StateHolder, contract Contract) (*ContractualState, error) {
	cs := &ContractualState{
		parent:   state,
		contract: contract,
		Provides: make(map[string]any),
		Consumes: make(map[string]any),
	}

	if err := cs.locks(state); err != nil {
		return nil, err
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
	// WillProvide is data that will be provided to the State after fulfillment
	WillProvide []Obligation
	// WillConsume is data that is given to the ContractualState from the State on creation
	WillConsume []Obligation
}

type ContractualState struct {
	// Provides is data that will be provided to the State after fulfillment
	Provides map[string]any
	// Consumes is data that is given to the ContractualState from the State
	Consumes map[string]any

	// parent should be a pointer
	parent   StateHolder
	contract Contract
}

// locks will check if the stateholder has any locks on data
func (s *ContractualState) locks(state StateHolder) error {
	locks := state.GetLocks()

	for _, obligation := range s.contract.WillProvide {
		key := obligation.Key

		if slices.Contains(locks, key) {
			return fmt.Errorf("key %s is locked", key)
		}
	}

	return nil
}

// consume will fetch data from the `state` and set it in the ContractualState
// state Must be a pointer
func (s *ContractualState) consume(state StateHolder) error {
	stateData := state.GetState()

	for _, obligation := range s.contract.WillConsume {
		key := obligation.Key
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
	staging := make(map[string]any)
	locks := make([]string, 0)

	for _, obligation := range s.contract.WillProvide {
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

	maps.Copy(stateData, staging)

	s.parent.AddLocks(locks)

	return nil
}
