package library

import "fmt"

type Manager struct{}

func (m Manager) Update() error {
	return fmt.Errorf("library update not implemented")
}
