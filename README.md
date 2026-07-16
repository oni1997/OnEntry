# OnEntry

A self-hosted, encrypted password manager I'm building as a learning project and because I wanted something I actually enjoy using every day.

## Why I'm building this

I've been wanting to learn Rust and Zig for a while, and I also wanted an excuse to get more comfortable with React, C#, and Go. Rather than build random todo apps, I figured I'd make something I'd actually use. Bitwarden is great, but I wanted my own setup that I could tweak and extend however I want. Plus, building a real system with encryption and multiple services sounded way more interesting than following tutorials.

## What it does

Right now it's an MVP. It stores passwords, generates strong ones, searches through them, and syncs across devices. Everything is encrypted. That's the goal, anyway.

## Tech stack

- React + TypeScript + TailwindCSS for the web UI
- Go for the API, auth, and vault management
- Rust for all the crypto stuff (Argon2, AES-256-GCM, password generation)
- Zig for utility jobs (backups, import/export, cleanup)
- C# WPF for a Windows desktop client
- PostgreSQL for storage
- Redis for caching (optional)
- Docker Compose to run it all

## Getting started

You'll need Docker and Docker Compose installed. Then:

```
docker compose up --build
```

The web UI runs on port 8081, the Go API on 8082, Rust crypto on 8083, Zig utilities on 8084, PostgreSQL on 5432, and Redis on 6379.

## Project structure

The code lives in these folders:

- apps/web - React frontend
- apps/desktop - C# Windows desktop client
- services/api-go - Go backend
- services/crypto-rust - Rust encryption service
- services/utility-zig - Zig utility service
- shared/proto - protocol buffer definitions
- shared/sdk - shared types
- database/migrations - PostgreSQL schema
- deploy/docker - Docker configs

## What's working

- User registration and login with JWT
- Vault encryption/decryption via Rust service
- Password generation with lots of options
- Search across your vault
- Dashboard with stats
- Import/export of vault data
- Basic rate limiting and audit logs

## What's still rough

This is very much an MVP. The desktop client is mostly scaffolding. Some features are simplified. The encryption flow works but I'm still testing edge cases. The UI is functional but basic. I'm actively working on it.

## Security notes

Master passwords never leave the client in plaintext. The vault is encrypted with AES-256-GCM using a key derived from the master password via Argon2. All crypto operations happen in the Rust service. The Go backend only ever sees encrypted data.

## Future plans

After the MVP stabilizes, I want to add:
- TOTP/2FA support
- Passkeys
- Browser extension
- Mobile app
- Secret sharing
- Team vaults

## Why this stack

Go feels like the right fit for the API layer - fast, simple, great concurrency. Rust is perfect for crypto since safety matters there. Zig is fun for utilities and makes me think about memory in a different way. React is what I know best for frontends. C# is just for the Windows desktop client because native Windows apps still matter. Each piece does what it's good at.

This is a learning project as much as it's a tool. I'm building it publicly so I can look back and see how far I've come with these languages.