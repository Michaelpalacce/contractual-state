package main

import (
	"fmt"
)

type State struct {
	State map[string]interface{} `json:"state"`

	locks []string
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

	parent   *State
	contract Contract
}

// consume will fetch data from the `state` and set it in the ContractualState
func (s *ContractualState) consume(state *State) error {
	for _, obligation := range s.contract.Consumes {
		key := obligation.Key

		for _, lockName := range state.locks {
			if lockName == key {
				return fmt.Errorf("Key %s is locked", key)
			}
		}

		value, ok := state.State[key]

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

	for key, value := range staging {
		s.parent.State[key] = value
	}

	s.parent.locks = append(s.parent.locks, locks...)

	return nil
}

// WithContract creates a new ContractualState, initializing the data
func (s *State) WithContract(contract Contract) (*ContractualState, error) {
	cs := &ContractualState{
		parent:   s,
		contract: contract,
		Provides: make(map[string]interface{}),
		Consumes: make(map[string]interface{}),
	}

	if err := cs.consume(s); err != nil {
		return nil, err
	}

	return cs, nil
}
