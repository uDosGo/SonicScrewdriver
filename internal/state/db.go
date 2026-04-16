package state

import "fmt"

type DB struct{}

func Open(path string) (*DB, error) {
	return nil, fmt.Errorf("sqlite backend not implemented: %s", path)
}
