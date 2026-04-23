package container

import (
	"time"
)

type Runtime interface {
	Start(name string) error
	Stop(name string) error
	List() ([]string, error)
	Remove(name string) error
	
	// Health monitoring methods
	CheckContainerHealth(name string) (*HealthStatus, error)
	GetAllContainerHealth() ([]HealthStatus, error)
	RestartContainer(name string) error
	StartHealthMonitoring()
}

// HealthStatus represents the health status of a container
type HealthStatus struct {
	Name      string
	Status    string
	Healthy   bool
	Error     string
	Timestamp time.Time
}
