# 🔄 Secret Rotation Guide

## 📋 Overview

Secret rotation is a critical security practice that involves periodically changing API keys, passwords, and other sensitive credentials. Sonic-Screwdriver provides built-in support for secret rotation with version history tracking.

## 🎯 Benefits of Secret Rotation

1. **Enhanced Security** - Regularly changing secrets reduces the window of opportunity for attackers
2. **Compliance** - Meets security best practices and regulatory requirements
3. **Audit Trail** - Maintains a history of secret changes for forensic analysis
4. **Minimal Downtime** - Rotation can be performed without interrupting services

## 🔧 Implementation

### **Rotation Mechanism**

Sonic-Screwdriver implements secret rotation with:

- **Current Value Storage** - The active secret value
- **History Tracking** - Previous values with timestamps and rotation actions
- **Atomic Updates** - Ensures secrets are fully updated or not at all

### **Storage Format**

Secrets are stored in an encrypted format with the following structure:

```
Main Secret: current_value
History Entry: previous_value|2026-04-22|rotated
```

## 🚀 Usage

### **CLI Commands**

#### **Rotate a Secret**

```bash
# Rotate a secret to a new value
sonic secret rotate openrouter_api_key --value "sk-new-key-here"

# Example output:
# Rotating secret openrouter_api_key...
# ✓ Secret openrouter_api_key rotated successfully
# 
# Rotation History:
#   1. rotated (on 2026-04-22)
```

#### **View Rotation History**

```bash
# Show rotation history for a secret
sonic secret history openrouter_api_key

# Example output:
# Showing rotation history for secret openrouter_api_key...
# 
# Rotation History for openrouter_api_key:
# ┌─────────┬──────────────────────────────────────────────┐
# │ #      │ Date                │ Action          │ Value (partial) │
# ├─────────┼──────────────────────────────────────────────┤
# │ 1      │ 2026-04-22           │ rotated         │ sk-old-key...   │
# └─────────┴──────────────────────────────────────────────┘
```

### **TUI Commands**

1. **Rotate Secret**
   - Navigate to: Secret Management → Rotate Secret
   - Enter secret name and new value
   - System shows rotation confirmation and history

2. **Show Rotation History**
   - Navigate to: Secret Management → Show Rotation History
   - Enter secret name
   - System displays formatted history table

## 📊 Best Practices

### **Rotation Frequency**

| Secret Type | Recommended Rotation | Notes |
|-------------|---------------------|-------|
| API Keys | Every 90 days | Critical production keys |
| Database Passwords | Every 180 days | Production databases |
| Service Accounts | Every 365 days | Internal service accounts |
| Development Keys | Every 30 days | Less critical, more frequent |

### **Rotation Process**

1. **Prepare** - Have new secret ready before rotation
2. **Rotate** - Use `sonic secret rotate` command
3. **Verify** - Test the new secret works
4. **Update** - Change the secret in all dependent services
5. **Monitor** - Watch for any issues post-rotation
6. **Document** - Record the rotation in your change log

## 🔐 Security Considerations

### **Before Rotation**

- **Backup** - Ensure you have a backup of current secrets
- **Test** - Verify the new secret works in a staging environment
- **Schedule** - Perform rotations during low-traffic periods

### **During Rotation**

- **Atomic** - Use atomic operations to prevent partial updates
- **Monitor** - Watch system logs for any authentication failures
- **Rollback** - Be prepared to revert if issues occur

### **After Rotation**

- **Verify** - Confirm all services are using the new secret
- **Audit** - Review rotation history for completeness
- **Document** - Update any documentation with new values

## 📋 Rotation Checklist

### **Pre-Rotation**

- [ ] Identify all secrets that need rotation
- [ ] Determine appropriate rotation schedule
- [ ] Generate new secret values securely
- [ ] Test new secrets in staging environment
- [ ] Notify team of upcoming rotation
- [ ] Schedule rotation during maintenance window

### **Rotation**

- [ ] Backup current secrets
- [ ] Execute rotation command
- [ ] Verify rotation was successful
- [ ] Update dependent services
- [ ] Monitor for authentication errors

### **Post-Rotation**

