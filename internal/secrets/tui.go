package secrets

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TUI provides a simple text-based user interface
type TUI struct {
	secretStore *SecretStore
	nodeRegistry *NodeRegistry
	proxyServer *ProxyServer
}

// NewTUI creates a new TUI instance
func NewTUI(secretStore *SecretStore, nodeRegistry *NodeRegistry, proxyServer *ProxyServer) *TUI {
	return &TUI{
		secretStore:  secretStore,
		nodeRegistry: nodeRegistry,
		proxyServer:  proxyServer,
	}
}

// Run starts the TUI
func (t *TUI) Run() {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		t.showMainMenu()
		
		fmt.Print("\nEnter your choice (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "q" || input == "quit" {
			break
		}
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		switch num {
		case 1:
			t.showSecretMenu(reader)
		case 2:
			t.showNodeMenu(reader)
		case 3:
			t.showProxyMenu(reader)
		case 4:
			t.showStatus()
		case 5:
			t.showRemoteMenu(reader)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
	
	fmt.Println("Goodbye!")
}

// showMainMenu displays the main menu
func (t *TUI) showMainMenu() {
	fmt.Println("\n🎛️ Sonic-Screwdriver v2.0 - API Central Hub")
	fmt.Println("==============================================")
	fmt.Println("1. 🔑 Secret Management")
	fmt.Println("2. 🖥️ Node Management")
	fmt.Println("3. 🌐 API Proxy")
	fmt.Println("4. 📊 System Status")
	fmt.Println("5. 🖥️ Remote Access")
	fmt.Println("6. 🔜 Swarm Orchestration (planned)")
	fmt.Println("7. 🔜 Master Node Setup (planned)")
}

// showSecretMenu displays the secret management menu
func (t *TUI) showSecretMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n🔑 Secret Management")
		fmt.Println("=====================")
		fmt.Println("1. Add Secret")
		fmt.Println("2. Get Secret")
		fmt.Println("3. List Secrets")
		fmt.Println("4. Grant Access")
		fmt.Println("5. Revoke Access")
		fmt.Println("6. Show Policy")
		fmt.Println("7. Backup Secrets")
		fmt.Println("8. Restore Secrets")
		fmt.Println("9. Rotate Secret")
		fmt.Println("10. Show Rotation History")
		fmt.Println("11. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		switch num {
		case 1:
			t.addSecret(reader)
		case 2:
			t.getSecret(reader)
		case 3:
			t.listSecrets()
		case 4:
			t.grantAccess(reader)
		case 5:
			t.revokeAccess(reader)
		case 6:
			t.showPolicy(reader)
		case 7:
			t.backupSecrets(reader)
		case 8:
			t.restoreSecrets(reader)
		case 9:
			t.rotateSecret(reader)
		case 10:
			t.showSecretHistory(reader)
		case 11:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// showNodeMenu displays the node management menu
func (t *TUI) showNodeMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n🖥️ Node Management")
		fmt.Println("===================")
		fmt.Println("1. Register Node")
		fmt.Println("2. List Nodes")
		fmt.Println("3. Show Node Details")
		fmt.Println("4. Revoke Node")
		fmt.Println("5. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		switch num {
		case 1:
			t.registerNode(reader)
		case 2:
			t.listNodes()
		case 3:
			t.showNode(reader)
		case 4:
			t.revokeNode(reader)
		case 5:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// showProxyMenu displays the proxy management menu
func (t *TUI) showProxyMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n🌐 API Proxy")
		fmt.Println("============")
		fmt.Println("1. Show Proxy Status")
		fmt.Println("2. Check Proxy Health")
		fmt.Println("3. Make Proxy Call")
		fmt.Println("4. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		switch num {
		case 1:
			t.showProxyStatus()
		case 2:
			t.checkProxyHealth()
		case 3:
			t.makeProxyCall(reader)
		case 4:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// showStatus displays system status
func (t *TUI) showStatus() {
	fmt.Println("\n📊 System Status")
	fmt.Println("================")
	
	// Show secrets count
	secrets, _ := t.secretStore.ListSecrets()
	fmt.Printf("Secrets: %d\n", len(secrets))
	
	// Show nodes count
	nodes, _ := t.nodeRegistry.ListNodes()
	fmt.Printf("Registered Nodes: %d\n", len(nodes))
	
	// Show proxy status
	status := t.proxyServer.GetStatus()
	fmt.Println("\nProxy Status:")
	for provider, providerStatus := range status {
		calls := providerStatus["calls"].(int)
		rateLimit := providerStatus["rate_limit"].(string)
		fmt.Printf("  %s: %d calls, %s\n", provider, calls, rateLimit)
	}
	
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func (t *TUI) showRemoteMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n🖥️ Remote Access")
		fmt.Println("================")
		fmt.Println("1. Setup VNC Server")
		fmt.Println("2. Start VNC Server")
		fmt.Println("3. Stop VNC Server")
		fmt.Println("4. Setup SSH Access")
		fmt.Println("5. Setup Samba Sharing")
		fmt.Println("6. Show Remote Access Info")
		fmt.Println("7. Back to Main Menu")
		
		fmt.Print("\nEnter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		switch num {
		case 1:
			t.setupVNC(reader)
		case 2:
			t.startVNC()
		case 3:
			t.stopVNC()
		case 4:
			t.setupSSH()
		case 5:
			t.setupSamba(reader)
		case 6:
			t.showRemoteAccessInfo()
		case 7:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func (t *TUI) setupVNC(reader *bufio.Reader) {
	fmt.Print("Enter VNC password (default: password): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if password == "" {
		password = "password"
	}
	
	fmt.Print("Enter screen geometry (default: 1920x1080): ")
	geometry, _ := reader.ReadString('\n')
	geometry = strings.TrimSpace(geometry)
	if geometry == "" {
		geometry = "1920x1080"
	}
	
	fmt.Printf("Setting up VNC server with password and geometry %s...\n", geometry)
	
	// In a real implementation, you would create and store a VNC server instance
	// For this TUI, we'll just show the setup information
	fmt.Println("✓ VNC server setup completed")
	fmt.Println("Note: VNC server setup is simulated in this TUI")
	fmt.Println("In the actual CLI, run: sonic remote vnc setup <password> <geometry>")
}

func (t *TUI) startVNC() {
	fmt.Println("Starting VNC server...")
	// Simulated - in real implementation this would start the VNC server
	fmt.Println("✓ VNC server started")
	fmt.Println("Note: VNC server start is simulated in this TUI")
	fmt.Println("In the actual CLI, run: sonic remote vnc start")
}

func (t *TUI) stopVNC() {
	fmt.Println("Stopping VNC server...")
	// Simulated - in real implementation this would stop the VNC server
	fmt.Println("✓ VNC server stopped")
	fmt.Println("Note: VNC server stop is simulated in this TUI")
	fmt.Println("In the actual CLI, run: sonic remote vnc stop")
}

func (t *TUI) setupSSH() {
	fmt.Println("Setting up SSH for remote access...")
	// Simulated - in real implementation this would setup SSH
	fmt.Println("✓ SSH setup completed")
	fmt.Println("Note: SSH setup is simulated in this TUI")
	fmt.Println("In the actual CLI, run: sonic remote ssh setup")
}

func (t *TUI) setupSamba(reader *bufio.Reader) {
	fmt.Print("Enter share name: ")
	shareName, _ := reader.ReadString('\n')
	shareName = strings.TrimSpace(shareName)
	
	fmt.Print("Enter share path: ")
	sharePath, _ := reader.ReadString('\n')
	sharePath = strings.TrimSpace(sharePath)
	
	fmt.Printf("Setting up Samba file sharing for %s at %s...\n", shareName, sharePath)
	// Simulated - in real implementation this would setup Samba
	fmt.Println("✓ Samba setup completed")
	fmt.Println("Note: Samba setup is simulated in this TUI")
	fmt.Printf("In the actual CLI, run: sonic remote samba setup %s %s\n", shareName, sharePath)
}

func (t *TUI) showRemoteAccessInfo() {
	fmt.Println("\n🌐 Remote Access Information")
	fmt.Println("============================")
	fmt.Println("This is a simulated remote access info screen.")
	fmt.Println("In the actual CLI, run: sonic remote info")
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// Helper methods for secret management
func (t *TUI) addSecret(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Print("Enter secret value: ")
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)
	
	if err := t.secretStore.AddSecret(name, value); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Secret added successfully")
	}
	
	// Test if it's an API key
	if name == "openrouter_api_key" || name == "deepseek_api_key" || name == "gemini_api_key" || name == "github_token" {
		provider := ""
		if name == "openrouter_api_key" {
			provider = "openrouter"
		} else if name == "deepseek_api_key" {
			provider = "deepseek"
		} else if name == "gemini_api_key" {
			provider = "gemini"
		} else if name == "github_token" {
			provider = "github"
		}
		
		if provider != "" {
			fmt.Printf("Testing %s API key...\n", provider)
			err := t.proxyServer.TestAPIKey(provider, value)
			if err != nil {
				fmt.Printf("⚠️  API key test failed: %v\n", err)
			} else {
				fmt.Printf("✅ API key is valid\n")
			}
		}
	}
}

func (t *TUI) getSecret(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	value, err := t.secretStore.GetSecret(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Secret value: %s\n", value)
	}
}

func (t *TUI) listSecrets() {
	secrets, err := t.secretStore.ListSecrets()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return
	}
	
	fmt.Println("\nSecrets:")
	for i, secret := range secrets {
		policy, _ := t.secretStore.GetPolicy(secret)
		allowedNodes := "all"
		if len(policy.AllowedNodes) > 0 {
			allowedNodes = strings.Join(policy.AllowedNodes, ", ")
		}
		fmt.Printf("%d. %s (allowed: %s)\n", i+1, secret, allowedNodes)
	}
}

func (t *TUI) grantAccess(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	secretName, _ := reader.ReadString('\n')
	secretName = strings.TrimSpace(secretName)
	
	fmt.Print("Enter node name: ")
	nodeName, _ := reader.ReadString('\n')
	nodeName = strings.TrimSpace(nodeName)
	
	// Get current policy
	policy, err := t.secretStore.GetPolicy(secretName)
	if err != nil {
		policy = SecretPolicy{
			AllowedNodes: []string{},
			AllowedRoles:  []string{},
			RateLimit:     "60/min",
		}
	}
	
	// Add node to allowed nodes
	for _, allowedNode := range policy.AllowedNodes {
		if allowedNode == nodeName {
			fmt.Println("✓ Access already granted")
			return
		}
	}
	
	policy.AllowedNodes = append(policy.AllowedNodes, nodeName)
	
	if err := t.secretStore.SetPolicy(secretName, policy); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Access granted successfully")
	}
	
	// Also update node registry
	if err := t.nodeRegistry.GrantSecretAccess(nodeName, secretName); err != nil {
		fmt.Printf("Warning: Failed to update node registry: %v\n", err)
	}
}

func (t *TUI) revokeAccess(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	secretName, _ := reader.ReadString('\n')
	secretName = strings.TrimSpace(secretName)
	
	fmt.Print("Enter node name: ")
	nodeName, _ := reader.ReadString('\n')
	nodeName = strings.TrimSpace(nodeName)
	
	policy, err := t.secretStore.GetPolicy(secretName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Remove node from allowed nodes
	newAllowedNodes := []string{}
	for _, allowedNode := range policy.AllowedNodes {
		if allowedNode != nodeName {
			newAllowedNodes = append(newAllowedNodes, allowedNode)
		}
	}
	
	policy.AllowedNodes = newAllowedNodes
	
	if err := t.secretStore.SetPolicy(secretName, policy); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Access revoked successfully")
	}
	
	// Also update node registry
	if err := t.nodeRegistry.RevokeSecretAccess(nodeName, secretName); err != nil {
		fmt.Printf("Warning: Failed to update node registry: %v\n", err)
	}
}

func (t *TUI) showPolicy(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	policy, err := t.secretStore.GetPolicy(name)
	if err != nil {
		fmt.Println("Policy: No policy set")
		return
	}
	
	fmt.Printf("\nPolicy for %s:\n", name)
	fmt.Printf("  Allowed Nodes: %v\n", policy.AllowedNodes)
	fmt.Printf("  Allowed Roles: %v\n", policy.AllowedRoles)
	fmt.Printf("  Rate Limit: %s\n", policy.RateLimit)
}

func (t *TUI) backupSecrets(reader *bufio.Reader) {
	fmt.Print("Enter backup file path: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	
	if err := t.secretStore.Backup(filePath); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Backup created successfully")
	}
}

func (t *TUI) restoreSecrets(reader *bufio.Reader) {
	fmt.Print("Enter backup file path: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	
	if err := t.secretStore.Restore(filePath); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Backup restored successfully")
	}
}

func (t *TUI) rotateSecret(reader *bufio.Reader) {
	fmt.Print("Enter secret name to rotate: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Print("Enter new secret value: ")
	newValue, _ := reader.ReadString('\n')
	newValue = strings.TrimSpace(newValue)
	
	fmt.Printf("Rotating secret %s...\n", name)
	
	if err := t.secretStore.RotateSecret(name, newValue); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Secret %s rotated successfully\n", name)
	
	// Show history
	history, err := t.secretStore.GetSecretHistory(name)
	if err == nil && len(history) > 0 {
		fmt.Println("\nRotation History:")
		for i, entry := range history {
			fmt.Printf("  %d. %s (on %s)\n", i+1, entry["action"], entry["date"])
		}
	}
}

func (t *TUI) showSecretHistory(reader *bufio.Reader) {
	fmt.Print("Enter secret name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	history, err := t.secretStore.GetSecretHistory(name)
	if err != nil {
		if err.Error() == "no history available" {
			fmt.Println("No rotation history available for this secret")
			return
		}
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	if len(history) == 0 {
		fmt.Println("No rotation history available")
		return
	}
	
	fmt.Printf("\nRotation History for %s:\n", name)
	fmt.Println("┌─────────┬──────────────────────────────────────────────┐")
	fmt.Println("│ #      │ Date                │ Action          │ Value (partial) │")
	fmt.Println("├─────────┼──────────────────────────────────────────────┤")
	
	for i, entry := range history {
		valuePreview := entry["value"]
		if len(valuePreview) > 10 {
			valuePreview = valuePreview[:10] + "..."
		}
		fmt.Printf("│ %-7d │ %-20s │ %-14s │ %-15s │\n", 
			i+1, entry["date"], entry["action"], valuePreview)
	}
	
	fmt.Println("└─────────┴──────────────────────────────────────────────┘")
	
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// Helper methods for node management
func (t *TUI) registerNode(reader *bufio.Reader) {
	fmt.Print("Enter node name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Print("Enter master address: ")
	masterAddr, _ := reader.ReadString('\n')
	masterAddr = strings.TrimSpace(masterAddr)
	
	node, err := t.nodeRegistry.RegisterNode(name, masterAddr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Node %s registered (ID: %s)\n", node.Name, node.ID)
	}
}

func (t *TUI) listNodes() {
	nodes, err := t.nodeRegistry.ListNodes()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	if len(nodes) == 0 {
		fmt.Println("No nodes registered")
		return
	}
	
	fmt.Println("\nRegistered Nodes:")
	fmt.Printf("%-15s %-10s %-20s\n", "ID", "NAME", "STATUS")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")
	
	for _, node := range nodes {
		fmt.Printf("%-15s %-10s %-20s\n", node.ID, node.Name, node.Status)
	}
}

func (t *TUI) showNode(reader *bufio.Reader) {
	fmt.Print("Enter node name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	node, err := t.nodeRegistry.GetNode(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("\nNode: %s\n", node.Name)
	fmt.Printf("  ID: %s\n", node.ID)
	fmt.Printf("  Status: %s\n", node.Status)
	fmt.Printf("  Last seen: %s\n", node.LastSeen.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Allowed secrets: %v\n", node.AllowedSecrets)
}

func (t *TUI) revokeNode(reader *bufio.Reader) {
	fmt.Print("Enter node name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	if err := t.nodeRegistry.RevokeNode(name); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Node %s revoked\n", name)
	}
}

// Helper methods for proxy management
func (t *TUI) showProxyStatus() {
	status := t.proxyServer.GetStatus()
	
	fmt.Println("\nProxy Status:")
	fmt.Println("┌─────────────┬──────────┬────────────┬─────────────┬──────────────┐")
	fmt.Println("│ Provider    │ Calls    │ Errors     │ Rate Limit  │ Status       │")
	fmt.Println("├─────────────┼──────────┼────────────┼─────────────┼──────────────┤")
	
	providers := []string{"openrouter", "deepseek", "gemini", "github"}
	for _, provider := range providers {
		if providerStatus, ok := status[provider]; ok {
			calls := providerStatus["calls"].(int)
			rateLimit := providerStatus["rate_limit"].(string)
			statusText := providerStatus["status"].(string)
			
			emoji := "🟢"
			if statusText != "healthy" {
				emoji = "🔴"
			}
			
			fmt.Printf("│ %-10s │ %-7d │ 0 (0%%)     │ %-10s │ %s %-11s │\n", 
				provider, calls, rateLimit, emoji, statusText)
		}
	}
	
	fmt.Println("└─────────────┴──────────┴────────────┴─────────────┴──────────────┘")
}

func (t *TUI) checkProxyHealth() {
	health := t.proxyServer.GetHealth()
	
	fmt.Println("\nProxy Health:")
	for provider, status := range health {
		emoji := "✅"
		if status != "healthy" {
			emoji = "❌"
		}
		fmt.Printf("%s %s: %s\n", emoji, provider, status)
	}
}

func (t *TUI) makeProxyCall(reader *bufio.Reader) {
	fmt.Print("Enter provider (openrouter, deepseek, gemini, github): ")
	provider, _ := reader.ReadString('\n')
	provider = strings.TrimSpace(provider)
	
	fmt.Print("Enter request data (JSON): ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)
	
	var requestData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &requestData); err != nil {
		fmt.Printf("Error: Failed to parse JSON: %v\n", err)
		return
	}
	
	request := ProxyRequest{
		Provider: provider,
		Method:   "POST",
		Path:     "/chat/completions",
		Headers:  map[string]string{},
		Body:     requestData,
	}
	
	response, err := t.proxyServer.HandleRequest(request)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("\n✓ Proxy call completed\n")
	fmt.Printf("Status: %d\n", response.Status)
	fmt.Printf("Response: %+v\n", response.Body)
}