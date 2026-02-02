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
```bash
# Clone the repo
git clone https://github.com/keysync/cli.git
cd cli

# Install to /usr/local/bin
make install
```

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
