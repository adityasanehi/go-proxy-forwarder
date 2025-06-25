# Deployment Guide

This document explains how to deploy the Go Proxy Rotator using GitHub Actions for automated releases.

## GitHub Actions Setup

The project includes two GitHub Actions workflows:

### 1. Build and Test Workflow (`.github/workflows/build.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` branch

**Actions:**
- Runs tests and code quality checks
- Builds binaries for multiple platforms
- Uploads build artifacts for testing

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of a git tag (e.g., `v1.0.0`)
- Manual workflow dispatch

**Actions:**
- Builds optimized binaries for all supported platforms
- Creates release packages with all necessary files
- Automatically creates GitHub release with download links

## Creating a Release

### Method 1: Git Tags (Recommended)

1. **Ensure your code is ready:**
   ```bash
   git add .
   git commit -m "Release v1.0.0"
   git push origin main
   ```

2. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all release assets

### Method 2: Manual Workflow Dispatch

1. **Go to your GitHub repository**
2. **Navigate to Actions tab**
3. **Select "Release" workflow**
4. **Click "Run workflow"**
5. **Enter the desired tag name (e.g., v1.0.0)**
6. **Click "Run workflow"**

### Method 3: GitHub CLI

```bash
# Create release with auto-generated notes
gh release create v1.0.0 --generate-notes

# Create release with custom notes
gh release create v1.0.0 --notes "Release notes here"
```

## Supported Platforms

The GitHub Actions build the following binaries:

| Platform | Architecture | Binary Name |
|----------|-------------|-------------|
| Linux | AMD64 | `go-proxy-rotator-linux-amd64` |
| Linux | ARM64 | `go-proxy-rotator-linux-arm64` |
| Windows | AMD64 | `go-proxy-rotator-windows-amd64.exe` |
| macOS | Intel | `go-proxy-rotator-darwin-amd64` |
| macOS | Apple Silicon | `go-proxy-rotator-darwin-arm64` |

## Release Assets

Each release includes:

1. **Compressed archives** containing:
   - Pre-compiled binary
   - Static web interface files (`static/` directory)
   - Sample proxy list (`sample_proxies.txt`)
   - Documentation (`README.md`)

2. **Archive formats:**
   - Linux/macOS: `.tar.gz`
   - Windows: `.zip`

## Local Building

For local development and testing:

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release packages
make package

# Run tests
make test
```

## Environment Variables

Configure the application using environment variables:

```bash
# Server configuration
export PORT=3000
export DATABASE_PATH=./proxies.db
export MAX_FILE_SIZE=10485760  # 10MB

# Health check configuration
export HEALTH_CHECK_URL=https://httpbin.org/ip
export LOG_LEVEL=info
```

## Docker Deployment

Alternative deployment using Docker:

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or build manually
docker build -t go-proxy-rotator .
docker run -p 3000:3000 -v proxy_data:/data go-proxy-rotator
```

## Security Considerations

1. **Database Security:**
   - Store the SQLite database in a secure location
   - Regular backups of the database file
   - Consider encryption for sensitive proxy credentials

2. **Network Security:**
   - Use HTTPS in production environments
   - Consider running behind a reverse proxy (nginx, Apache)
   - Implement rate limiting if needed

3. **Access Control:**
   - The web interface has no built-in authentication
   - Consider adding authentication middleware for production use
   - Restrict access using firewall rules or VPN

## Monitoring

Monitor the application using:

1. **Health endpoint:** `GET /health`
2. **Statistics endpoint:** `GET /api/v1/proxies/stats`
3. **Application logs** (stdout/stderr)
4. **Database file size and integrity**

## Troubleshooting

### Common Issues

1. **Build failures:**
   - Check Go version compatibility (requires Go 1.19+)
   - Ensure all dependencies are available
   - Check for syntax errors in code

2. **Release failures:**
   - Verify GitHub Actions permissions
   - Check if tag already exists
   - Ensure repository has proper access tokens

3. **Runtime issues:**
   - Check database file permissions
   - Verify static files are included in release
   - Check environment variable configuration

### Debug Commands

```bash
# Check build locally
go build -v .

# Run with debug logging
LOG_LEVEL=debug ./go-proxy-rotator

# Test database connectivity
sqlite3 proxies.db ".tables"
```

## Version Management

The application uses semantic versioning (SemVer):

- **Major version** (v1.0.0 → v2.0.0): Breaking changes
- **Minor version** (v1.0.0 → v1.1.0): New features, backward compatible
- **Patch version** (v1.0.0 → v1.0.1): Bug fixes, backward compatible

Build information is embedded in the binary and accessible via the `/health` endpoint.