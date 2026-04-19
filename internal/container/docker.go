package container

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	dockerContainer "github.com/docker/docker/api/types/container"
	dockerFilters "github.com/docker/docker/api/types/filters"
	dockerImage "github.com/docker/docker/api/types/image"
)

// DockerRuntime implements the Runtime interface using Docker SDK
type DockerRuntime struct {
	client *client.Client
}

func NewDockerRuntime() (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &DockerRuntime{client: cli}, nil
}

func (d *DockerRuntime) Start(name string) error {
	log.Printf("Starting container: %s", name)
	
	ctx := context.Background()
	
	// Check if container exists
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("name", name)),
	})
	
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}
	
	if len(containers) == 0 {
		// Container doesn't exist, try to create and start it
		return d.createAndStartContainer(ctx, name)
	}
	
	// Container exists, just start it
	containerID := containers[0].ID
	if containers[0].State == "running" {
		log.Printf("Container %s is already running", name)
		return nil
	}
	
	if err := d.client.ContainerStart(ctx, containerID, dockerContainer.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %w", name, err)
	}
	
	log.Printf("Started container: %s (ID: %s)", name, containerID)
	return nil
}

func (d *DockerRuntime) createAndStartContainer(ctx context.Context, name string) error {
	// Use a default image for now - in production this would come from the game manifest
	imageName := "alpine:latest"
	
	// Pull the image
	_, _, err := d.client.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			log.Printf("Pulling image: %s", imageName)
			_, err = d.client.ImagePull(ctx, imageName, dockerImage.PullOptions{})
			if err != nil {
				return fmt.Errorf("failed to pull image %s: %w", imageName, err)
			}
		} else {
			return fmt.Errorf("failed to inspect image %s: %w", imageName, err)
		}
	}
	
	// Create container
	resp, err := d.client.ContainerCreate(ctx, &dockerContainer.Config{
		Image: imageName,
		Cmd:  []string{"sleep", "infinity"}, // Keep container running
		Labels: map[string]string{
			"com.sonic.game": name,
		},
	}, &dockerContainer.HostConfig{
		AutoRemove: false,
	}, nil, nil, name)
	
	if err != nil {
		return fmt.Errorf("failed to create container %s: %w", name, err)
	}
	
	// Start container
	if err := d.client.ContainerStart(ctx, resp.ID, dockerContainer.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %w", name, err)
	}
	
	log.Printf("Created and started container: %s (ID: %s)", name, resp.ID)
	return nil
}

func (d *DockerRuntime) Stop(name string) error {
	log.Printf("Stopping container: %s", name)
	
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
		log.Printf("Container %s not found", name)
		return nil // Not an error if container doesn't exist
	}
	
	containerID := containers[0].ID
	if containers[0].State != "running" {
		log.Printf("Container %s is already stopped", name)
		return nil
	}
	
	// Stop the container
	timeout := 10 // seconds
	if err := d.client.ContainerStop(ctx, containerID, dockerContainer.StopOptions{Timeout: &timeout}); err != nil {
		if !errdefs.IsNotFound(err) {
			return fmt.Errorf("failed to stop container %s: %w", name, err)
		}
	}
	
	log.Printf("Stopped container: %s (ID: %s)", name, containerID)
	return nil
}

func (d *DockerRuntime) List() ([]string, error) {
	log.Printf("Listing containers")
	
	ctx := context.Background()
	
	containers, err := d.client.ContainerList(ctx, dockerContainer.ListOptions{
		All: true,
		Filters: dockerFilters.NewArgs(dockerFilters.Arg("label", "com.sonic.game")),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	
	var names []string
	for _, c := range containers {
		// Extract the name without the leading slash
		name := strings.TrimPrefix(c.Names[0], "/")
		names = append(names, name)
	}
	
	return names, nil
}

func (d *DockerRuntime) Remove(name string) error {
	log.Printf("Removing container: %s", name)
	
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
		log.Printf("Container %s not found", name)
		return nil // Not an error if container doesn't exist
	}
	
	containerID := containers[0].ID
	
	// Remove the container
	options := dockerContainer.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}
	
	if err := d.client.ContainerRemove(ctx, containerID, options); err != nil {
		if !errdefs.IsNotFound(err) {
			return fmt.Errorf("failed to remove container %s: %w", name, err)
		}
	}
	
	log.Printf("Removed container: %s (ID: %s)", name, containerID)
	return nil
}

func (d *DockerRuntime) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
