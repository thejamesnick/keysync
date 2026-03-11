# KeySync: The SSH-Native Secret Manager 🔐

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Tagline:** Sync your secrets securely, SSH-style. 🔐

KeySync is a free, developer-first tool designed for remote teams to manage environment variables and secrets securely. It uses **SSH keys as identity**, ensuring that only authorized team members can decrypt and access secrets. 

**Zero Knowledge:** The server never sees plaintext secrets.
**Local-First:** Encryption happens on your machine.

---

## 🛠 Project Status
Current Phase: **Repo Setup & Architecture Plan**

See detailed documentation in the [`goal/`](./goal) folder:
*   [`goal/keysync.txt`](./goal/keysync.txt) - Core philosophy & overview
*   [`goal/tech-stack.txt`](./goal/tech-stack.txt) - Go & age encryption stack
*   [`goal/plan.txt`](./goal/plan.txt) - Build roadmap
*   [`goal/api.txt`](./goal/api.txt) - CLI & API reference
*   [`goal/analytics.txt`](./goal/analytics.txt) - Metadata-only analytics plan

---

## 🚀 Quick Start (Coming Soon)

### Installation
**One-line install (Mac & Linux):**
```bash
curl -sL https://raw.githubusercontent.com/thejamesnick/keysync/main/install.sh | bash
```

**One-line install (Windows, PowerShell):**
```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/thejamesnick/keysync/main/install.ps1 | iex"
```

**One-line install (Windows, PowerShell 7+):**
```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/thejamesnick/keysync/main/install.ps1 | iex"
```

**Install from source (Mac/Linux/Windows):**
```bash
go install github.com/thejamesnick/keysync@latest
```

**Build and run locally (Mac/Linux/Windows):**
```bash
go test ./...
go run ./cmd/keysync --help
```

### Windows (PowerShell) setup

**1) Install Go**

Option A (recommended):
```powershell
winget install -e --id GoLang.Go
```

Then restart your terminal and verify:
```powershell
go version
```

**2) Get KeySync**
```powershell
git clone https://github.com/thejamesnick/keysync.git
cd .\keysync
```

**3) Build/test locally**
```powershell
go test ./...
go run .\cmd\keysync --help
```

**4) SSH keys on Windows**
- Default SSH folder: `C:\Users\<you>\.ssh\`
- Common keys: `id_ed25519`, `id_rsa` (private) and `.pub` (public)
- If you don't have keys yet, Windows 10/11 typically includes OpenSSH:

```powershell
ssh-keygen -t ed25519
```

**Troubleshooting**
- If PowerShell says `go` is not recognized, your PATH hasn't refreshed—close and reopen the terminal (or sign out/in).

### Usage
```bash
# 1. Setup your identity
keysync generate --email me@example.com   # (If you don't have keys)
keysync signup --email me@example.com --me # Auto-finds your key

# 2. Create a project
keysync init

# 3. Add team members (Magic!)
keysync add-key github:username           # Import from GitHub
keysync add-key --me                      # Add yourself quickly
keysync add-key bob.pub                   # Or use a file

# 4. Push encrypted secrets
keysync push   # Encrypts .env -> secrets.enc
keysync pull   # Decrypts secrets.enc -> .env
```
**Find your own keys:**
```bash
keysync whoami
```

---

## 🏗 Architecture

### Account & Identity
*   **Authentication:** Challenge-response via SSH keys. No passwords.
*   **Access Control:** Per-project/environment authorization.

### Encryption Model
*   Uses **age** / Go crypto libraries.
*   Secrets are encrypted *independently* for every authorized public key.
*   Server stores only encrypted blobs.

---

## 🤝 Contributing
KeySync is built in public. Check out our [Build Plan](./goal/plan.txt) to see what we're working on next.

---

MIT License © 2026 KeySync
