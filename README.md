# Notifox CLI

Pipe output from any command into an alert. That's it.

```bash
echo "hello" | notifox send -a john -c email
```

## Installation

### macOS / Linux

```bash
curl -fsSL https://notifox.com/install.sh | sh
```

This installs to `~/.local/bin` (or `/usr/local/bin` if you have write access). Add to your PATH if needed.

**Or download manually** from the [releases page](https://github.com/notifoxhq/notifox-cli/releases).

### Windows
Download `notifox-cli_windows_amd64.zip` from the [releases](https://github.com/notifoxhq/notifox-cli/releases) page and extract `notifox.exe` to a folder on your PATH.

### Build from source

```bash
git clone https://github.com/notifoxhq/notifox-cli.git
cd notifox-cli
go build -ldflags "-X main.version=$(git describe --tags --always)" -o notifox .
```

## Setup

Set your API key:

```bash
export NOTIFOX_API_KEY="your_key"
```

Optionally set defaults so you don't have to type them every time:

```bash
export NOTIFOX_AUDIENCE="john"
export NOTIFOX_CHANNEL="sms"
```

Flags override environment variables, same as AWS CLI.

## Examples

Pipe stuff into it:

```bash
kubectl get pods -A | notifox send -a platform -c email
tail -n 200 app.log | notifox send -a oncall -c sms
```

Or use the message flag:

```bash
notifox send -a mathis -c sms -m "DB is down"
```

For email, you can set a subject with `-s` or `--subject` (SMS ignores it):

```bash
notifox send -a oncall -c email -s "Server alert" -m "Disk at 95%"
```

Add `-v` for verbose output (message ID, cost, parts).

If you set the env vars, you can skip the flags:

```bash
echo "server down" | notifox send
```

That's about it. Use it in cron jobs, CI pipelines, whatever. It's just a simple way to send alerts.
