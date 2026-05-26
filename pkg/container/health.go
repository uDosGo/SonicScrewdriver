package container

import (
	"context"
	"fmt"
	"log"
	"time"

	dockerContainer "github.com/docker/docker/api/types/container"
	dockerFilters "github.com/docker/docker/api/types/filters"
)

// HealthCheckInterval defines how often to check container health
const HealthCheckInterval = 30 * time.Second

// CheckContainerHealth checks the health status of a container
func (d *DockerRuntime) CheckContainerHealth(name string) (*HealthStatus, error) {
	ctx := context.Background()

	// Find the container
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("name", name)),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		return &HealthStatus{
			Name:      name,
			Status:    "not_found",
			Healthy:   false,
			Error:     "Container not found",
			Timestamp: time.Now(),
		}, nil
	}

	container := containers[0]
	containerID := container.ID

	// Check container state
	if container.State != "running" {
		return &HealthStatus{
			Name:      name,
			Status:    container.State,
			Healthy:   false,
			Error:     fmt.Sprintf("Container is %s", container.State),
			Timestamp: time.Now(),
		}, nil
	}

	// Inspect container for detailed health information
	inspect, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	// Check health status if health check is configured
	if inspect.State.Health != nil {
		return &HealthStatus{
			Name:      name,
			Status:    inspect.State.Health.Status,
			Healthy:   inspect.State.Health.Status == "healthy",
			Error:     inspect.State.Health.Log[len(inspect.State.Health.Log)-1].Output,
			Timestamp: time.Now(),
		}, nil
	}

	// Basic health check - container is running
	return &HealthStatus{
		Name:      name,
		Status:    "running",
		Healthy:   true,
		Error:     "",
		Timestamp: time.Now(),
	}, nil
}

// MonitorContainerHealth starts a health monitoring goroutine
func (d *DockerRuntime) MonitorContainerHealth(name string, recoveryCallback func(string)) {
	go func() {
		for {
			status, err := d.CheckContainerHealth(name)
			if err != nil {
				log.Printf("Health check error for %s: %v", name, err)
			} else if !status.Healthy {
				log.Printf("Container %s unhealthy: %s (status: %s)", name, status.Error, status.Status)
				if recoveryCallback != nil {
					recoveryCallback(name)
				}
			}
			
			// Wait for next health check
			time.Sleep(HealthCheckInterval)
		}
	}()
}

// StartHealthMonitoring starts health monitoring for all Sonic containers
func (d *DockerRuntime) StartHealthMonitoring() {
	ctx := context.Background()

	// Find all Sonic containers
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("label", "com.sonic.game")),
	})

	if err != nil {
		log.Printf("Failed to list containers for health monitoring: %v", err)
		return
	}

	log.Printf("Starting health monitoring for %d containers", len(containers))

	for _, container := range containers {
		// Extract name without leading slash
		name := container.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		// Start monitoring for each container
		d.MonitorContainerHealth(name, func(containerName string) {
			log.Printf("Attempting to recover container: %s", containerName)
			
			// Try to restart the container
			err := d.RestartContainer(containerName)
			if err != nil {
				log.Printf("Failed to restart container %s: %v", containerName, err)
			} else {
				log.Printf("Successfully restarted container: %s", containerName)
			}
		})
	}
}

// RestartContainer attempts to restart a container
func (d *DockerRuntime) RestartContainer(name string) error {
	ctx := context.Background()

	// Find the container
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("name", name)),
	})

	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		return fmt.Errorf("container %s not found", name)
	}

	containerID := containers[0].ID
	
	// Stop the container
	timeout := 10 // seconds
	if err := d.client.ContainerStop(ctx, containerID, dockerContainer.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Start the container
	if err := d.client.ContainerStart(ctx, containerID, dockerContainer.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Restarted container: %s (ID: %s)", name, containerID)
	return nil
}

// GetAllContainerHealth gets health status for all Sonic containers
func (d *DockerRuntime) GetAllContainerHealth() ([]HealthStatus, error) {
	ctx := context.Background()

	// Find all Sonic containers
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("label", "com.sonic.game")),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var healthStatuses []HealthStatus

	for _, container := range containers {
		// Extract name without leading slash
		name := container.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		// Get health status
		status, err := d.CheckContainerHealth(name)
		if err != nil {
			log.Printf("Failed to check health for %s: %v", name, err)
			status = &HealthStatus{
				Name:      name,
				Status:    "error",
				Healthy:   false,
				Error:     err.Error(),
				Timestamp: time.Now(),
			}
		}

		healthStatuses = append(healthStatuses, *status)
	}

	return healthStatuses, nil
}
