package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

type queueData struct {
	data []string
}

func createQueueFile(location string) error {
	if _, err := os.Stat(location); os.IsExist(err) {
		return fmt.Errorf("file already exists: %s", location)
	}

	return ioutil.WriteFile(location, []byte{}, 0644)
}

// loadQueue will load a queue's data
func loadQueue(location string) (q queueData, e error) {
	b, err := ioutil.ReadFile(location)
	if err != nil {
		return q, err
	}

	bLines := bytes.Split(b, []byte{'\n'})
	for _, bLine := range bLines {
		line := string(bLine)
		q.data = append(q.data, line)
	}

	return q, nil
}

func (q queueData) atIndex(i int) string {
	return q.data[i]
}
