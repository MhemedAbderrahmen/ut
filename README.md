# UploadThing CLI (ut)

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A powerful command-line interface for [UploadThing](https://uploadthing.com), enabling seamless file uploads and downloads directly from your terminal.

## Features

- Fast file uploads** to UploadThing with progress tracking
- Download files** with custom paths and progress indication
- **Private file support** with API key authentication
- **List uploaded files** with detailed metadata
- **Easy configuration** management
- **Progress tracking** for large file operations
- **Flexible output options** for downloads
- **Force overwrite** capabilities

## Installation

### Prerequisites

- Go 1.24 or higher
- UploadThing account and API key

### From Source

```bash
git clone https://github.com/username/ut.git
cd uploadthing-cli
go build -o ut .
```


## Quick Start

### 1. Configure Your API Key

```bash
ut config
```

### 2. Upload a File

```bash
ut push document.pdf
```

### 3. List Your Files

```bash
ut list
```

### 4. Download a File

```bash
ut fetch your-file-key.pdf
```

## Usage

### Configuration

Configure your UploadThing secret key:

```bash
# Set your secret key
ut config set-secret sk_your_secret_key_here

# View current configuration
ut config show

```

### File Upload

Upload files to UploadThing:

```bash
# Basic upload
ut push image.jpg

# Upload with progress tracking (automatic for large files)
ut push large-video.mp4
```

**Supported file types:** Images (JPG, PNG, GIF), Documents (PDF, TXT, JSON, XML, CSV), and more.

### File Download

Download files from UploadThing:

```bash
# Download to current directory
ut fetch abc123-example.jpg

# Download with custom filename
ut fetch abc123-example.jpg -o myfile.jpg

# Download to specific directory
ut fetch abc123-example.jpg -o ./downloads/

# Download with progress bar
ut fetch abc123-example.jpg --progress

# Download private file (requires API key)
ut fetch abc123-example.jpg --private

# Force overwrite existing files
ut fetch abc123-example.jpg --force
```

### List Files

View your uploaded files:

```bash
# List all files
ut list

# List with file details
ut list --verbose
```

## Configuration

The CLI stores configuration in `~/.ut-cli/config.yml` by default. You can customize the location:

## 📋 Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `ut config` | Set your UploadThing secret key | `ut config` |
| `ut push <file>` | Upload a file to UploadThing | `ut push document.pdf` |
| `ut fetch <filekey>` | Download a file by file key | `ut fetch abc123-file.jpg` |
| `ut list` | List all uploaded files | `ut list` |

### Command Options

#### `ut fetch` options:
- `-o, --output`: Custom output path or directory
- `-f, --force`: Overwrite existing files without prompt
- `-p, --progress`: Show download progress
- `--private`: Download private file (requires API key)

#### `ut list` options:
- `-v, --verbose`: Show detailed file information

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/ut.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
## Acknowledgments - [UploadThing](https://uploadthing.com) for their excellent file upload service [Cobra](https://github.com/spf13/cobra) for the CLI framework
- All contributors who help improve this tool

## Support

-  **Bug Reports**: [GitHub Issues](https://github.com/MhemedAbderrahmen/ut/issues)
-  **Feature Requests**: [GitHub Discussions](https://github.com/MhemedAbderrahmen/ut/discussions)
-  **Documentation**: [Wiki](https://github.com/MhemedAbderrahmen/ut/wiki)

## Links
- [UploadThing Documentation](https://docs.uploadthing.com) [UploadThing Dashboard](https://uploadthing.com/dashboard)
- [API Reference](https://docs.uploadthing.com/api-reference)

---

**Made with ❤️ for the UploadThing community** 
