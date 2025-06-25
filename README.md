# Go Proxy Forwarder

An advanced proxy forwarder built with Go and Fiber. It acts as a proxy server, forwarding your requests through randomly selected proxies from a persistent database with health monitoring and management features.

## Features

### Core Functionality
- 🔄 **Smart Proxy Rotation**: Automatically rotates through healthy proxies with intelligent selection
- 📊 **Health Monitoring**: Real-time tracking of proxy response times, failure rates, and availability
- 💾 **Persistent Storage**: SQLite database with proper indexing for reliable proxy management
- 🏥 **Automated Health Checks**: Continuous monitoring with configurable test URLs and intervals

### User Interface & Experience
- 🌐 **Modern Web Interface**: Professional dashboard with responsive design and real-time updates
- 📤 **Drag & Drop Upload**: Intuitive file upload with progress indicators and validation
- 📈 **Live Statistics**: Real-time metrics with auto-refresh and visual indicators
- 🔔 **Toast Notifications**: Instant feedback for all user actions with auto-dismiss
- 💡 **Tooltips & Shortcuts**: Helpful hints and keyboard shortcuts (Ctrl+R, Ctrl+U)
- 📱 **Mobile Responsive**: Fully optimized for mobile devices and tablets

### API & Integration
- 🔧 **Comprehensive REST API**: Complete CRUD operations with proper HTTP status codes
- 📝 **Multiple Input Formats**: Support for various proxy list formats and URL schemes
- 🔐 **Authentication Support**: HTTP Basic Auth and SOCKS5 proxy authentication
- 📋 **Detailed Documentation**: Complete API documentation with examples and SDKs

### Advanced Features
- ⚡ **Performance Optimization**: Response time-based selection and failure tracking
- 🎯 **Load Balancing**: Intelligent proxy selection based on health and performance
- 🔍 **Advanced Filtering**: Filter proxies by status, protocol, performance, and more
- 📊 **Analytics**: Detailed statistics with success rates and performance metrics

## Quick Start

### Option 1: Download Pre-built Binary (Recommended)

1. **Download the latest release:**
   - Go to the [Releases page](../../releases)
   - Download the appropriate binary for your platform:
     - Linux: `go-proxy-forwarder-linux-amd64.tar.gz`
     - Windows: `go-proxy-forwarder-windows-amd64.exe.zip`
     - macOS Intel: `go-proxy-forwarder-darwin-amd64.tar.gz`
     - macOS Apple Silicon: `go-proxy-forwarder-darwin-arm64.tar.gz`

2. **Extract and run:**
   ```bash
   # Linux/macOS
   tar -xzf go-proxy-forwarder-*.tar.gz
   cd go-proxy-forwarder-*
   ./go-proxy-forwarder
   
   # Windows
   # Extract the ZIP file and run go-proxy-forwarder.exe
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
   cd go-proxy-forwarder
   ```

2. **Build and run:**
   ```bash
   # Using Make (recommended)
   make run
   
   # Or manually
   go mod download
   go build -o go-proxy-forwarder .
   ./go-proxy-forwarder
   ```

3. **Access the web interface:**
   Open `http://localhost:3000` in your browser

## Usage

### Web Interface

The web interface provides:
- 📊 Real-time proxy statistics
- 📤 Drag-and-drop file upload for proxy lists
- 📋 Proxy list management with health status
- 🔄 Manual health checks and proxy refresh
- 🗑️ Individual and bulk proxy deletion

### Proxy Usage

Point your HTTP client to `http://localhost:3000` to use the proxy forwarder:

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

## API Endpoints

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

## API Documentation

### Interactive Documentation
- **Swagger UI**: `http://localhost:3000/docs` - Interactive API documentation with try-it-out functionality
- **OpenAPI Spec**: `http://localhost:3000/docs/swagger.yaml` - Machine-readable API specification

### Static Documentation
- **Comprehensive Guide**: [API Docs](docs/API.md) - Detailed API documentation with examples

## Configuration

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

## Development & Releases

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

## 🗺️ Roadmap

### Phase 1: Core Stability ✅ 
- [x] Basic proxy rotation with health monitoring
- [x] SQLite database with proper schema
- [x] REST API with CRUD operations
- [x] File upload functionality
- [x] Basic web interface

### Phase 2: Enhanced UX ✅ 
- [x] Modern, responsive web interface
- [x] Real-time statistics and auto-refresh
- [x] Drag & drop file uploads
- [x] Toast notifications and user feedback
- [x] Mobile-optimized design
- [x] Comprehensive API documentation

### Phase 3: Advanced Features 🚧
- [ ] **Geo-location Support**: Country/city detection for proxies
- [ ] **Advanced Analytics**: Historical performance graphs and trends
- [ ] **Notification System**: Email/webhook alerts for proxy failures
- [ ] **Load Balancing**: Weighted round-robin and least connections
- [ ] **Proxy Chaining**: Multi-hop proxy support
- [ ] **Custom Health Checks**: Configurable test URLs and intervals

### Phase 4: Enterprise Features 🔮 
- [ ] **User Authentication**: Login system with role-based access
- [ ] **Multi-tenancy**: Support for multiple users/organizations
- [ ] **Audit Logging**: Comprehensive activity logs
- [ ] **High Availability**: Clustering and failover support
- [ ] **Monitoring Integration**: Prometheus/Grafana metrics
- [ ] **API Rate Limiting**: Request throttling and quotas
- [ ] **Docker Swarm Support**: Native swarm deployment
- [ ] **IPv6 Support**: Full IPv6 proxy compatibility
- [ ] **SOCKS4 Support**: Additional protocol support
- [ ] **Proxy Rotation Algorithms**: Custom rotation strategies
- [ ] **Backup/Restore**: Database backup and migration tools
- [ ] **Performance Benchmarking**: Built-in speed testing

## 📁 Project Structure

```
go-proxy-forwarder/
├── .github/workflows/   # GitHub Actions CI/CD
├── config/             # Configuration management
├── database/           # Database operations
├── docs/               # API documentation
│   ├── swagger.yaml    # OpenAPI 3.0 specification
|   ├── API.md # Comprehensive API guide
|   ├── DEPLOYMENT.md # Comprehensive Deployment guide
|   └── UI.md # Comprehensive UI updates
├── handlers/           # HTTP handlers
│   ├── proxy_handlers.go
│   └── swagger_handler.go
├── models/             # Data models
├── services/           # Business logic
├── static/             # Web interface files
├── tools/              # Development tools
│   └── swagger-gen.go  # OpenAPI validation/conversion
├── main.go             # Application entry point
├── Makefile            # Build automation
├── Dockerfile          # Docker configuration
├── docker-compose.yml  # Docker Compose setup
└── README.md           # Project documentation
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
