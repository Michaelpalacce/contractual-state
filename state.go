package main

// State is a simple implementation of the StateHolder interface
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

func (s *State) GetLocks() []string {
	if s.locks == nil {
		s.locks = make([]string, 0)
	}

	return s.locks
}

func (s *State) AddLocks(locks []string) {
	s.locks = append(s.locks, locks...)
}
