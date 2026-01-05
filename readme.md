# 🔐 PassManager

A secure, efficient, and easy-to-use command-line password manager built with Go and PocketBase.

![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)

## 📋 Table of Contents

- [Features](#-features)
- [Security Architecture](#-security-architecture)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [PocketBase Setup](#-pocketbase-setup)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
- [Configuration](#-configuration)
- [Security Best Practices](#-security-best-practices)
- [Troubleshooting](#-troubleshooting)
- [Development](#-development)
- [Contributing](#-contributing)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔒 **Military-Grade Encryption** | AES-256-GCM authenticated encryption |
| 🔑 **Secure Key Derivation** | Argon2id (winner of Password Hashing Competition) |
| 📋 **Clipboard Integration** | Copy passwords without displaying them |
| 🎲 **Password Generator** | Cryptographically secure random passwords |
| 🔍 **Smart Search** | Search across titles, usernames, and URLs |
| 📁 **Categories** | Organize credentials by category |
| 🌐 **Self-Hosted Backend** | PocketBase for complete data ownership |
| 💻 **Cross-Platform** | Works on Linux, macOS, and Windows |
| 🚀 **Fast & Lightweight** | Single binary, minimal dependencies |
| 📝 **Secure Notes** | Encrypted notes for each credential |

---

## 🏗 Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        PassManager CLI                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐   │
│  │   Master    │───▶│   Argon2id   │───▶│  Derived Key    │   │
│  │  Password   │    │     KDF      │    │   (256-bit)     │   │
│  └─────────────┘    └──────────────┘    └────────┬────────┘   │
│                                                   │             │
│                            ┌──────────────────────┘             │
│                            ▼                                    │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐   │
│  │  Plaintext  │───▶│  AES-256-GCM │───▶│   Ciphertext    │   │
│  │  Password   │    │  Encryption  │    │   (Base64)      │   │
│  └─────────────┘    └──────────────┘    └────────┬────────┘   │
│                                                   │             │
└───────────────────────────────────────────────────┼─────────────┘
                                                    │
                                                    ▼
                                          ┌─────────────────┐
                                          │   PocketBase    │
                                          │    Database     │
                                          │  (Encrypted     │
                                          │   Data Only)    │
                                          └─────────────────┘
```

### Cryptographic Details

| Component | Algorithm | Parameters |
|-----------|-----------|------------|
| **Key Derivation** | Argon2id | Time: 3, Memory: 64MB, Threads: 4, KeyLen: 32 |
| **Encryption** | AES-256-GCM | 256-bit key, 96-bit nonce, authenticated |
| **Salt** | CSPRNG | 128 bits (16 bytes) |
| **Password Hash** | SHA-256 | Of derived key (for verification only) |

### Security Properties

- ✅ **Zero-Knowledge**: Server never sees plaintext passwords
- ✅ **Forward Secrecy**: Unique nonce per encryption
- ✅ **Authenticated Encryption**: Detects tampering
- ✅ **Memory-Hard KDF**: Resistant to GPU/ASIC attacks
- ✅ **No Password Storage**: Master password never stored

---

## 📦 Prerequisites

### Required

- **Go** 1.21 or higher
- **PocketBase** 0.20 or higher

### Optional

- **xclip** or **xsel** (Linux clipboard support)
- **pbcopy** (macOS - included by default)

### Check Prerequisites

```bash
# Check Go version
go version

# Check PocketBase
./pocketbase --version
```

---

## 🚀 Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/anujbarve/passmanager.git
cd passmanager

# Download dependencies
go mod download

# Build the binary
go build -ldflags="-s -w" -o passmanager

# (Optional) Install globally
sudo mv passmanager /usr/local/bin/
```

### Option 2: Go Install

```bash
go install github.com/anujbarve/passmanager@latest
```

### Option 3: Download Pre-built Binary

```bash
# Linux (amd64)
curl -L https://github.com/anujbarve/passmanager/releases/latest/download/passmanager-linux-amd64 -o passmanager
chmod +x passmanager

# macOS (amd64)
curl -L https://github.com/anujbarve/passmanager/releases/latest/download/passmanager-darwin-amd64 -o passmanager
chmod +x passmanager

# Windows (amd64)
curl -L https://github.com/anujbarve/passmanager/releases/latest/download/passmanager-windows-amd64.exe -o passmanager.exe
```

### Verify Installation

```bash
passmanager --version
passmanager --help
```

---

## 🗄 PocketBase Setup

### Step 1: Download PocketBase

```bash
# Linux
wget https://github.com/pocketbase/pocketbase/releases/latest/download/pocketbase_0.22.0_linux_amd64.zip
unzip pocketbase_0.22.0_linux_amd64.zip

# macOS
wget https://github.com/pocketbase/pocketbase/releases/latest/download/pocketbase_0.22.0_darwin_amd64.zip
unzip pocketbase_0.22.0_darwin_amd64.zip

# Windows
# Download from: https://github.com/pocketbase/pocketbase/releases
```

### Step 2: Start PocketBase

```bash
./pocketbase serve
```

PocketBase will start on `http://127.0.0.1:8090`

### Step 3: Create Admin Account

1. Open http://127.0.0.1:8090/_/ in your browser
2. Create your admin account (email + password)

### Step 4: Create Collections

Navigate to **Collections** in the admin panel and create:

#### Collection 1: `vault_config`

| Field Name | Type | Required | Options |
|------------|------|----------|---------|
| `salt` | Text | ✅ | - |
| `password_hash` | Text | ✅ | - |

**API Rules (set all to empty for admin-only access):**
- List/Search: `@request.auth.id != ""`
- View: `@request.auth.id != ""`
- Create: `@request.auth.id != ""`
- Update: `` (disabled)
- Delete: `` (disabled)

#### Collection 2: `credentials`

| Field Name | Type | Required | Options |
|------------|------|----------|---------|
| `title` | Text | ✅ | Max: 255 |
| `username` | Text | ❌ | Max: 255 |
| `encrypted_password` | Text | ✅ | - |
| `url` | URL | ❌ | - |
| `notes` | Text | ❌ | - |
| `category` | Text | ❌ | Max: 50 |

**API Rules:**
- List/Search: `@request.auth.id != ""`
- View: `@request.auth.id != ""`
- Create: `@request.auth.id != ""`
- Update: `@request.auth.id != ""`
- Delete: `@request.auth.id != ""`

### Step 5: (Alternative) Import Schema

Save this as `pb_schema.json` and import via Admin UI → Settings → Import Collections:

```json
[
    {
        "name": "vault_config",
        "type": "base",
        "schema": [
            {
                "name": "salt",
                "type": "text",
                "required": true
            },
            {
                "name": "password_hash",
                "type": "text",
                "required": true
            }
        ]
    },
    {
        "name": "credentials",
        "type": "base",
        "schema": [
            {
                "name": "title",
                "type": "text",
                "required": true,
                "options": {
                    "max": 255
                }
            },
            {
                "name": "username",
                "type": "text",
                "required": false,
                "options": {
                    "max": 255
                }
            },
            {
                "name": "encrypted_password",
                "type": "text",
                "required": true
            },
            {
                "name": "url",
                "type": "url",
                "required": false
            },
            {
                "name": "notes",
                "type": "text",
                "required": false
            },
            {
                "name": "category",
                "type": "text",
                "required": false,
                "options": {
                    "max": 50
                }
            }
        ]
    }
]
```

---

## ⚡ Quick Start

```bash
# 1. Start PocketBase (in a separate terminal)
./pocketbase serve

# 2. Initialize your vault
passmanager init

# 3. Add your first credential
passmanager add -t "GitHub" -u "your@email.com" -l "https://github.com" -g

# 4. List all credentials
passmanager list

# 5. Retrieve a credential
passmanager get -i <credential-id> -c
```

---

## 📖 Usage

### Initialize Vault

Set up the password manager for first use:

```bash
passmanager init
```

You'll be prompted for:
- PocketBase URL (e.g., `http://127.0.0.1:8090`)
- Admin email and password
- Master password (minimum 12 characters)

### Add Credential

```bash
# Interactive mode
passmanager add -t "Service Name"

# With all options
passmanager add \
  --title "GitHub" \
  --username "your@email.com" \
  --password "your-password" \
  --url "https://github.com" \
  --notes "Personal account" \
  --category "development"

# Generate random password
passmanager add -t "AWS Console" -u "admin" -g --length 24

# Short flags
passmanager add -t "Twitter" -u "myhandle" -l "https://twitter.com" -c "social" -g
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--title` | `-t` | Credential title (required) | - |
| `--username` | `-u` | Username or email | - |
| `--password` | `-p` | Password (prompts if not provided) | - |
| `--url` | `-l` | Website URL | - |
| `--notes` | `-n` | Additional notes (encrypted) | - |
| `--category` | `-c` | Category for organization | "general" |
| `--generate` | `-g` | Generate random password | false |
| `--length` | - | Generated password length | 20 |

### List Credentials

```bash
# List all
passmanager list

# Search by title, username, or URL
passmanager list -s "github"
passmanager list --search "email@example.com"
```

**Output:**
```
🔐 Stored Credentials
=====================
ID                   TITLE                     USERNAME                       CATEGORY
-------------------- ------------------------- ------------------------------ ---------------
abc123def456         GitHub                    your@email.com                 development
xyz789ghi012         AWS Console               admin                          cloud

Total: 2 credential(s)
```

### Get Credential

```bash
# View credential (password hidden)
passmanager get -i abc123def456

# Show password in terminal
passmanager get -i abc123def456 --show

# Copy password to clipboard
passmanager get -i abc123def456 --copy

# Both show and copy
passmanager get -i abc123def456 -s -c
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--id` | `-i` | Credential ID (required) | - |
| `--show` | `-s` | Display password | false |
| `--copy` | `-c` | Copy password to clipboard | false |

### Delete Credential

```bash
passmanager delete -i abc123def456
```

You'll be prompted for confirmation.

### Generate Password

Generate secure passwords without storing them:

```bash
# Default (20 chars with symbols)
passmanager generate

# Custom length
passmanager generate -l 32

# Without symbols
passmanager generate -l 16 --symbols=false

# Copy to clipboard
passmanager generate -l 24 -c
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--length` | `-l` | Password length | 20 |
| `--symbols` | `-s` | Include symbols | true |
| `--copy` | `-c` | Copy to clipboard | false |

---

## ⚙️ Configuration

Configuration is stored at `~/.passmanager/config.json`:

```json
{
  "pocketbase_url": "http://127.0.0.1:8090",
  "admin_email": "admin@example.com",
  "session_timeout_minutes": 5
}
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `pocketbase_url` | PocketBase server URL | - |
| `admin_email` | Admin email for authentication | - |
| `session_timeout_minutes` | Session cache duration | 5 |

### Environment Variables

```bash
# Override PocketBase URL
export PASSMANAGER_PB_URL="http://localhost:8090"

# Enable debug mode
export PASSMANAGER_DEBUG="true"
```

---

## 🛡 Security Best Practices

### Master Password Guidelines

- ✅ Use at least 16 characters
- ✅ Include uppercase, lowercase, numbers, and symbols
- ✅ Use a passphrase (e.g., "correct-horse-battery-staple")
- ❌ Don't reuse passwords from other services
- ❌ Don't store master password digitally

### Operational Security

```bash
# Clear clipboard after use (Linux)
sleep 30 && xclip -selection clipboard < /dev/null

# Clear bash history
history -c && history -w

# Use in private terminal
# Avoid using in shared screen sessions
```

### PocketBase Hardening

```bash
# Run with HTTPS (production)
./pocketbase serve --https=yourdomain.com

# Restrict to localhost only
./pocketbase serve --http=127.0.0.1:8090

# Regular backups
./pocketbase backup
```

### Data Backup

```bash
# Backup PocketBase data
cp -r pb_data/ pb_data_backup_$(date +%Y%m%d)/

# Backup config
cp ~/.passmanager/config.json ~/.passmanager/config.json.backup
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. "Vault not initialized"

```bash
# Solution: Run init command
passmanager init
```

#### 2. "Failed to connect to PocketBase"

```bash
# Check if PocketBase is running
curl http://127.0.0.1:8090/api/health

# Start PocketBase
./pocketbase serve
```

#### 3. "Authentication failed"

```bash
# Verify credentials in PocketBase Admin UI
# Check PocketBase version for correct auth endpoint

# PocketBase 0.23+: Uses _superusers collection
# PocketBase < 0.23: Uses /api/admins/auth-with-password
```

#### 4. "Collection not found"

```bash
# Create required collections in PocketBase Admin UI
# See "PocketBase Setup" section
```

#### 5. "Invalid master password"

```
# Master password cannot be recovered
# You must reinitialize the vault (data will be lost)
passmanager init
```

#### 6. "Clipboard not working" (Linux)

```bash
# Install xclip
sudo apt install xclip

# Or xsel
sudo apt install xsel
```

### Debug Mode

```bash
# Enable verbose output
export PASSMANAGER_DEBUG=true
passmanager list
```

### Logs

```bash
# PocketBase logs
tail -f pb_data/logs.db

# Application logs (if enabled)
tail -f ~/.passmanager/passmanager.log
```

---

## 🛠 Development

### Project Structure

```
passmanager/
├── cmd/
│   ├── root.go           # Root command & CLI setup
│   ├── init.go           # Vault initialization
│   ├── add.go            # Add credentials
│   ├── get.go            # Retrieve credentials
│   ├── list.go           # List credentials
│   ├── delete.go         # Delete credentials
│   └── generate.go       # Password generator
├── internal/
│   ├── crypto/
│   │   └── crypto.go     # Encryption & key derivation
│   ├── database/
│   │   └── pocketbase.go # PocketBase client
│   ├── models/
│   │   └── credential.go # Data models
│   └── config/
│       └── config.go     # Configuration management
├── main.go               # Entry point
├── go.mod
├── go.sum
└── README.md
```

### Building

```bash
# Development build
go build -o passmanager

# Production build (optimized)
go build -ldflags="-s -w" -o passmanager

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o passmanager-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o passmanager-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o passmanager-windows-amd64.exe
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/crypto/...

# Verbose output
go test -v ./...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Dependencies

```bash
# Update dependencies
go get -u ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

### 1. Fork the Repository

```bash
git clone https://github.com/anujbarve/passmanager.git
cd passmanager
```

### 2. Create a Branch

```bash
git checkout -b feature/amazing-feature
```

### 3. Make Changes

- Follow Go best practices
- Add tests for new features
- Update documentation

### 4. Commit Changes

```bash
git commit -m "feat: add amazing feature"
```

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance

### 5. Push and Create PR

```bash
git push origin feature/amazing-feature
```

### Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the problem, not the person

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 Acknowledgments

- [PocketBase](https://pocketbase.io/) - Backend database
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Go Crypto](https://pkg.go.dev/golang.org/x/crypto) - Cryptographic primitives
- [Argon2](https://github.com/P-H-C/phc-winner-argon2) - Password hashing algorithm

---

## 📊 Roadmap

- [ ] **v1.1** - Session caching (avoid repeated password entry)
- [ ] **v1.2** - Import/Export (encrypted JSON/CSV)
- [ ] **v1.3** - Password strength analyzer
- [ ] **v1.4** - TOTP/2FA support
- [ ] **v1.5** - Browser extension integration
- [ ] **v2.0** - Multi-user support with sharing

---

## 📞 Support

- 📧 Email: anujbarve27[at]gmail[dot]com
- 🐛 Issues: [GitHub Issues](https://github.com/anujbarve/passmanager/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/anujbarve/passmanager/discussions)

---

<p align="center">
  Made with ❤️ and 🔐 by <a href="https://github.com/anujbarve">Anuj Barve</a>
</p>

<p align="center">
  <a href="#-passmanager">Back to top ⬆️</a>
</p>


---

## Additional Files to Include

### .gitignore

```gitignore
# Binaries
passmanager
passmanager.exe
*.exe
*.dll
*.so
*.dylib

# Build output
/dist/
/bin/

# Test binary
*.test

# Coverage
coverage.out
coverage.html

# Go workspace
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# PocketBase
pb_data/
pocketbase

# Config (contains sensitive paths)
.passmanager/

# Environment
.env
.env.local
```

### Makefile

```makefile
.PHONY: build clean test lint run install

BINARY_NAME=passmanager
VERSION=$(shell git describe --tags --always --dirty)
BUILD_FLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

run: build
	./$(BINARY_NAME)

install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Cross-compilation
build-all: clean
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe
```

### LICENSE

```
MIT License

Copyright (c) 2024 Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

