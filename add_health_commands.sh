#!/bin/bash

# Add health and repair commands to sonic CLI
# This script inserts the new commands after the "logs" case

FILE="/home/wizard/code-vault/sonic-screwdriver/cmd/sonic/main.go"

# Create a temporary file
TMPFILE=$(mktemp)

# Find the line number where we want to insert (after "logs" case)
LINE_NUMBER=$(grep -n 'case "config":' "$FILE" | cut -d: -f1)

# Copy everything before the config case to temp file
head -n $((LINE_NUMBER - 1)) "$FILE" > "$TMPFILE"

# Add the new health and repair commands
cat >> "$TMPFILE" << 'EOF'
		case "health":
			if len(os.Args) < 3 {
				fmt.Println("Usage: sonic health <game>|--all")
				os.Exit(1)
			}
			if os.Args[2] == "--all" {
				healthStatuses, err := runtime.GetAllContainerHealth()
				if err != nil {
					log.Fatalf("Failed to get container health: %v", err)
				}
				if len(healthStatuses) == 0 {
					fmt.Println("No Sonic containers found")
					return
				}
				fmt.Println("Container Health Status:")
				for _, status := range healthStatuses {
					healthIcon := "✅"
					if !status.Healthy {
						healthIcon = "❌"
					}
					fmt.Printf("  %s %s: %s", healthIcon, status.Name, status.Status)
					if status.Error != "" {
						fmt.Printf(" (%s)", status.Error)
					}
					fmt.Println()
				}
			} else {
				gameName := os.Args[2]
				status, err := runtime.CheckContainerHealth(gameName)
				if err != nil {
					log.Fatalf("Failed to check health: %v", err)
				}
				healthIcon := "✅"
				if !status.Healthy {
					healthIcon = "❌"
				}
				fmt.Printf("%s Health Status:\n", gameName)
				fmt.Printf("  Status: %s\n", status.Status)
				fmt.Printf("  Healthy: %t\n", status.Healthy)
				if status.Error != "" {
					fmt.Printf("  Error: %s\n", status.Error)
				}
				fmt.Printf("  Last Check: %s\n", status.Timestamp.Format("2006-01-02 15:04:05"))
			}
		case "repair":
			if len(os.Args) < 3 {
				fmt.Println("Usage: sonic repair <game>|--all")
				os.Exit(1)
			}
			if os.Args[2] == "--all" {
				healthStatuses, err := runtime.GetAllContainerHealth()
				if err != nil {
					log.Fatalf("Failed to get container health: %v", err)
				}
				repairedCount := 0
				for _, status := range healthStatuses {
					if !status.Healthy {
						fmt.Printf("Attempting to repair %s...\n", status.Name)
						err := runtime.RestartContainer(status.Name)
						if err != nil {
							log.Printf("Failed to repair %s: %v", status.Name, err)
						} else {
							fmt.Printf("✅ Repaired %s\n", status.Name)
							repairedCount++
						}
					}
				}
				if repairedCount == 0 {
					fmt.Println("All containers are healthy, no repairs needed")
				} else {
					fmt.Printf("Repaired %d container(s)\n", repairedCount)
				}
			} else {
				gameName := os.Args[2]
				fmt.Printf("Attempting to repair %s...\n", gameName)
				status, err := runtime.CheckContainerHealth(gameName)
				if err != nil {
					log.Fatalf("Failed to check health: %v", err)
				}
				if status.Healthy {
					fmt.Printf("%s is already healthy\n", gameName)
					return
				}
				err = runtime.RestartContainer(gameName)
				if err != nil {
					log.Fatalf("Failed to repair: %v", err)
				}
				fmt.Printf("✅ Repaired %s\n", gameName)
			}

EOF

# Append the rest of the file (from config case onwards)
tail -n +$LINE_NUMBER "$FILE" >> "$TMPFILE"

# Update help text to include new commands
sed -i 's/sonic logs <game>        Show container logs/sonic logs <game>        Show container logs\n  sonic health <game>|--all  Check container health status\n  sonic repair <game>|--all  Repair unhealthy containers/' "$TMPFILE"

# Replace original file
mv "$TMPFILE" "$FILE"

echo "✅ Added health and repair commands to sonic CLI"
