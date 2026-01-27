# Notifox CLI Usage Guide

This document provides comprehensive documentation for using the Notifox CLI tool.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Command Reference](#command-reference)
4. [Message Input Methods](#message-input-methods)
5. [Environment Variables](#environment-variables)
6. [Examples](#examples)
7. [Error Handling](#error-handling)
8. [Integration Examples](#integration-examples)
9. [Exit Codes](#exit-codes)
10. [Troubleshooting](#troubleshooting)

---

## Installation

### Method 1: Download Pre-built Binary

The easiest way to install the CLI is to download a pre-built binary from the GitHub releases page.

#### macOS (Apple Silicon / ARM64)

```bash
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_darwin_arm64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/
```

#### macOS (Intel / AMD64)

```bash
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_darwin_amd64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/
```

#### Linux (AMD64)

```bash
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_linux_amd64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/
```

#### Linux (ARM64)

```bash
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_linux_arm64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/
```

#### Windows

1. Visit the [releases page](https://github.com/notifoxhq/notifox-cli/releases)
2. Download `notifox-cli_windows_amd64.zip`
3. Extract the zip file
4. Move `notifox.exe` to a directory in your PATH (e.g., `C:\Windows\System32`)

### Method 2: Build from Source

If you have Go installed (version 1.23 or later), you can build from source:

```bash
git clone https://github.com/notifoxhq/notifox-cli.git
cd notifox-cli
go build -o notifox .
sudo mv notifox /usr/local/bin/  # On Unix systems
```

### Verify Installation

After installation, verify the CLI is working:

```bash
notifox send -h
```

This should display the help text for the `send` command.

---

## Configuration

### Required: API Key

The CLI requires a Notifox API key to authenticate requests. Set it as an environment variable:

```bash
export NOTIFOX_API_KEY="your_api_key_here"
```

**Important:** The API key is required for all operations. If it's not set, the CLI will exit with an error.

### Optional: Default Audience and Channel

You can set default values for audience and channel to avoid typing them repeatedly:

```bash
export NOTIFOX_AUDIENCE="oncall"
export NOTIFOX_CHANNEL="sms"
```

When these are set, you can omit the `-a` and `-c` flags:

```bash
echo "Server is down" | notifox send
```

**Flag Precedence:** Command-line flags always override environment variables. This follows the same behavior as AWS CLI.

### Optional: Custom Base URL

For development or testing against a different Notifox instance:

```bash
export NOTIFOX_BASE_URL="https://api-staging.notifox.com"
```

This is useful for:
- Testing against staging environments
- Using a self-hosted Notifox instance
- Development and debugging

---

## Command Reference

### `notifox send`

Sends an alert to a specified audience via SMS or email.

#### Syntax

```bash
notifox send [flags]
```

#### Flags

| Flag | Short | Required | Description | Example |
|------|-------|-----------|-------------|----------|
| `--audience` | `-a` | Yes* | The audience to send the alert to | `-a oncall` |
| `--channel` | `-c` | Yes* | Channel to use: `sms` or `email` | `-c sms` |
| `--message` | `-m` | No** | Message body to send | `-m "Server down"` |

\* Required if not set via environment variables (`NOTIFOX_AUDIENCE`, `NOTIFOX_CHANNEL`)  
\** Required if no input is provided via stdin

#### Flag Details

**`-a, --audience`**

- Specifies which audience should receive the alert
- Audience names are configured in your Notifox account
- Common examples: `oncall`, `platform`, `ops`, `mathis`
- Can be set via `NOTIFOX_AUDIENCE` environment variable
- Flag value overrides environment variable

**`-c, --channel`**

- Specifies the delivery channel
- Valid values: `sms` or `email`
- Case-sensitive
- Can be set via `NOTIFOX_CHANNEL` environment variable
- Flag value overrides environment variable

**`-m, --message`**

- Provides the message body directly
- If not provided, the CLI reads from stdin
- If both `-m` flag and stdin are provided, the flag takes precedence
- Message content is sent as-is (no automatic formatting)

#### Examples

Basic usage with flags:

```bash
notifox send -a oncall -c sms -m "Database connection failed"
```

Using environment variables:

```bash
export NOTIFOX_AUDIENCE="oncall"
export NOTIFOX_CHANNEL="sms"
notifox send -m "Database connection failed"
```

Piping input:

```bash
echo "Database connection failed" | notifox send -a oncall -c sms
```

---

## Message Input Methods

The CLI supports two methods for providing the message content:

### Method 1: Command-line Flag (`-m` or `--message`)

Use the `-m` flag to provide the message directly:

```bash
notifox send -a oncall -c sms -m "Server is down"
```

**Advantages:**
- Simple one-liner commands
- No need for pipes
- Works well in scripts

### Method 2: Standard Input (stdin)

Pipe content into the CLI:

```bash
echo "Server is down" | notifox send -a oncall -c sms
```

Or pipe output from other commands:

```bash
kubectl get pods | notifox send -a platform -c email
tail -n 50 /var/log/app.log | notifox send -a ops -c sms
```

**Advantages:**
- Works naturally with Unix pipelines
- Can send large amounts of text
- Integrates with existing command-line workflows

### Precedence Rules

1. If `-m` or `--message` flag is provided, it is used (highest priority)
2. If no flag is provided, stdin is read (if available)
3. If neither is provided, the command fails with an error

**Important Notes:**

- Stdin is only read if it's a pipe (not a terminal)
- If you run `notifox send` without piping and without `-m`, it will wait for input or fail
- Empty messages are not allowed and will cause an error

---

## Environment Variables

### `NOTIFOX_API_KEY` (Required)

Your Notifox API key for authentication.

```bash
export NOTIFOX_API_KEY="sk_live_abc123xyz"
```

**Security Note:** Never commit API keys to version control. Use environment variables or secret management tools.

### `NOTIFOX_AUDIENCE` (Optional)

Default audience name. Can be overridden with `-a` or `--audience` flag.

```bash
export NOTIFOX_AUDIENCE="oncall"
```

### `NOTIFOX_CHANNEL` (Optional)

Default channel (`sms` or `email`). Can be overridden with `-c` or `--channel` flag.

```bash
export NOTIFOX_CHANNEL="sms"
```

### `NOTIFOX_BASE_URL` (Optional)

Override the default Notifox API base URL. Useful for staging environments or self-hosted instances.

```bash
export NOTIFOX_BASE_URL="https://api-staging.notifox.com"
```

**Default:** Uses the production Notifox API URL.

---

## Examples

### Basic Alert

Send a simple alert:

```bash
notifox send -a oncall -c sms -m "Production deployment completed"
```

### Piping Command Output

Send the output of a command:

```bash
df -h | notifox send -a ops -c email
```

### Sending Log Files

Send the last 100 lines of a log file:

```bash
tail -n 100 /var/log/nginx/error.log | notifox send -a platform -c email
```

### Using Environment Variables

Set defaults and send without flags:

```bash
export NOTIFOX_AUDIENCE="oncall"
export NOTIFOX_CHANNEL="sms"
echo "Server restart required" | notifox send
```

### Overriding Environment Variables

Flags override environment variables:

```bash
export NOTIFOX_AUDIENCE="oncall"
export NOTIFOX_CHANNEL="sms"
# This sends to "mathis" via "email", not "oncall" via "sms"
notifox send -a mathis -c email -m "Urgent: Check database"
```

### Multi-line Messages

Send multi-line content via stdin:

```bash
cat <<EOF | notifox send -a ops -c email
System Status Report:
- CPU: 85%
- Memory: 70%
- Disk: 60%
EOF
```

### Conditional Alerts

Use in shell scripts with conditional logic:

```bash
#!/bin/bash
if ! systemctl is-active --quiet nginx; then
    echo "Nginx is down!" | notifox send -a ops -c sms
fi
```

---

## Error Handling

The CLI provides clear error messages for common issues:

### Missing API Key

```
Error: NOTIFOX_API_KEY environment variable is required
```

**Solution:** Set the `NOTIFOX_API_KEY` environment variable.

### Missing Audience

```
Error: audience is required (use -a/--audience or set NOTIFOX_AUDIENCE)
```

**Solution:** Provide audience via `-a` flag or `NOTIFOX_AUDIENCE` environment variable.

### Missing Channel

```
Error: channel is required (use -c/--channel or set NOTIFOX_CHANNEL)
```

**Solution:** Provide channel via `-c` flag or `NOTIFOX_CHANNEL` environment variable.

### Missing Message

```
Error: message is required (provide via -m/--message or stdin)
```

**Solution:** Provide message via `-m` flag or pipe input to stdin.

### Invalid Channel

```
Error: channel must be 'sms' or 'email', got 'text'
```

**Solution:** Use `sms` or `email` (case-sensitive).

### Authentication Error

```
Error: authentication failed: Invalid API key (status: 401)
```

**Solution:** Check that your API key is correct and active.

### Insufficient Balance

```
Error: insufficient balance: Your account balance is too low
```

**Solution:** Add credits to your Notifox account.

### Rate Limit Exceeded

```
Error: rate limit exceeded: Too many requests
```

**Solution:** Wait before sending more alerts, or upgrade your plan.

### Connection Error

```
Error: connection error: dial tcp: lookup api.notifox.com: no such host
```

**Solution:** Check your internet connection and DNS settings.

### API Error

```
Error: API error: Invalid audience name (status: 400)
```

**Solution:** Verify the audience name exists in your Notifox account.

---

## Integration Examples

### Cron Jobs

Send alerts from scheduled cron jobs:

```bash
# /etc/cron.daily/backup-check
#!/bin/bash
if ! /usr/local/bin/backup.sh; then
    echo "Backup failed at $(date)" | notifox send -a ops -c sms
fi
```

### GitHub Actions

Use in CI/CD pipelines:

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy
        run: ./deploy.sh
      - name: Notify on failure
        if: failure()
        env:
          NOTIFOX_API_KEY: ${{ secrets.NOTIFOX_API_KEY }}
        run: |
          echo "Deployment failed" | notifox send -a platform -c email
```

### Kubernetes Jobs

Monitor Kubernetes resources:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: health-check
spec:
  template:
    spec:
      containers:
      - name: checker
        image: your-image
        command:
        - /bin/sh
        - -c
        - |
          if ! curl -f http://app:8080/health; then
            echo "Health check failed" | notifox send -a oncall -c sms
          fi
        env:
        - name: NOTIFOX_API_KEY
          valueFrom:
            secretKeyRef:
              name: notifox-secret
              key: api-key
      restartPolicy: Never
```

### Shell Scripts

Integrate into existing scripts:

```bash
#!/bin/bash
set -e

# Your script logic here
result=$(some-command)

if [ $? -ne 0 ]; then
    echo "Command failed: $result" | notifox send -a ops -c sms
    exit 1
fi
```

### Monitoring Scripts

Monitor system resources:

```bash
#!/bin/bash
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)

if (( $(echo "$CPU_USAGE > 90" | bc -l) )); then
    echo "High CPU usage: ${CPU_USAGE}%" | notifox send -a ops -c sms
fi
```

---

## Exit Codes

The CLI uses standard Unix exit codes:

- **0**: Success - Alert was sent successfully
- **1**: Failure - An error occurred

This makes it suitable for use in scripts and conditional logic:

```bash
if notifox send -a oncall -c sms -m "Test"; then
    echo "Alert sent successfully"
else
    echo "Failed to send alert"
    exit 1
fi
```

Common failure scenarios:
- Missing required flags or environment variables
- Invalid API key
- Network errors
- API errors (invalid audience, insufficient balance, etc.)
- Invalid channel value

---

## Troubleshooting

### CLI Not Found

**Problem:** `command not found: notifox`

**Solutions:**
1. Verify the binary is in your PATH: `which notifox`
2. Add the binary location to your PATH
3. Use the full path: `/usr/local/bin/notifox send -h`

### Permission Denied

**Problem:** `permission denied` when running the CLI

**Solution:**
```bash
chmod +x /usr/local/bin/notifox
```

### Environment Variables Not Working

**Problem:** Environment variables are not being picked up

**Solutions:**
1. Verify the variable is set: `echo $NOTIFOX_API_KEY`
2. Export the variable in your shell session
3. For persistent setup, add to `~/.bashrc`, `~/.zshrc`, or `/etc/environment`
4. In scripts, ensure variables are exported before running the CLI

### Stdin Not Being Read

**Problem:** Piped input is not being used

**Solutions:**
1. Ensure you're actually piping data: `echo "test" | notifox send -a x -c sms`
2. Check that stdin is a pipe, not a terminal
3. Use the `-m` flag if you want to provide the message directly

### API Connection Issues

**Problem:** Connection errors or timeouts

**Solutions:**
1. Check internet connectivity
2. Verify DNS resolution: `nslookup api.notifox.com`
3. Check firewall rules
4. Use `NOTIFOX_BASE_URL` if using a custom endpoint
5. Verify the API endpoint is accessible from your network

### Authentication Failures

**Problem:** Getting authentication errors

**Solutions:**
1. Verify API key is correct: `echo $NOTIFOX_API_KEY`
2. Check for extra spaces or newlines in the API key
3. Ensure the API key is active in your Notifox account
4. Regenerate the API key if necessary

### Invalid Audience Errors

**Problem:** "Invalid audience name" errors

**Solutions:**
1. Verify the audience name exists in your Notifox account
2. Check for typos in the audience name
3. Ensure the audience is properly configured
4. Use the exact audience name as shown in your Notifox dashboard

### Message Too Long

**Problem:** Messages are being truncated or failing

**Note:** The CLI doesn't enforce message length limits. If you encounter issues:
1. Check Notifox's message length limits for SMS vs Email
2. SMS messages may be split into multiple parts (shown in the "Parts" output)
3. Very long messages may need to be sent via email instead

---

## Best Practices

1. **Store API keys securely**: Use environment variables or secret management tools, never hardcode
2. **Use environment variables for defaults**: Set `NOTIFOX_AUDIENCE` and `NOTIFOX_CHANNEL` for common use cases
3. **Handle errors in scripts**: Check exit codes when using in automation
4. **Use appropriate channels**: SMS for urgent alerts, email for detailed reports
5. **Test before production**: Verify your setup with a test audience first
6. **Monitor usage**: Keep track of alert volume to avoid rate limits
7. **Use meaningful messages**: Include context like timestamps, server names, or error codes
8. **Set up monitoring**: Monitor the CLI itself to ensure alerts are being sent

---

## Additional Resources

- GitHub Repository: https://github.com/notifoxhq/notifox-cli
- Notifox Documentation: https://notifox.com/docs
- Notifox Go SDK: https://github.com/notifoxhq/notifox-go

---

## Support

For issues, questions, or contributions:
- Open an issue on GitHub: https://github.com/notifoxhq/notifox-cli/issues
- Check existing documentation and examples
- Review error messages for specific guidance
