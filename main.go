package main

import (
	"fmt"
)

type State struct {
	State map[string]interface{} `json:"state"`

	locks []string
}

func (s *State) GetState() map[string]interface{} {
	if s.State == nil {
		s.State = make(map[string]interface{})
	}

	return s.State
}

func (s *State) GetLocks() *[]string {
	if s.locks == nil {
		s.locks = make([]string, 0)
	}

	return &s.locks
}

func main() {
	state := &State{
		State: map[string]interface{}{
			"test.property": "Hello",
			"test.struct": struct{ Test string }{
				Test: "123",
			},
		},
	}

	contract := Contract{
		Consumes: []Obligation{
			{Required: true, Key: "test.property"},
			{Required: true, Key: "test.struct"},
		},
		Provides: []Obligation{
			{Required: true, Key: "provided.key"},
			{Required: true, Key: "test.struct"},
			{Required: true, Key: "test.property", Lock: true},
		},
	}

	cs, err := WithContract(state, contract)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumed data before modification:", cs.Provides)

	cs.Provides["test.property"] = "Hello World"
	cs.Provides["test.struct"] = nil
	cs.Provides["provided.key"] = "Some new Key"

	fmt.Println("Consumed data after modification:", cs.Provides)
	fmt.Println("Original state (unchanged):", state.State)

	cs.Fulfill()

	fmt.Println("Original state after fulfillment:", state.State)

	overlappingContract := Contract{
		Provides: []Obligation{
			{Required: true, Key: "test.property", Lock: true},
		},
	}

	cs, err = WithContract(state, overlappingContract)
	if err != nil {
		panic(err)
	}

	fmt.Println(cs)
}
