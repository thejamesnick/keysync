# KeySync API & CLI Reference Book 🔐

This document is the **single source of truth** for how KeySync works from a developer’s point of view.

It defines:

- The philosophy behind the API
- Authentication and security model
- All API endpoints
- The CLI commands and how they map to the API

This is written with one goal:

> **If you understand this book, you understand KeySync.**

---

## 1. Design Principles

KeySync is built on a few non‑negotiable principles:

1. **Secrets are encrypted locally**
2. **The server never sees plaintext**
3. **SSH keys are identity**
4. **Projects are metadata, secrets are blobs**
5. **The CLI is the primary interface**

The API exists to coordinate access, not to hold power.

---

## 2. Core Concepts

### Account

An account represents a human or organization.

- Identified by a verified email
- Owns one or more projects
- Has one or more **admin public keys**

---

### Project

A project is a container for secrets.

Each project has:

- A unique project ID
- A name (human‑friendly)
- One or more environments (dev, prod, etc.)
- A list of authorized public keys

Projects contain **no plaintext secrets**.

---

### Environment

Environments are logical namespaces inside a project.

Examples:

- dev
- staging
- prod

Each environment maps to one encrypted blob per version.

---

### Blob

A blob is an encrypted binary object.

- Produced on the developer machine
- Stored as‑is on the server
- Never decrypted server‑side

---

## 3. Authentication Model

KeySync uses **challenge–response authentication** with SSH keys.

### Login Flow

1. CLI requests a challenge
2. Server returns a nonce
3. CLI signs the nonce using the SSH private key
4. Server verifies the signature
5. Server issues a short‑lived session token

No passwords. No permanent tokens.

---

## 4. API Overview

Base URL:

```
https://api.keysync.dev
```

All requests:

- Use HTTPS
- Are signed or token‑authenticated
- Return JSON

---

## 5. Authentication Endpoints

### POST /auth/challenge

Request a login challenge.

**Request**

```
{ "email": "user@example.com" }
```

**Response**

```
{ "challenge": "random_nonce" }
```

---

### POST /auth/verify

Verify signed challenge.

**Request**

```
{
  "email": "user@example.com",
  "public_key": "ssh-ed25519 AAAA...",
  "signature": "base64_signature"
}
```

**Response**

```
{ "token": "session_token" }
```

---

## 6. Project Endpoints

### GET /projects

List all projects accessible by the account.

**Response**

```
[
  { "id": "proj_1", "name": "payments-api" },
  { "id": "proj_2", "name": "auth-service" }
]
```

---

### POST /projects

Create a new project.

**Request**

```
{ "name": "my-project" }
```

---

### GET /projects/{id}

Get project metadata.

---

### POST /projects/{id}/keys

Add a public key to a project.

**Request**

```
{ "public_key": "ssh-ed25519 AAAA...", "role": "dev" }
```

---

### DELETE /projects/{id}/keys/{fingerprint}

Revoke access for a key.

---

## 7. Secret Endpoints

### POST /secrets/push

Upload an encrypted blob.

**Request**

```
{
  "project_id": "proj_1",
  "environment": "prod",
  "blob": "base64_encrypted_data",
  "checksum": "sha256"
}
```

---

### GET /secrets/pull

Retrieve the latest encrypted blob.

**Response**

```
{
  "blob": "base64_encrypted_data",
  "checksum": "sha256"
}
```

---

## 8. CLI Overview

The CLI is the primary way users interact with KeySync.

Command structure:

```
keysync <resource> <action> [flags]
```

---

## 9. Authentication Commands

### keysync login

Authenticate the CLI using SSH.

```
keysync login
```

---

## 10. Project Commands

### List projects

```
keysync projects
```

---

### Create project

```
keysync project create my-project
```

---

### Project info

```
keysync project info my-project
```

---

## 11. Key Management Commands

### Add key

```
keysync project key add my-project ~/.ssh/id_ed25519.pub
```

---

### Remove key

```
keysync project key remove my-project <fingerprint>
```

---

## 12. Secret Commands

### Push secrets

```
keysync secrets push --env prod
```

---

### Pull secrets

```
keysync secrets pull --env prod
```

---

## 13. Error Model

All errors are explicit and human‑readable.

Example:

```
ERROR: access denied — key not authorized for this project
```

No silent failures.

---

## 14. Security Guarantees

KeySync guarantees:

- Server cannot decrypt secrets
- Revoked keys lose access instantly
- Encrypted blobs are immutable
- All actions are auditable (metadata only)

---

## 15. Philosophy (Final Word)

KeySync is intentionally boring.

No dashboards. No clutter. No magic.

Just:

- Keys
- Projects
- Secrets

> **The best security system is the one you don’t have to think about.**

End of book.

