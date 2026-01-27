

# Notifox CLI

A small, Unix-friendly CLI for sending **internal alerts** through Notifox (SMS + Email).

The core goal is simple:

✅ **Pipe output from any command into an alert.**

```bash
echo "hello" | notifox send -a mathis -c email
```

---

## Installation

(TBD)

---

## Quick start

### 1) Set your API key
The CLI reads your API key from an environment variable:

```bash
export NOTIFOX_API_KEY="YOUR_API_KEY"
```

### 2) (Optional) Set default audience and channel
You can set default values for audience and channel (flags will override these):

```bash
export NOTIFOX_AUDIENCE="mathis"
export NOTIFOX_CHANNEL="email"
```

With these set, you can omit the flags:

```bash
echo "hello" | notifox send
```

### 3) Send an alert (stdin)
If text is piped into `notifox`, the CLI uses stdin as the message body:

```bash
echo "hello" | notifox send -a mathis -c email
```

This makes it work naturally with Linux / macOS workflows:

```bash
kubectl get pods -A | notifox send -a platform -c email
tail -n 200 app.log | notifox send -a oncall -c sms
```

### 4) Send an alert (message flag)
If you prefer sending a one-liner without piping stdin, use `-m / --message`:

```bash
notifox send -a mathis -c sms -m "DB is down"
```

---

## Command: `notifox send`

### Required flags

- `-a, --audience <name>`  
  The audience to send the alert to (example: `mathis`, `oncall`, `platform`)  
  Can also be set via `NOTIFOX_AUDIENCE` environment variable (flag overrides env var)

- `-c, --channel <sms|email>`  
  Which channel to send through  
  Can also be set via `NOTIFOX_CHANNEL` environment variable (flag overrides env var)

**Note:** Flags take precedence over environment variables (like AWS CLI behavior). If you set both, the flag value is used.

### Message input (stdin vs `-m`)
The message body can come from either:

1. `-m / --message`
2. stdin (if piped)

The precedence is:

✅ `--message` wins  
✅ otherwise, stdin is used  
❌ if neither is provided, the command fails with a helpful error

Examples:

```bash
# stdin
echo "backup failed" | notifox send -a oncall -c sms

# message flag
notifox send -a oncall -c sms -m "backup failed"
```

---

## Configuration

### Environment variables

- `NOTIFOX_API_KEY` *(required)*  
  API key used to authenticate requests

- `NOTIFOX_AUDIENCE` *(optional)*  
  Default audience to use if `-a/--audience` flag is not provided  
  Flags override environment variables (AWS CLI-style behavior)

- `NOTIFOX_CHANNEL` *(optional)*  
  Default channel to use if `-c/--channel` flag is not provided  
  Flags override environment variables (AWS CLI-style behavior)

Optional (mostly useful for development/testing):

- `NOTIFOX_BASE_URL` *(optional)*  
  Overrides the API base URL

### Examples

```bash
# Set defaults via environment variables
export NOTIFOX_API_KEY="your_key"
export NOTIFOX_AUDIENCE="oncall"
export NOTIFOX_CHANNEL="sms"

# Now you can send without flags
echo "server down" | notifox send

# Flags still override env vars
echo "server down" | notifox send -a mathis -c email  # Uses mathis/email, not oncall/sms
```

---

## What this is good for

The CLI is designed for **internal alerting**, not customer marketing messages.

Common use cases:

- cron jobs / scheduled scripts  
- CI pipelines (GitHub Actions, GitLab CI, etc.)
- server maintenance scripts
- Kubernetes jobs
- “send me the output of this command if it fails”

Examples:

```bash
# Send output of a command
df -h | notifox send -a ops -c email

# Send an error log excerpt
tail -n 50 /var/log/nginx/error.log | notifox send -a ops -c sms
```

---

## Internals / Architecture

The CLI intentionally follows a small layered design:

### CLI layer
- parses args (`send`, `--audience`, `--channel`, etc.)
- reads defaults from environment variables (`NOTIFOX_AUDIENCE`, `NOTIFOX_CHANNEL`)
- applies flag precedence (flags override env vars)
- reads message input from stdin or `--message`
- prints success/error output and exits with a proper exit code

### App layer
- validates inputs (missing flags, empty message, etc.)
- converts CLI input into a send request

### Notifox SDK layer
The CLI sends alerts using the official Notifox Go SDK:

- `notifox-go`: https://github.com/notifoxhq/notifox-go

This keeps the CLI thin and focused:
it mainly translates command-line input → SDK call.

---

## Exit codes

The CLI exits non-zero on failure, so it can be used reliably in scripts.

Examples:
- missing required flags
- missing message (no stdin + no `-m`)
- invalid API key
- API/network error

---

## Roadmap (future ideas)

- `notifox configure` (save config/profile)
- `notifox audiences list`
- `notifox credits`
- `--json` output mode
- `--subject` for email alerts
- `--dry-run` to print payload without sending
- `--file` to send file contents

---