- [ ] Verify all services are operational
- [ ] Update documentation
- [ ] Record rotation in change log
- [ ] Schedule next rotation
- [ ] Review rotation process for improvements

## 🔄 Advanced Features

### **Automated Rotation**

While Sonic-Screwdriver currently supports manual rotation, you can automate the process using scripts:

```bash
#!/bin/bash
# Automated rotation script

# Generate new key (example)
NEW_KEY=$(openssl rand -hex 32)

# Rotate the secret
sonic secret rotate openrouter_api_key --value "sk-$NEW_KEY"

# Restart services that use this key
systemctl restart my-service

# Log the rotation
echo "$(date): Rotated openrouter_api_key" >> /var/log/secret_rotations.log
```

### **Bulk Rotation**

For rotating multiple secrets:

```bash
#!/bin/bash
# Bulk rotation script

SECRETS=("openrouter_api_key" "deepseek_api_key" "github_token")

for secret in "${SECRETS[@]}"; do
  echo "Rotating $secret..."
  # Generate and rotate each secret
  sonic secret rotate "$secret" --value "sk-$(openssl rand -hex 32)"
done
```

## 📚 Examples

### **Example 1: Rotating an API Key**

```bash
# Current key is expiring, rotate to new key
sonic secret rotate openrouter_api_key --value "sk-new-key-from-openrouter"

# Verify the rotation
sonic secret get openrouter_api_key

# Check history
sonic secret history openrouter_api_key
```

### **Example 2: Emergency Rotation**

```bash
# Key compromised, immediate rotation needed
sonic secret rotate github_token --value "ghp-new-token-from-github"

# Verify all services still work
sonic proxy health

# Check which services need the new token
sonic node list
```

### **Example 3: Scheduled Rotation**

```bash
# Monthly rotation script
#!/bin/bash

# Rotate all production keys
for secret in $(sonic secret list | grep prod); do
  sonic secret rotate "$secret" --value "$(generate_new_key)"
done

# Test all proxies
sonic proxy health

# Backup secrets
sonic secret backup /backups/secrets-$(date +%Y%m%d).json
```

## 🏥 Troubleshooting

### **Common Issues**

| Issue | Cause | Solution |
|-------|-------|----------|
| Rotation fails | Invalid new value | Check value format and length |
| History not showing | First rotation | History starts after first rotation |
| Services fail after rotation | Service not updated | Update service with new secret |
| Permission denied | Insufficient rights | Use appropriate admin credentials |

### **Recovery Procedures**

1. **Check History** - View previous values with `sonic secret history`
2. **Restore Backup** - Use `sonic secret restore` if needed
3. **Manual Override** - Temporarily use cached values with `--cached` flag
4. **Revert** - Rotate back to previous value if available

## 📈 Monitoring and Maintenance

### **Rotation Logs**

All rotations are logged in the secret store history. Regularly review:

```bash
# Check all secrets with history
for secret in $(sonic secret list); do
  echo "=== $secret ==="
  sonic secret history "$secret"
done
```

### **Backup Strategy**

1. **Before Rotation** - Always backup secrets
2. **After Rotation** - Verify backup includes new values
3. **Regular Backups** - Schedule automated backups
4. **Offsite Storage** - Store backups securely

```bash
# Backup before rotation
sonic secret backup /backups/secrets-pre-rotation.json

# Rotate secrets
sonic secret rotate my_secret --value "new_value"

# Backup after rotation
sonic secret backup /backups/secrets-post-rotation.json
```

## 🎯 Best Practices Summary

1. **Rotate Regularly** - Follow recommended schedules
2. **Test First** - Verify new secrets before full rotation
3. **Backup Always** - Never rotate without backup
4. **Monitor** - Watch for issues post-rotation
5. **Document** - Keep records of all rotations
6. **Automate** - Use scripts for repetitive rotations
7. **Review** - Regularly audit rotation history

## 📚 References

- **Sonic-Screwdriver Documentation** - Complete CLI reference
- **Security Best Practices** - Industry standards for secret management
- **Compliance Guidelines** - Regulatory requirements for your industry
- **Incident Response** - Procedures for compromised secrets

---

**Version:** 1.0.0  
**Last Updated:** 2026-04-22  
**Status:** Active  
**Owner:** Security Team