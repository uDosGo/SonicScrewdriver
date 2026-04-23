package secrets

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Node represents a registered node
type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"last_seen"`
	AllowedSecrets []string `json:"allowed_secrets"`
}

// NodeRegistry manages registered nodes
type NodeRegistry struct {
	nodes map[string]Node
	filePath string
	mu sync.RWMutex
}

// NewNodeRegistry creates a new node registry
func NewNodeRegistry(filePath string) (*NodeRegistry, error) {
	registry := &NodeRegistry{
		nodes: make(map[string]Node),
		filePath: filePath,
	}

	// Load existing nodes if file exists
	if _, err := os.Stat(filePath); err == nil {
		if err := registry.load(); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	return registry, nil
}

// load loads nodes from file
func (r *NodeRegistry) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var nodes map[string]Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = nodes

	return nil
}

// save saves nodes to file
func (r *NodeRegistry) save() error {
	r.mu.RLock()
	data, err := json.Marshal(r.nodes)
	r.mu.RUnlock()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

// RegisterNode registers a new node
func (r *NodeRegistry) RegisterNode(name, masterAddr string) (Node, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if node already exists
	for _, node := range r.nodes {
		if node.Name == name {
			return node, errors.New("node already registered")
		}
	}

	// Generate node ID (simple implementation)
	nodeID := "node_" + generateRandomID(8)

	node := Node{
		ID:        nodeID,
		Name:      name,
		Status:    "online",
		LastSeen:  time.Now(),
		AllowedSecrets: []string{},
	}

	r.nodes[nodeID] = node
	if err := r.save(); err != nil {
		return Node{}, err
	}

	return node, nil
}

// GetNode retrieves a node by name
func (r *NodeRegistry) GetNode(name string) (Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, node := range r.nodes {
		if node.Name == name {
			return node, nil
		}
	}

	return Node{}, errors.New("node not found")
}

// ListNodes returns all registered nodes
func (r *NodeRegistry) ListNodes() ([]Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nodes := make([]Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// RevokeNode revokes a node's access
func (r *NodeRegistry) RevokeNode(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, node := range r.nodes {
		if node.Name == name {
			delete(r.nodes, id)
			return r.save()
		}
	}

	return errors.New("node not found")
}

// GrantSecretAccess grants a node access to a secret
func (r *NodeRegistry) GrantSecretAccess(nodeName, secretName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, node := range r.nodes {
		if node.Name == nodeName {
			// Check if already granted
			for _, secret := range node.AllowedSecrets {
				if secret == secretName {
					return nil // Already granted
				}
			}
			node.AllowedSecrets = append(node.AllowedSecrets, secretName)
			r.nodes[id] = node
			return r.save()
		}
	}

	return errors.New("node not found")
}

// RevokeSecretAccess revokes a node's access to a secret
func (r *NodeRegistry) RevokeSecretAccess(nodeName, secretName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, node := range r.nodes {
		if node.Name == nodeName {
			newSecrets := []string{}
			for _, secret := range node.AllowedSecrets {
				if secret != secretName {
					newSecrets = append(newSecrets, secret)
				}
			}
			node.AllowedSecrets = newSecrets
			r.nodes[id] = node
			return r.save()
		}
	}

	return errors.New("node not found")
}

// generateRandomID generates a simple random ID (placeholder implementation)
func generateRandomID(length int) string {
	// TODO: Implement proper random ID generation
	return "abc123" // Placeholder
}