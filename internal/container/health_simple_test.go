package container

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHealthStatusStruct tests the HealthStatus struct
func TestHealthStatusStruct(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		status   HealthStatus
		expected string
	}{
		{
			name: "Healthy container",
			status: HealthStatus{
				Name:      "test-container",
				Status:    "healthy",
				Healthy:   true,
				Error:     "",
				Timestamp: now,
			},
			expected: "healthy",
		},
		{
			name: "Unhealthy container",
			status: HealthStatus{
				Name:      "test-container",
				Status:    "unhealthy",
				Healthy:   false,
				Error:     "container crashed",
				Timestamp: now,
			},
			expected: "unhealthy",
		},
		{
			name: "Not found container",
			status: HealthStatus{
				Name:      "test-container",
				Status:    "not_found",
				Healthy:   false,
				Error:     "Container not found",
				Timestamp: now,
			},
			expected: "not_found",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.Status)
			assert.Equal(t, "test-container", tt.status.Name)
			assert.Equal(t, tt.status.Healthy, tt.status.Healthy)
			assert.Equal(t, tt.status.Error, tt.status.Error)
			assert.WithinDuration(t, now, tt.status.Timestamp, 1*time.Second)
		})
	}
}

// TestHealthStatusEquality tests equality of HealthStatus structs
func TestHealthStatusEquality(t *testing.T) {
	now := time.Now()
	
	status1 := HealthStatus{
		Name:      "container1",
		Status:    "running",
		Healthy:   true,
		Error:     "",
		Timestamp: now,
	}
	
	status2 := HealthStatus{
		Name:      "container1",
		Status:    "running",
		Healthy:   true,
		Error:     "",
		Timestamp: now,
	}
	
	// Test that statuses with same values are equal
	assert.Equal(t, status1.Name, status2.Name)
	assert.Equal(t, status1.Status, status2.Status)
	assert.Equal(t, status1.Healthy, status2.Healthy)
	assert.Equal(t, status1.Error, status2.Error)
	assert.WithinDuration(t, status1.Timestamp, status2.Timestamp, 1*time.Second)
}

// TestHealthStatusModification tests modification of HealthStatus
func TestHealthStatusModification(t *testing.T) {
	now := time.Now()
	
	status := HealthStatus{
		Name:      "test-container",
		Status:    "running",
		Healthy:   true,
		Error:     "",
		Timestamp: now,
	}
	
	// Test that we can modify the status
	status.Healthy = false
	status.Error = "container failed"
	status.Status = "exited"
	
	assert.False(t, status.Healthy)
	assert.Equal(t, "container failed", status.Error)
	assert.Equal(t, "exited", status.Status)
}

// TestHealthStatusArray tests array of HealthStatus structs
func TestHealthStatusArray(t *testing.T) {
	now := time.Now()
	
	statuses := []HealthStatus{
		{
			Name:      "container1",
			Status:    "running",
			Healthy:   true,
			Error:     "",
			Timestamp: now,
		},
		{
			Name:      "container2",
			Status:    "exited",
			Healthy:   false,
			Error:     "container crashed",
			Timestamp: now,
		},
		{
			Name:      "container3",
			Status:    "running",
			Healthy:   true,
			Error:     "",
			Timestamp: now,
		},
	}
	
	assert.Len(t, statuses, 3)
	
	// Count healthy vs unhealthy
	healthyCount := 0
	for _, status := range statuses {
		if status.Healthy {
			healthyCount++
		}
	}
	
	assert.Equal(t, 2, healthyCount)
	assert.Equal(t, 1, len(statuses)-healthyCount)
}

// TestHealthStatusTime tests time handling in HealthStatus
func TestHealthStatusTime(t *testing.T) {
	now := time.Now()
	
	status := HealthStatus{
		Name:      "test-container",
		Status:    "running",
		Healthy:   true,
		Error:     "",
		Timestamp: now,
	}
	
	// Test that timestamp is set correctly
	assert.WithinDuration(t, now, status.Timestamp, 1*time.Second)
	
	// Test that we can update timestamp
	later := now.Add(1 * time.Minute)
	status.Timestamp = later
	assert.Equal(t, later, status.Timestamp)
}
