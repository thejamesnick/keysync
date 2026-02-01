# KeySync Account & Workspace Model

## Purpose of an Account in KeySync

In KeySync, an **account** represents a **workspace namespace** for projects and encrypted secrets. It is intentionally lightweight and **not a traditional password-based user system**.

The account exists to:

- Uniquely identify a project owner
- Namespace projects in the cloud
- Enable remote collaboration
- Manage admin and developer SSH keys

Security **does not rely on the account itself**, but on SSH key cryptography.

---

## What an Account Is (and Is Not)

### ✅ What it IS

- A **workspace identifier** (email-based)
- A container for multiple projects
- A reference for cloud storage ownership
- A way to manage admin and developer keys

### ❌ What it is NOT

- Not a password-based auth system
- Not a secret holder
- Not involved in encryption/decryption

---

## Account Identity Model

### Account Identifier

- **Email address** (e.g. `alice@example.com`)
- Used only for:
  - Workspace naming
  - Project ownership
  - Cloud namespace isolation

### Account Authentication (Critical)

Accounts are authenticated using **SSH challenge–response**, not passwords.

**Why:**

- Email alone is not proof of ownership
- SSH keys already provide strong cryptographic identity

### Authentication Flow

1. User submits:
   - Email
   - Public SSH key (admin key)
2. Server issues a random challenge (nonce)
3. CLI signs the challenge using the private SSH key
4. Server verifies signature using the public key
5. Account is now authenticated and bound to that key

> The private key never leaves the user’s machine.

---

## Admin Keys vs Developer Keys

### Admin Keys

- Represent **ownership and control** of the account
- Can:
  - Create projects
  - Add/remove SSH keys
  - Upload and rotate secrets
  - Manage environments

Best practice:

- 1–2 admin keys per account
- One per trusted machine

### Developer Keys

- Granted access **per project/environment**
- Can:
  - Pull encrypted secrets
  - Decrypt locally using private key
- Cannot:
  - Modify keys
  - Upload secrets
  - Manage projects

---

## One Account → Many Projects

A single KeySync account can own **multiple independent projects**.

Example:

```
alice@example.com
├── api-server
├── mobile-app
├── landing-site
└── internal-tools
```

Each project is:

- Cryptographically isolated
- Independently managed
- Scoped by environment

---

## Cloud Namespace Structure

All data is stored under the account namespace:

```
/keysync/
└── alice@example.com/
    ├── api-server/
    │   ├── dev/
    │   │   ├── secrets.enc
    │   │   └── keys/
    │   └── prod/
    ├── mobile-app/
    └── landing-site/
```

Notes:

- Only encrypted blobs are stored
- Public keys are metadata only
- No plaintext secrets ever exist in cloud storage

---

## Local Account Representation

On the user’s machine:

```
~/.keysync/
├── account.json      # email, admin key fingerprint
├── projects/
│   ├── api-server/
│   ├── mobile-app/
│   └── landing-site/
└── cache/
```

The account file **does not contain secrets**.

---

## Multi-Machine & Recovery Model

An account may have multiple admin keys:

- Laptop admin key
- Desktop admin key
- Secure backup key (offline)

If a machine is lost:

1. Remove the compromised admin key
2. Rotate affected secrets
3. Access is fully revoked

---

## Why This Model Works

- No shared passwords
- No central trust in server
- SSH keys are proven, stable identity
- Cloud is zero-knowledge by design
- Scales from solo devs to remote teams

---

## Mental Model (TL;DR)

```
Email        → Workspace namespace
Admin key    → Ownership & control
Project      → Secret boundary
Environment  → Scope
SSH keys     → Access
Cloud        → Encrypted transport only
```

KeySync accounts are **identity and organization**, not security bottlenecks. Security lives entirely in cryptography and local private keys.

