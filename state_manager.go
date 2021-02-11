package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	stateFileLocation = fmt.Sprintf("%s/.queuer_state", os.Getenv("HOME"))
)

type state struct {
	m map[string]*queueState
}

type queueState struct {
	CurrentIndex int    `json:"currentLine"`
	Location     string `json:"location"`
}

func loadState() (s state, e error) {
	b, err := ioutil.ReadFile(stateFileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return s, fmt.Errorf("create state file `%s` first", stateFileLocation)
		}
		return s, err
	}

	if len(b) == 0 {
		fmt.Println("Forming new queue state...")
		s.m = map[string]*queueState{}
		return s, nil
	}

	err = json.Unmarshal(b, &s.m)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *state) getQueue(name string) (*queueState, bool) {
	queue, ok := s.m[name]
	return queue, ok
}

func (s *state) addQueue(name, location string) error {
	if _, exists := s.getQueue(name); exists {
		return fmt.Errorf("queue with name %s already exists", name)
	}

	s.m[name] = &queueState{
		CurrentIndex: 0,
		Location:     location,
	}

	return nil
}

func (s *state) save() error {
	b, err := json.Marshal(s.m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(stateFileLocation, b, 0644)
}
