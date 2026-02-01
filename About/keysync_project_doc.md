# KeySync Project Documentation

## Project Overview
**Name:** KeySync  
**Type:** Free Developer Tool / SSH-Native Secret Manager  
**Platform:** CLI first, optional web dashboard (future)  
**Infrastructure:** Maxcloud & NX for hosting, storage, and distribution  

**Tagline:** *Sync your secrets securely, SSH-style.*

**Description:**
KeySync is a free, developer-first tool that allows remote teams to manage environment variables and secrets securely. Using **SSH keys as identity**, only authorized team members can decrypt and access secrets. The tool enables full project and environment-based secret management, letting developers work offline, locally, or in a remote workspace while ensuring no plaintext secrets are ever exposed to servers or external services.

KeySync is designed to be simple, intuitive, and fast. Developers can upload an entire `.env` file at once, and secrets are encrypted per authorized key. Team admins manage onboarding and offboarding, controlling who can access which secrets for each project and environment.

---

## Core Features

1. **SSH-Native Identity**
   - Uses developer SSH keys for authentication.
   - Ensures per-dev access control.

2. **Project & Environment Based Secrets**
   - Supports multiple projects (repos/apps).
   - Supports multiple environments: dev, staging, prod, custom.
   - Secrets are scoped per environment.

3. **Full `.env` Upload**
   - Developers can upload an entire `.env` file at once.
   - CLI parses, encrypts, and stores each secret securely.

4. **Local-first Workflow**
   - Secrets decrypted locally on dev machine.
   - `.env` file generated or exported in-memory for shell usage.

5. **Access Management**
   - Admins can add/remove SSH public keys per project/environment.
   - Supports onboarding/offboarding securely.
   - Secret rotation supported after revocation.

6. **Encryption & Security**
   - End-to-end encryption using SSH-compatible methods (age / Go crypto).
   - Server or shared storage only stores encrypted blobs.
   - Zero knowledge: server never sees plaintext secrets.

7. **CLI Commands (MVP)**
   ```bash
   # Initialize project
   keysync init project-alpha

   # Upload .env
   keysync upload --project=project-alpha --env=dev .env

   # Pull secrets locally
   keysync pull --project=project-alpha --env=dev

   # Add a new key
   keysync add-key --project=project-alpha --env=dev alice.pub

   # Remove key / revoke access
   keysync remove-key --project=project-alpha --env=prod bob.pub

   # Rotate secrets
   keysync rotate --project=project-alpha --env=prod
   ```

8. **Tech Stack**
   - **CLI:** Go (cross-platform)
   - **Encryption:** age / Go crypto libraries
   - **Local Storage:** SQLite / JSON
   - **Server (optional):** Maxcloud / NX blob store for encrypted secrets
   - **CI/CD Integration:** GitHub Actions, GitLab, Vercel

9. **Naming & Conventions**
   - Secrets use uppercase letters and underscores.
   - Short and clear names (e.g., `DATABASE_URL`, `JWT_SECRET`).
   - Optional project prefixes for clarity if needed.

---

## Benefits

- Fully free and open-source CLI tool.
- Dev-native: integrates into existing SSH workflows.
- Local-first: devs can work offline securely.
- Secure per-dev access control, zero-knowledge server.
- Supports multiple projects, environments, and secret rotation.
- Fast onboarding/offboarding for remote teams.
- CI/CD ready for automated deployments.

---

## Future Enhancements (Optional)

- Web dashboard for admin & audit
- Secret versioning and rollback
- Conflict resolution for simultaneous updates
- Grouped namespaces for large projects
- Web3-specific secret management (wallets, RPC URLs)
- Multi-machine dev sync via Maxcloud

---

**Target Audience:** Remote dev teams, Web3 developers, startups, open-source projects, and anyone who wants a secure, developer-native workflow for managing secrets.

**Free & Open-Source Philosophy:** KeySync prioritizes simplicity, security, and accessibility, enabling dev teams to manage secrets without expensive subscriptions or complex infra.

