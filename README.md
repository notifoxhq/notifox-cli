# Notifox CLI

Pipe output from any command into an alert. That's it.

```bash
echo "hello" | notifox send -a mathis -c email
```

## Installation

Download from the [releases page](https://github.com/notifoxhq/notifox-cli/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_darwin_arm64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_darwin_amd64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/

# Linux
curl -L https://github.com/notifoxhq/notifox-cli/releases/latest/download/notifox-cli_linux_amd64.tar.gz | tar -xz
sudo mv notifox /usr/local/bin/
```

Windows users: download the zip from the releases page and extract it manually.

Or build from source:

```bash
git clone https://github.com/notifoxhq/notifox-cli.git
cd notifox-cli
make build  # or: go build -ldflags "-X main.version=$(git describe --tags --always)" -o notifox .
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

If you set the env vars, you can skip the flags:

```bash
echo "server down" | notifox send
```

That's about it. Use it in cron jobs, CI pipelines, whatever. It's just a simple way to send alerts.
