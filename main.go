package main

import (
	"fmt"
)

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
		WillConsume: []Obligation{
			{Required: true, Key: "test.property"},
			{Required: true, Key: "test.struct"},
		},
		WillProvide: []Obligation{
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
		WillProvide: []Obligation{
			{Required: true, Key: "test.property", Lock: true},
		},
	}

	cs, err = WithContract(state, overlappingContract)
	if err != nil {
		panic(err)
	}

	fmt.Println(cs)
}
