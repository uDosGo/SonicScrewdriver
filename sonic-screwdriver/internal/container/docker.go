package container

import "fmt"

// DockerRuntime is a scaffold-only Docker backend.
type DockerRuntime struct{}

func (d DockerRuntime) Start(name string) error {
	return fmt.Errorf("docker start not implemented for %s", name)
}

func (d DockerRuntime) Stop(name string) error {
	return fmt.Errorf("docker stop not implemented for %s", name)
}
