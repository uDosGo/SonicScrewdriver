package remote

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// VNCServer manages VNC server functionality
type VNCServer struct {
	port        int
	password    string
	geometry    string
	display     int
	isRunning   bool
}

// NewVNCServer creates a new VNC server instance
func NewVNCServer(port int, password, geometry string) *VNCServer {
	return &VNCServer{
		port:     port,
		password: password,
		geometry: geometry,
		display:  1, // Default to display :1
	}
}

// SetupVNC sets up VNC server on the system
func (v *VNCServer) SetupVNC() error {
	// Check if required packages are installed
	if err := v.checkDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %v", err)
	}

	// Set VNC password
	if err := v.setVNCPassword(); err != nil {
		return fmt.Errorf("failed to set VNC password: %v", err)
	}

	// Configure VNC server
	if err := v.configureVNC(); err != nil {
		return fmt.Errorf("failed to configure VNC: %v", err)
	}

	return nil
}

// StartVNC starts the VNC server
func (v *VNCServer) StartVNC() error {
	if v.isRunning {
		return fmt.Errorf("VNC server is already running")
	}

	// Start VNC server
	cmd := exec.Command("vncserver", fmt.Sprintf(":%d", v.display), "-geometry", v.geometry, "-localhost", "no")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VNC server: %v", err)
	}

	v.isRunning = true
	log.Printf("VNC server started on display :%d (port %d)", v.display, v.port)
	log.Printf("Connect using: vncviewer %s:%d", GetLocalIP(), v.display)

	return nil
}

// StopVNC stops the VNC server
func (v *VNCServer) StopVNC() error {
	if !v.isRunning {
		return fmt.Errorf("VNC server is not running")
	}

	// Stop VNC server
	cmd := exec.Command("vncserver", "-kill", fmt.Sprintf(":%d", v.display))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop VNC server: %v", err)
	}

	v.isRunning = false
	log.Printf("VNC server stopped")

	return nil
}

// checkDependencies checks if required packages are installed
func (v *VNCServer) checkDependencies() error {
	packages := []string{"tightvncserver", "x11vnc", "vnc4server"}
	
	for _, pkg := range packages {
		if _, err := exec.LookPath(pkg); err == nil {
			log.Printf("Found VNC package: %s", pkg)
			return nil
		}
	}

	return fmt.Errorf("no VNC server package found. Please install one of: %s", strings.Join(packages, ", "))
}

// setVNCPassword sets the VNC password
func (v *VNCServer) setVNCPassword() error {
	if v.password == "" {
		return fmt.Errorf("VNC password cannot be empty")
	}

	// Create a temporary password file
	tempFile := "/tmp/vncpwd.txt"
	if err := os.WriteFile(tempFile, []byte(v.password+"\n"+v.password+"\n"), 0600); err != nil {
		return fmt.Errorf("failed to create password file: %v", err)
	}
	defer os.Remove(tempFile)

	// Set password using vncpasswd
	cmd := exec.Command("vncpasswd", tempFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set VNC password: %v", err)
	}

	return nil
}

// configureVNC configures the VNC server
func (v *VNCServer) configureVNC() error {
	// Create VNC config directory
	configDir := fmt.Sprintf("/home/%s/.vnc", GetCurrentUser())
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Create xstartup file
	xstartupContent := `#!/bin/sh
unset SESSION_MANAGER
unset DBUS_SESSION_BUS_ADDRESS

[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources

xsetroot -solid grey
vncconfig -iconic &

# Start your window manager or desktop environment
exec startxfce4 &
`

	xstartupPath := fmt.Sprintf("%s/xstartup", configDir)
	if err := os.WriteFile(xstartupPath, []byte(xstartupContent), 0755); err != nil {
		return fmt.Errorf("failed to create xstartup file: %v", err)
	}

	return nil
}

// GetLocalIP gets the local IP address
func GetLocalIP() string {
	// Simple implementation - in real code you'd use net package
	return "localhost" // Placeholder
}

// GetCurrentUser gets the current username
func GetCurrentUser() string {
	return os.Getenv("USER")
}

// SetupSSH sets up SSH for remote access
func SetupSSH() error {
	log.Println("Setting up SSH for remote access...")
	
	// Check if SSH is installed
	if _, err := exec.LookPath("ssh"); err != nil {
		return fmt.Errorf("SSH is not installed")
	}

	// Enable SSH service
	cmd := exec.Command("sudo", "systemctl", "enable", "ssh")
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to enable SSH service: %v", err)
	}

	// Start SSH service
	cmd = exec.Command("sudo", "systemctl", "start", "ssh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start SSH service: %v", err)
	}

	log.Println("SSH service is running")
	log.Printf("Connect using: ssh %s@%s", GetCurrentUser(), GetLocalIP())

	return nil
}

// SetupSamba sets up Samba for file sharing
func SetupSamba(shareName, sharePath string) error {
	log.Println("Setting up Samba for file sharing...")
	
	// Check if Samba is installed
	if _, err := exec.LookPath("smbd"); err != nil {
		return fmt.Errorf("Samba is not installed")
	}

	// Create Samba configuration
	config := fmt.Sprintf(`[%s]
	path = %s
	writable = yes
	browseable = yes
	guest ok = yes
	create mask = 0777
	directory mask = 0777
`, shareName, sharePath)

	// Append to smb.conf
	configFile := "/etc/samba/smb.conf"
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open Samba config: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(config); err != nil {
		return fmt.Errorf("failed to write Samba config: %v", err)
	}

	// Restart Samba service
	cmd := exec.Command("sudo", "systemctl", "restart", "smbd")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart Samba: %v", err)
	}

	log.Println("Samba file sharing is set up")
	log.Printf("Access shared files at: smb://%s/%s", GetLocalIP(), shareName)

	return nil
}

// GetRemoteAccessInfo returns information about remote access
func GetRemoteAccessInfo() string {
	user := GetCurrentUser()
	ip := GetLocalIP()
	
	info := "🌐 Remote Access Information:\n"
	info += fmt.Sprintf("  SSH: ssh %s@%s\n", user, ip)
	info += fmt.Sprintf("  VNC: vncviewer %s:1\n", ip)
	info += fmt.Sprintf("  Samba: smb://%s/shared\n", ip)
	info += "\n📝 Note: Make sure to configure your firewall to allow these connections."
	
	return info
}