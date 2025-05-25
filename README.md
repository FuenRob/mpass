# 🔐 MPASS (Manager Passwords)

A secure command-line password manager written in Go. Mpass stores your passwords encrypted locally and provides quick access without displaying passwords on screen.

## ✨ Features

- **🛡️ Robust encryption**: AES-256-GCM with PBKDF2 key derivation (100,000 iterations)
- **📋 Automatic clipboard**: Passwords are copied directly to clipboard
- **🔍 Smart search**: Search by username, URL, or both
- **🎯 Multiple selector**: Elegant handling of multiple matches
- **💾 Database-free**: Encrypted local file storage
- **🌐 Cross-platform**: Compatible with Windows, macOS, and Linux
- **🚫 No traces**: Passwords never displayed on screen
- **🔒 Local security**: Unique salt + master password

## 🚀 Installation

### Prerequisites

- Go 1.21 or higher
- For Linux: `xclip` or `xsel` (for clipboard functionality)

```bash
# Ubuntu/Debian
sudo apt install xclip

# Fedora/RHEL
sudo dnf install xclip

# Arch Linux
sudo pacman -S xclip
```

### Build

```bash
# Clone the repository
git clone <repository-url>
cd mpass

# Install dependencies
go mod tidy

# Build
go build -o mpass

# (Optional) Install globally
sudo mv mpass /usr/local/bin/
```

## 📖 Usage

### Available commands

| Command | Description |
|---------|-------------|
| `add` | Add a new password |
| `get -u <username>` | Search by username |
| `get -l <url>` | Search by URL |
| `get -u <username> -l <url>` | Search by username AND URL |
| `list` | List all entries (without showing passwords) |

### Usage examples

#### ➕ Add a password

```bash
$ ./mpass add
Enter master password: ********
Username: rob@example.com
URL: https://github.com
Password: ********
✅ Password entry added successfully!
```

#### 🔍 Search by username

```bash
$ ./mpass get -u rob
Enter master password: ********
✅ Password for rob@github.com copied to clipboard!
```

#### 🌐 Search by URL

```bash
$ ./mpass get -l github.com
Enter master password: ********
✅ Password for rob@github.com copied to clipboard!
```

#### 📋 List all entries

```bash
$ ./mpass list
Enter master password: ********
📚 Found 3 password entries:

1. rob@github.com
2. alice@gitlab.com  
3. admin@company.com
```

#### 🎯 Multiple matches

```bash
$ ./mpass get -u rob
Enter master password: ********

🔍 Multiple entries found:
────────────────────────────
1. rob@github.com
2. rob@gitlab.com

Select entry number (1-2): 1
✅ Password for rob@github.com copied to clipboard!
```

## 🏗️ Architecture

```
mpass/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command
│   ├── add.go             # Add command
│   ├── get.go             # Get command
│   └── list.go            # List command
├── internal/              # Internal code
│   ├── crypto/            # Encryption functions
│   ├── storage/           # Vault management
│   ├── models/            # Data structures
│   └── ui/                # User interface
├── pkg/                   # Public packages
│   └── clipboard/         # Clipboard utilities
├── go.mod
├── go.sum
└── main.go
```

## 🔒 Security

### Encryption

- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key derivation**: PBKDF2 with SHA-256
- **Iterations**: 100,000
- **Salt**: 256-bit random unique per vault
- **Nonce**: Randomly generated for each operation

### Storage

- **Location**: `~/.mpass/vault.enc`
- **Permissions**: 600 (read/write for owner only)
- **Format**: Salt (32 bytes) + Encrypted data

### Security flow

1. **First time**: A unique random salt is generated
2. **Each operation**:
    - Master password + salt → PBKDF2 → Encryption key
    - Automatic integrity verification with GCM
    - Fails on incorrect password or corrupted data

## 🛠️ Development

### Code structure

- **Modular**: Clear separation of responsibilities
- **Testable**: Well-defined interfaces
- **Scalable**: Easy to add new functionality
- **Maintainable**: Clean and documented code

### Dependencies

```go
require (
    github.com/spf13/cobra v1.8.0    // CLI framework
    golang.org/x/crypto v0.17.0      // Cryptographic functions
    golang.org/x/term v0.15.0        // Password input
)
```

### Testing

```bash
# Run tests
go test ./...

# Tests with coverage
go test -cover ./...
```

## 🚨 Security considerations

### ✅ Best practices

- **Strong master password**: Use a unique and complex password
- **Secure backup**: Back up your vault in a secure location
- **Memory cleanup**: Passwords are cleared from memory after use
- **No logging**: Passwords are not recorded in system logs

### ⚠️ Limitations

- **Single-user**: Designed for individual use
- **No synchronization**: Local storage only
- **Clipboard dependent**: Requires system clipboard access

## 🐛 Troubleshooting

### Error: "failed to copy to clipboard"

**Linux:**
```bash
# Install xclip or xsel
sudo apt install xclip
# or
sudo apt install xsel
```

**WSL (Windows Subsystem for Linux):**
```bash
# Install clip.exe
echo "alias clip='clip.exe'" >> ~/.bashrc
source ~/.bashrc
```

### Error: "failed to decrypt vault (wrong password?)"

- Verify you're using the correct master password
- The vault file might be corrupted - restore from backup

### Error: "permission denied"

```bash
# Check directory permissions
ls -la ~/.mpass/
# Fix permissions if necessary
chmod 700 ~/.mpass/
chmod 600 ~/.mpass/vault.enc
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 Roadmap

- [ ] Integrated password generator
- [ ] Export/import backup
- [ ] Categories and tags
- [ ] Cloud synchronization (encrypted)
- [ ] Optional web or desktop interface
- [ ] Browser extension
- [ ] Mobile application

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## 👨‍💻 Author

Developed by Roberto Morais with ❤️ using Go

---

**⚠️ Important**: This software is provided "as is". Always maintain secure backups of your important data.