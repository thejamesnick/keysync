# KeySync Product Documentation (Remote-First with SSH-Authenticated Accounts)

## Project Overview

**Name:** KeySync\
**Type:** Free Developer Tool / SSH-Native Secret Manager\
**Platform:** CLI first, optional web dashboard\
**Infrastructure:** Maxcloud & NX for hosting, encrypted storage, and distribution

**Tagline:** *Sync your secrets securely, anywhere, SSH-style.*

**Description:** KeySync is a free, developer-first tool designed for **remote teams** to manage environment variables and secrets securely. Using **SSH keys as identity**, only authorized team members can decrypt and access secrets. KeySync supports project and environment-based secret management, multi-machine workflows, and optional cloud syncing — making remote collaboration seamless while ensuring **zero-knowledge security**.

The system ties projects to **SSH-authenticated accounts** (using email as identifier), ensuring that project ownership is verified, and all encrypted secrets are safely shared across the team.

---

## Core Features

### 1. SSH-Authenticated Accounts

- Account is tied to **email** for identification and namespace
- Ownership is **verified via SSH challenge-response** (private key proves identity)
- Optional: OAuth login (GitHub/Google) for account creation
- Provides **project workspace mapping** and team management for cloud storage

### 2. Project & Environment-Based Secrets

- Multiple projects (repos/apps) supported
- Multiple environments: dev, staging, prod, custom
- Secrets scoped per environment and encrypted per dev key

### 3. Cloud Sync (Remote-First)

- Optional but **required for remote teams**
- Owner uploads encrypted secrets to Maxcloud/NX blob storage
- Only authorized devs with matching private keys can decrypt
- Enables multi-machine and cross-country collaboration

### 4. Full `.env` Upload

- Upload entire `.env` file at once
- CLI encrypts per authorized dev key and pushes to remote
- Generates local `.env` for dev or exports in-memory for shell usage

### 5. Access Management

- Admin can add/remove public keys per project/environment
- Supports onboarding/offboarding securely
- Secret rotation supported after revocation

### 6. Encryption & Security

- End-to-end encryption with `age` / Go crypto
- Server stores **only encrypted blobs** — zero knowledge
- SSH keys control access, not passwords
- Challenge-response ensures account authenticity

### 7. CLI Commands (MVP)

```bash
# Sign up with SSH-authenticated account
keysync signup --email alice@example.com --key ~/.ssh/id_rsa.pub

# Initialize project
keysync init project-alpha

# Add team member keys
keysync add-key bob.pub

# Upload .env to cloud
keysync upload .env --env=dev --remote

# Pull secrets remotely
keysync pull --project=project-alpha --env=dev --remote

# Remove a dev key / revoke access
keysync remove-key bob.pub --project=project-alpha --env=prod

# Rotate secrets
keysync rotate --project=project-alpha --env=prod
```

### 8. Tech Stack

- **CLI:** Go (cross-platform binaries)
- **Encryption:** age / Go crypto libraries
- **Local Storage:** SQLite / JSON
- **Cloud Storage:** Maxcloud / NX encrypted blob storage
- **Optional Web UI:** React / Next.js for dashboard, audit logs, key management
- **CI/CD Integration:** GitHub Actions, GitLab, Vercel

### 9. Naming & Conventions

- Uppercase letters and underscores: e.g., `DATABASE_URL`, `JWT_SECRET`
- Short and clear, optional project prefixes

---

## Remote-First Workflow

### 1. Account Creation & Authentication

1. User provides **email + public SSH key**
2. Server issues **challenge**
3. CLI signs challenge with private key
4. Server verifies signature → account is authenticated

### 2. Project Creation & Upload

1. Owner initializes project and environment
2. Adds authorized dev keys
3. Uploads `.env` → encrypted per dev key → stored in cloud

### 3. Dev Pull & Decrypt

1. Developer pulls encrypted secrets
2. CLI finds copy encrypted for their public key
3. Decrypts locally with private key → generates `.env` or exports to shell

### 4. Offboarding / Rotation

- Remove dev key → rotate affected secrets → revoked dev can no longer decrypt new secrets

---

## Benefits

- Fully free and secure for remote teams
- SSH keys as identity + authentication, zero-knowledge storage
- Multi-machine, multi-environment, multi-project ready
- Easy onboarding/offboarding
- Supports CI/CD injection securely
- Optional web dashboard for admin, audit, and rotation management

---

## Future Enhancements (Optional)

- Full web dashboard with project workspace view
- Secret versioning and rollback
- Conflict resolution for simultaneous updates
- Grouped namespaces for large projects
- Web3-specific secret management (wallets, RPC URLs)
- Team analytics & audit logs

---

**Target Audience:** Remote developer teams, startups, open-source projects, Web3 devs — anyone needing secure, SSH-based remote secret sharing.

**Security Philosophy:** Only authorized SSH keys can decrypt secrets. Cloud storage sees **only encrypted blobs**. Project ownership is verified via SSH challenge-response linked to the account email.

