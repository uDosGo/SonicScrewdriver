package container

import (
	"fmt"
	"log"
)

// DockerRuntime is a scaffold-only Docker backend.
type DockerRuntime struct{}

func (d DockerRuntime) Start(name string) error {
	log.Printf("Starting container: %s", name)
	// TODO: Implement actual Docker start logic
	// Temporary: Return nil for testing
	return nil
}

func (d DockerRuntime) Stop(name string) error {
	log.Printf("Stopping container: %s", name)
	// TODO: Implement actual Docker stop logic
	// Temporary: Return nil for testing
	return nil
}

func (d DockerRuntime) List() ([]string, error) {
	log.Printf("Listing containers")
	// TODO: Implement actual Docker list logic
	return []string{}, fmt.Errorf("docker list not implemented")
}

func (d DockerRuntime) Remove(name string) error {
	log.Printf("Removing container: %s", name)
	// TODO: Implement actual Docker remove logic
	// Temporary: Return nil for testing
	return nil
}
