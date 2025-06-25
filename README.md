# Go Proxy Forwarder

An advanced proxy rotator built with Go and Fiber. It acts as a proxy server, forwarding your requests through randomly selected proxies from a persistent database with health monitoring and management features.

## Features

- 🔄 **Smart Proxy Rotation**: Automatically rotates through healthy proxies
- 📊 **Health Monitoring**: Tracks proxy response times and failure rates
- 📤 **File Upload**: Upload proxy lists via web interface or API
- 💾 **Persistent Storage**: SQLite database for proxy management
- 🌐 **Web Interface**: User-friendly dashboard for proxy management
- 🔧 **REST API**: Complete API for proxy CRUD operations
- 📈 **Statistics**: Real-time proxy statistics and health metrics
- 🏥 **Health Checks**: Automated proxy health verification
- 🔐 **Authentication Support**: HTTP and SOCKS5 proxy authentication
- 📝 **Multiple Formats**: Support for various proxy list formats

## 🚀 Quick Start

### Option 1: Download Pre-built Binary (Recommended)

1. **Download the latest release:**
   - Go to the [Releases page](../../releases)
   - Download the appropriate binary for your platform:
     - Linux: `go-proxy-rotator-linux-amd64.tar.gz`
     - Windows: `go-proxy-rotator-windows-amd64.exe.zip`
     - macOS Intel: `go-proxy-rotator-darwin-amd64.tar.gz`
     - macOS Apple Silicon: `go-proxy-rotator-darwin-arm64.tar.gz`

2. **Extract and run:**
   ```bash
   # Linux/macOS
   tar -xzf go-proxy-rotator-*.tar.gz
   cd go-proxy-rotator-*
   ./go-proxy-rotator
   
   # Windows
   # Extract the ZIP file and run go-proxy-rotator.exe
   ```

3. **Access the web interface:**
   Open `http://localhost:3000` in your browser

### Option 2: Build from Source

#### Prerequisites
- [Go](https://golang.org/doc/install) 1.19+ installed on your machine

#### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd go-proxy-rotator
   ```

2. **Build and run:**
   ```bash
   # Using Make (recommended)
   make run
   
   # Or manually
   go mod download
   go build -o go-proxy-rotator .
   ./go-proxy-rotator
   ```

3. **Access the web interface:**
   Open `http://localhost:3000` in your browser

## 📖 Usage

### Web Interface

The web interface provides:
- 📊 Real-time proxy statistics
- 📤 Drag-and-drop file upload for proxy lists
- 📋 Proxy list management with health status
- 🔄 Manual health checks and proxy refresh
- 🗑️ Individual and bulk proxy deletion

### Proxy Usage

Point your HTTP client to `http://localhost:3000` to use the proxy rotator:

```bash
curl -x http://localhost:3000 https://httpbin.org/ip
```

Each request will be forwarded through a randomly selected healthy proxy.

### Supported Proxy Formats

The application supports multiple proxy list formats:

1. **Host:Port:Username:Password**
   ```
   192.168.1.100:8080:user1:pass1
   192.168.1.101:8080:user2:pass2
   ```

2. **Host:Port** (no authentication)
   ```
   203.0.113.1:3128
   203.0.113.2:3128
   ```

3. **URL Format**
   ```
   http://user:pass@proxy.example.com:8080
   socks5://user:pass@proxy.example.com:1080
   ```

## 🔧 API Endpoints

### Proxy Management

- `POST /api/v1/proxies/upload` - Upload proxy list file
- `GET /api/v1/proxies` - Get all proxies
- `GET /api/v1/proxies/active` - Get active proxies only
- `POST /api/v1/proxies` - Add single proxy
- `DELETE /api/v1/proxies/:id` - Delete specific proxy
- `DELETE /api/v1/proxies` - Clear all proxies
- `GET /api/v1/proxies/stats` - Get proxy statistics
- `POST /api/v1/proxies/health-check` - Run health check

### Health Check

- `GET /health` - Application health status

### Example API Usage

**Upload proxy list:**
```bash
curl -X POST -F "file=@proxies.txt" http://localhost:3000/api/v1/proxies/upload
```

**Add single proxy:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"host":"192.168.1.100","port":8080,"username":"user","password":"pass"}' \
  http://localhost:3000/api/v1/proxies
```

**Get statistics:**
```bash
curl http://localhost:3000/api/v1/proxies/stats
```

## ⚙️ Configuration

Environment variables:

- `PORT` - Server port (default: 3000)
- `DATABASE_PATH` - SQLite database path (default: ./proxies.db)
- `MAX_FILE_SIZE` - Maximum upload file size in bytes (default: 10MB)
- `HEALTH_CHECK_URL` - URL for proxy health checks (default: https://httpbin.org/ip)
- `LOG_LEVEL` - Logging level (default: info)

## 🐳 Docker

You can also run the server using Docker:

```bash
# Build and run
docker-compose up --build

# Run in background
docker-compose up -d

# Stop
docker-compose down
```

The server will be accessible at `http://localhost:3000`.

## 🚀 Development & Releases

### Building Locally

Use the provided Makefile for common development tasks:

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run with development settings
make dev

# Create release packages
make package

# Clean build artifacts
make clean
```

### GitHub Actions

This project uses GitHub Actions for automated building and releasing:

#### Automatic Builds
- **Trigger**: Push to `main` or `develop` branches, or pull requests to `main`
- **Actions**: 
  - Run tests and linting
  - Build binaries for multiple platforms
  - Upload build artifacts

#### Releases
- **Trigger**: Push a git tag (e.g., `v1.0.0`) or manual workflow dispatch
- **Actions**:
  - Build binaries for all supported platforms
  - Create release packages with static files
  - Automatically create GitHub release with binaries

#### Creating a Release

1. **Tag and push:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Or use GitHub CLI:**
   ```bash
   gh release create v1.0.0 --generate-notes
   ```

3. **Manual release via GitHub Actions:**
   - Go to Actions tab in GitHub
   - Select "Release" workflow
   - Click "Run workflow"
   - Enter the desired tag name

#### Supported Platforms

The GitHub Actions automatically build for:
- Linux (AMD64, ARM64)
- Windows (AMD64)
- macOS (Intel, Apple Silicon)

Each release includes:
- Pre-compiled binaries
- Static web interface files
- Sample proxy list
- Documentation

## 📁 Project Structure

```
go-proxy-rotator/
├── config/          # Configuration management
├── database/        # Database operations
├── handlers/        # HTTP handlers
├── models/          # Data models
├── services/        # Business logic
├── static/          # Web interface files
├── main.go          # Application entry point
├── sample_proxies.txt # Example proxy list
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
