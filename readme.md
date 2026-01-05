# 🔐 PassManager

A secure, efficient, and easy-to-use **interactive** command-line password manager built with Go and PocketBase.

![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)

```
   ____                __  __                                   
  |  _ \ __ _ ___ ___ |  \/  | __ _ _ __   __ _  __ _  ___ _ __ 
  | |_) / _' / __/ __|| |\/| |/ _' | '_ \ / _' |/ _' |/ _ \ '__|
  |  __/ (_| \__ \__ \| |  | | (_| | | | | (_| | (_| |  __/ |   
  |_|   \__,_|___/___/|_|  |_|\__,_|_| |_|\__,_|\__, |\___|_|   
                                                |___/           
```

## 📋 Table of Contents

- [Features](#-features)
- [Demo](#-demo)
- [Security Architecture](#-security-architecture)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [PocketBase Setup](#-pocketbase-setup)
- [Quick Start](#-quick-start)
- [Interactive Interface](#-interactive-interface)
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
| 🖥️ **Interactive CLI** | Beautiful menu-driven interface with arrow key navigation |
| 🔒 **Military-Grade Encryption** | AES-256-GCM authenticated encryption |
| 🔑 **Secure Key Derivation** | Argon2id (winner of Password Hashing Competition) |
| ⏱️ **Session Management** | Auto-lock vault after configurable timeout |
| 📋 **Clipboard Integration** | Copy passwords with auto-clear timeout |
| 🎲 **Password Generator** | Cryptographically secure random passwords |
| 🔍 **Smart Search** | Search across titles, usernames, and URLs |
| 📁 **Categories** | Organize credentials by category |
| ✏️ **Edit Credentials** | Modify existing passwords and details |
| 📤 **Export Vault** | Create encrypted backups |
| 🔐 **Change Master Password** | Re-encrypt all data with new password |
| 🌐 **Self-Hosted Backend** | PocketBase for complete data ownership |
| 💻 **Cross-Platform** | Works on Linux, macOS, and Windows |
| 🚀 **Fast & Lightweight** | Single binary, minimal dependencies |
| ⚙️ **Configurable Settings** | Customize timeouts, defaults, and more |

---

## 🎬 Demo

### Main Menu
```
╭──────────────────────────────────────────────────────────────────╮
│                                                                  │
│     ____                __  __                                   │
│    |  _ \ __ _ ___ ___ |  \/  | __ _ _ __   __ _  __ _  ___ _ __ │
│    | |_) / _' / __/ __|| |\/| |/ _' | '_ \ / _' |/ _' |/ _ \ '__││
│    |  __/ (_| \__ \__ \| |  | | (_| | | | | (_| | (_| |  __/ |   │
│    |_|   \__,_|___/___/|_|  |_|\__,_|_| |_|\__,_|\__, |\___|_|   │
│                                                  |___/           │
│                                                                  │
│    Secure Password Manager v1.0.0                                │
│    Your passwords, encrypted locally, stored securely.          │
│                                                                  │
│  🔓 Session active (expires in 4m 32s)                           │
│                                                                  │
│  🔐 Main Menu                                                    │
│  ────────────────────────────────────────────────────────────    │
│                                                                  │
│  ▸ ➕  Add Credential (Store a new password)                     │
│    📋  List Credentials (View all stored passwords)              │
│    🔍  Search Credentials (Find a specific password)             │
│    🔑  Get Credential (Retrieve a password by ID)                │
│    ✏️   Edit Credential (Modify an existing password)             │
│    🎲  Generate Password (Create a secure password)              │
│    🗑️   Delete Credential (Remove a stored password)              │
│    🔐  Change Master Password (Update your master password)      │
│    📤  Export Vault (Export encrypted backup)                    │
│    🔒  Lock Vault (Lock and require re-authentication)           │
│    ⚙️   Settings (Configure application settings)                 │
│    ❓  Help (Show help information)                              │
│    🚪  Exit (Close the application)                              │
│                                                                  │
╰──────────────────────────────────────────────────────────────────╯
```

### Credential List
```
  ╭─────────────────┬───────────────────────────┬────────────────────────────────┬─────────────────╮
  │ ID              │ TITLE                     │ USERNAME                       │ CATEGORY        │
  ├─────────────────┼───────────────────────────┼────────────────────────────────┼─────────────────┤
  │ abc123def456    │ GitHub                    │ user@example.com               │ development     │
  │ xyz789ghi012    │ AWS Console               │ admin                          │ cloud           │
  │ mno345pqr678    │ Gmail                     │ myemail@gmail.com              │ personal        │
  ╰─────────────────┴───────────────────────────┴────────────────────────────────┴─────────────────╯
```

### Credential Card
```
  ┌─────────────────────────────────────────────────┐
  │ GitHub                                          │
  ├─────────────────────────────────────────────────┤
  │  ID:          abc123def456                      │
  │  Username:    user@example.com                  │
  │  URL:         https://github.com                │
  │  Category:    development                       │
  │  Password:    ••••••••••••                      │
  └─────────────────────────────────────────────────┘
```

---

## 🏗 Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     PassManager Interactive CLI                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐    │
│  │   Master    │───▶│   Argon2id   │───▶│  Derived Key    │    │
│  │  Password   │    │     KDF      │    │   (256-bit)     │    │
│  └─────────────┘    └──────────────┘    └────────┬────────┘    │
│                                                   │              │
│  ┌─────────────┐                    ┌─────────────▼──────────┐  │
│  │   Session   │◀───────────────────│   Crypto Service       │  │
│  │   Manager   │    (Holds Key)     │   (In-Memory Only)     │  │
│  └──────┬──────┘                    └────────────────────────┘  │
│         │                                        │              │
│         │ Auto-Lock                              │              │
│         │ After Timeout                          ▼              │
│         │                           ┌──────────────────────┐    │
│         │                           │    AES-256-GCM       │    │
│         │                           │    Encrypt/Decrypt   │    │
│         ▼                           └──────────┬───────────┘    │
│  ┌─────────────┐                               │                │
│  │   Secure    │                               │                │
│  │   Logout    │                               ▼                │
│  │  (Zero Key) │                    ┌──────────────────────┐    │
│  └─────────────┘                    │   Encrypted Data     │    │
│                                     │   (Base64 Encoded)   │    │
│                                     └──────────┬───────────┘    │
│                                                │                │
└────────────────────────────────────────────────┼────────────────┘
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
- ✅ **Session Auto-Lock**: Automatic lockout after inactivity
- ✅ **Secure Memory Clear**: Keys zeroed on logout
- ✅ **Clipboard Auto-Clear**: Passwords removed from clipboard after timeout

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
go mod tidy

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

# macOS (arm64 / Apple Silicon)
curl -L https://github.com/anujbarve/passmanager/releases/latest/download/passmanager-darwin-arm64 -o passmanager
chmod +x passmanager

# Windows (amd64)
curl -L https://github.com/anujbarve/passmanager/releases/latest/download/passmanager-windows-amd64.exe -o passmanager.exe
```

### Verify Installation

```bash
./passmanager
```

---

## 🗄 PocketBase Setup

### Step 1: Download PocketBase

```bash
# Linux (amd64)
wget https://github.com/pocketbase/pocketbase/releases/download/v0.23.4/pocketbase_0.23.4_linux_amd64.zip
unzip pocketbase_0.23.4_linux_amd64.zip

# macOS (amd64)
wget https://github.com/pocketbase/pocketbase/releases/download/v0.23.4/pocketbase_0.23.4_darwin_amd64.zip
unzip pocketbase_0.23.4_darwin_amd64.zip

# macOS (arm64)
wget https://github.com/pocketbase/pocketbase/releases/download/v0.23.4/pocketbase_0.23.4_darwin_arm64.zip
unzip pocketbase_0.23.4_darwin_arm64.zip

# Windows - Download from:
# https://github.com/pocketbase/pocketbase/releases
```

### Step 2: Start PocketBase

```bash
./pocketbase serve
```

PocketBase will start on `http://127.0.0.1:8090`

### Step 3: Create Admin Account

1. Open http://127.0.0.1:8090/_/ in your browser
2. Create your superuser/admin account (email + password)

### Step 4: Create Collections

Navigate to **Collections** in the admin panel and create:

#### Collection 1: `vault_config`

| Field Name | Type | Required |
|------------|------|----------|
| `salt` | Plain text | ✅ |
| `password_hash` | Plain text | ✅ |

**API Rules:** Leave all empty (admin-only access)

#### Collection 2: `credentials`

| Field Name | Type | Required | Options |
|------------|------|----------|---------|
| `title` | Plain text | ✅ | Max: 255 |
| `username` | Plain text | ❌ | Max: 255 |
| `encrypted_password` | Plain text | ✅ | - |
| `url` | URL | ❌ | - |
| `notes` | Plain text | ❌ | - |
| `category` | Plain text | ❌ | Max: 50 |

**API Rules:** Leave all empty (admin-only access)

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
                "options": {"max": 255}
            },
            {
                "name": "username",
                "type": "text",
                "required": false,
                "options": {"max": 255}
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
                "options": {"max": 50}
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

# 2. Run PassManager
./passmanager

# 3. First run will launch the Setup Wizard
#    - Enter PocketBase URL
#    - Enter Admin credentials
#    - Create Master Password

# 4. After setup, you'll see the Main Menu
#    - Use arrow keys to navigate
#    - Press Enter to select
#    - Ctrl+C to exit
```

---

## 🖥️ Interactive Interface

### First Time Setup

On first run, PassManager will guide you through setup:

```
🔐 Password Manager Setup
========================

→ Let's set up your secure password vault.

PocketBase URL: http://127.0.0.1:8090
✓ Connected to PocketBase

Admin Email: admin@example.com
Admin Password: ••••••••
✓ Authenticated successfully

→ Now create your master password.
  This password encrypts all your data locally.
  It cannot be recovered if lost!

Master Password (min 12 chars): ••••••••••••
Confirm Master Password: ••••••••••••

✓ Vault created successfully!
⚠ IMPORTANT: Remember your master password!
```

### Unlocking Vault

```
▶ Unlock Vault
──────────────────────────────────────────────────────────────

Admin Password: ••••••••
⠋ Connecting...
Master Password: ••••••••••••

✓ Vault unlocked!
```

### Menu Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate menu items |
| `Enter` | Select option |
| `Ctrl+C` | Lock vault and exit |
| Type | Filter menu items |

### Add Credential

```
▶ Add New Credential
──────────────────────────────────────────────────────────────

Title: GitHub
Username/Email: user@example.com
URL: https://github.com
Category: development
Notes (optional): Personal account

Password:
  ▸ 🎲 Generate secure password
    ✏️  Enter password manually

Password length: 20

🔑 Generated: X#k9$mP2@nL5vB8&qR4w

⠋ Encrypting and saving...
✓ Credential saved! ID: abc123def456

Copy password to clipboard? [y/N]: y
✓ Password copied! (clipboard will clear in 30s)
```

### View Credential

```
  ┌─────────────────────────────────────────────────┐
  │ GitHub                                          │
  ├─────────────────────────────────────────────────┤
  │  ID:          abc123def456                      │
  │  Username:    user@example.com                  │
  │  URL:         https://github.com                │
  │  Category:    development                       │
  │  Password:    ••••••••••••                      │
  └─────────────────────────────────────────────────┘

  Notes: Personal account

Action:
  ▸ 👁️  Show password
    📋 Copy password to clipboard
    📋 Copy username to clipboard
    🔙 Go back
```

### Settings Menu

```
▶ Settings
──────────────────────────────────────────────────────────────

  1. Session Timeout: 5 minutes
  2. Clipboard Timeout: 30 seconds
  3. Default Category: general
  4. Default Password Length: 20
  5. Include Symbols by Default: true
  6. Back to Main Menu

Select option (1-6): _
```

---

## ⚙️ Configuration

Configuration is stored at `~/.passmanager/config.json`:

```json
{
  "pocketbase_url": "http://127.0.0.1:8090",
  "admin_email": "admin@example.com",
  "initialized": true,
  "settings": {
    "session_timeout_minutes": 5,
    "clipboard_timeout_seconds": 30,
    "default_category": "general",
    "password_length": 20,
    "include_symbols": true
  }
}
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `pocketbase_url` | PocketBase server URL | - |
| `admin_email` | Admin email for authentication | - |
| `session_timeout_minutes` | Auto-lock after inactivity | 5 |
| `clipboard_timeout_seconds` | Clear clipboard after copy | 30 |
| `default_category` | Default category for new credentials | "general" |
| `password_length` | Default generated password length | 20 |
| `include_symbols` | Include symbols in generated passwords | true |

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
# Use shorter session timeout for shared computers
# Settings → Session Timeout → 1 minute

# Enable clipboard auto-clear
# Settings → Clipboard Timeout → 15 seconds

# Lock vault when stepping away
# Main Menu → Lock Vault (or Ctrl+C)
```

### PocketBase Hardening

```bash
# Run with HTTPS (production)
./pocketbase serve --https=yourdomain.com

# Restrict to localhost only
./pocketbase serve --http=127.0.0.1:8090

# Regular backups
cp -r pb_data/ pb_data_backup_$(date +%Y%m%d)/
```

### Data Backup

```bash
# Use Export feature in PassManager
# Main Menu → Export Vault

# Backup PocketBase data
cp -r pb_data/ pb_data_backup_$(date +%Y%m%d)/

# Backup config
cp ~/.passmanager/config.json ~/.passmanager/config.json.backup
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. "Configuration not found"

```bash
# First run - setup wizard will start automatically
./passmanager
```

#### 2. "Cannot connect to PocketBase"

```bash
# Check if PocketBase is running
curl http://127.0.0.1:8090/api/health

# Start PocketBase
./pocketbase serve
```

#### 3. "Authentication failed"

```
# Verify credentials in PocketBase Admin UI
# http://127.0.0.1:8090/_/

# For PocketBase 0.23+: Uses _superusers collection
# For PocketBase < 0.23: Uses /api/admins/auth-with-password
```

#### 4. "Collection not found"

```bash
# Create required collections in PocketBase Admin UI
# See "PocketBase Setup" section above
```

#### 5. "Invalid master password"

```
⚠️ Master password cannot be recovered!

# If forgotten, you must reinitialize (data will be lost):
rm -rf ~/.passmanager/
./passmanager
```

#### 6. "Clipboard not working" (Linux)

```bash
# Install xclip
sudo apt install xclip

# Or xsel
sudo apt install xsel

# Or wl-clipboard for Wayland
sudo apt install wl-clipboard
```

#### 7. Session expires too quickly

```
# Increase session timeout in Settings
# Main Menu → Settings → Session Timeout → 15
```

---

## 🛠 Development

### Project Structure

```
passmanager/
├── internal/
│   ├── crypto/
│   │   └── crypto.go         # Encryption, key derivation, password generation
│   ├── database/
│   │   └── pocketbase.go     # PocketBase API client
│   ├── models/
│   │   └── credential.go     # Data models
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── session/
│   │   └── session.go        # Session & authentication state
│   └── ui/
│       ├── colors.go         # ANSI color codes
│       ├── ui.go             # UI helpers, banners, cards
│       └── menu.go           # Interactive menu system
├── main.go                   # Application entry point & handlers
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

### Dependencies

```go
require (
    github.com/atotto/clipboard v0.1.4
    github.com/briandowns/spinner v1.23.0
    github.com/manifoldco/promptui v0.9.0
    golang.org/x/crypto v0.17.0
    golang.org/x/term v0.15.0
)
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
GOOS=darwin GOARCH=arm64 go build -o passmanager-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o passmanager-windows-amd64.exe
```

### Using Makefile

```bash
# Build
make build

# Run
make run

# Install globally
make install

# Build all platforms
make build-all

# Clean
make clean

# Run tests
make test

# Run linter
make lint
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/crypto/...

# Verbose output
go test -v ./...
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

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Anuj Barve

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
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [spinner](https://github.com/briandowns/spinner) - Terminal spinners
- [Go Crypto](https://pkg.go.dev/golang.org/x/crypto) - Cryptographic primitives
- [Argon2](https://github.com/P-H-C/phc-winner-argon2) - Password hashing algorithm

---

## 📊 Roadmap

- [x] **v1.0** - Interactive CLI with session management
- [x] **v1.0** - Edit credentials
- [x] **v1.0** - Export vault (encrypted backup)
- [x] **v1.0** - Change master password
- [x] **v1.0** - Configurable settings
- [x] **v1.0** - Clipboard auto-clear
- [ ] **v1.1** - Import from other password managers
- [ ] **v1.2** - Password strength analyzer
- [ ] **v1.3** - TOTP/2FA support
- [ ] **v1.4** - Password breach checking (HaveIBeenPwned)
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

## Additional Files

### .gitignore

```gitignore
# Binaries
passmanager
passmanager.exe
passmanager-*
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
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# PocketBase
pb_data/
pocketbase
pocketbase.exe
*.zip

# Config (contains sensitive paths)
.passmanager/

# Environment
.env
.env.*
!.env.example

```

### Makefile

```makefile
.PHONY: build clean test lint run install build-all help

BINARY_NAME=passmanager
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_FLAGS=-ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME)
	@echo "✓ Build complete: ./$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Build and run
run: build
	./$(BINARY_NAME)

# Install to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin/$(BINARY_NAME)..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed successfully"

# Cross-compile for all platforms
build-all: clean
	@echo "Building for all platforms..."
	
	@echo "  → Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64
	
	@echo "  → Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64
	
	@echo "  → macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64
	
	@echo "  → macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64
	
	@echo "  → Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe
	
	@echo "✓ All builds complete"
	@ls -la $(BINARY_NAME)-*

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod verify
	@echo "✓ Dependencies updated"

# Show help
help:
	@echo "PassManager Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build      Build for current platform"
	@echo "  clean      Remove build artifacts"
	@echo "  test       Run tests"
	@echo "  coverage   Run tests with coverage report"
	@echo "  lint       Run golangci-lint"
	@echo "  run        Build and run"
	@echo "  install    Install to /usr/local/bin"
	@echo "  build-all  Cross-compile for all platforms"
	@echo "  deps       Update dependencies"
	@echo "  help       Show this help"
```

### go.mod

```go
module passmanager

go 1.21

require (
	github.com/atotto/clipboard v0.1.4
	github.com/briandowns/spinner v1.23.0
	github.com/manifoldco/promptui v0.9.0
	golang.org/x/crypto v0.17.0
	golang.org/x/term v0.15.0
)

require (
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
```