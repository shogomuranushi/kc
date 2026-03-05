# kc

Stop writing secrets in `.env`. Let Keychain guard them with Touch ID.

## Problem

- API keys stored as plaintext in `.env` files
- One `.gitignore` mistake away from leaking secrets
- Sharing `.env` files over Slack or notes
- AI agents (Claude Code, Cline, etc.) get unrestricted access to your secrets
- A compromised dependency or malicious supply chain code can silently read `.env` at build time

## Solution

`kc` is a secret management CLI backed by macOS Keychain / Windows Credential Manager / Linux Secret Service.
Your `.env` contains only `kc://service/key` references — secrets are resolved from the Keychain at runtime.
Touch ID / password authentication is required to retrieve secrets, so nothing is accessed without human approval.

```bash
# .env (safe to commit to git)
ANTHROPIC_API_KEY=kc://anthropic/api_key
AWS_SECRET_ACCESS_KEY=kc://aws/secret_access_key
PORT=3000
```

```bash
# Just run commands through kc (Touch ID prompt appears)
kc claude
kc npm run dev
```

## Install

```bash
# macOS (Apple Silicon)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-darwin-arm64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc

# macOS (Intel)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-darwin-amd64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc

# Linux (amd64)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-linux-amd64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc
```

## Usage

### 1. Store secrets

```bash
kc set anthropic api_key sk-ant-xxxx    # Pass as argument
kc set github token                      # Interactive prompt (no echo)
echo "sk-xxx" | kc set stripe secret_key # Pipe
```

### 2. Write references in `.env`

```bash
ANTHROPIC_API_KEY=kc://anthropic/api_key
GITHUB_TOKEN=kc://github/token
PORT=3000
```

### 3. Run commands

```bash
kc claude                    # Auto-discovers .env and resolves secrets
kc npm run dev
kc docker compose up
kc run --env-file .env.prod -- npm start  # Explicit .env file
```

### Migrate existing `.env`

```bash
kc migrate            # Auto-discovers .env in current directory
kc migrate .env.local # Specify file
```

Interactively moves plaintext values into Keychain and rewrites `.env` with `kc://` references.

### Other commands

```bash
kc get github token          # Print value to stdout (pipe-friendly)
kc get github token | pbcopy # Copy to clipboard
kc delete github token       # Delete a secret
kc list                      # List all secrets
kc list aws                  # Filter by service
```

## Design

- **stdout is values only, all messages go to stderr** — safe for piping
- **`.env` is auto-discovered upward** — works from subdirectories
- **`kc <command>`** — anything not `set`/`get`/`delete`/`list`/`run`/`migrate` is executed as an external command
- **`kc://service/key`** — service/key allows `[a-zA-Z0-9_.-]` only
- **Keychain service name** — namespaced with `kc-cli:` prefix

## Threat model

| What it prevents | What it doesn't prevent |
|---|---|
| Leaking secrets via static `.env` scanning | Reading env vars from a running process |
| Accessing secrets without human approval | Malicious code reading env vars within a running process |
| Accidentally committing secrets to git | |

## Build

```bash
go build -o kc .
```
