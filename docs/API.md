# Go Proxy Rotator API Documentation

## Overview

The Go Proxy Rotator provides a comprehensive REST API for managing proxies, monitoring health, and accessing statistics. All API endpoints return JSON responses and follow RESTful conventions.

**Base URL**: `http://localhost:3000/api/v1`

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "data": {...},
  "message": "Success message",
  "count": 10
}
```

### Error Response
```json
{
  "error": "Error description"
}
```

## HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - No proxies available

---

## Endpoints

### Proxy Management

#### Upload Proxy List

Upload a text file containing proxy configurations.

**Endpoint**: `POST /api/v1/proxies/upload`

**Content-Type**: `multipart/form-data`

**Parameters**:
- `file` (file, required) - Text file containing proxy list

**Supported File Formats**:
```
# Format 1: host:port:username:password
192.168.1.100:8080:user1:pass1
192.168.1.101:8080:user2:pass2

# Format 2: host:port (no authentication)
203.0.113.1:3128
203.0.113.2:3128

# Format 3: URL format
http://user:pass@proxy.example.com:8080
socks5://user:pass@proxy.example.com:1080
```

**Example Request**:
```bash
curl -X POST \
  -F "file=@proxies.txt" \
  http://localhost:3000/api/v1/proxies/upload
```

**Example Response**:
```json
{
  "message": "Proxies uploaded successfully",
  "total_parsed": 50,
  "total_added": 45,
  "total_skipped": 5
}
```

---

#### Get All Proxies

Retrieve all proxies from the database.

**Endpoint**: `GET /api/v1/proxies`

**Example Request**:
```bash
curl http://localhost:3000/api/v1/proxies
```

**Example Response**:
```json
{
  "proxies": [
    {
      "id": 1,
      "host": "192.168.1.100",
      "port": 8080,
      "username": "user1",
      "password": "pass1",
      "protocol": "http",
      "is_active": true,
      "last_checked": "2024-01-15T10:30:00Z",
      "response_time": 250,
      "fail_count": 0,
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

#### Get Active Proxies

Retrieve only active and healthy proxies.

**Endpoint**: `GET /api/v1/proxies/active`

**Example Request**:
```bash
curl http://localhost:3000/api/v1/proxies/active
```

**Example Response**:
```json
{
  "proxies": [
    {
      "id": 1,
      "host": "192.168.1.100",
      "port": 8080,
      "username": "user1",
      "password": "pass1",
      "protocol": "http",
      "is_active": true,
      "last_checked": "2024-01-15T10:30:00Z",
      "response_time": 250,
      "fail_count": 0,
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

#### Add Single Proxy

Add a single proxy to the database.

**Endpoint**: `POST /api/v1/proxies`

**Content-Type**: `application/json`

**Request Body**:
```json
{
  "host": "192.168.1.100",
  "port": 8080,
  "username": "user1",
  "password": "pass1",
  "protocol": "http"
}
```

**Required Fields**:
- `host` (string) - Proxy host/IP address
- `port` (integer) - Proxy port number

**Optional Fields**:
- `username` (string) - Authentication username
- `password` (string) - Authentication password
- `protocol` (string) - Protocol type (http, https, socks5). Default: "http"

**Example Request**:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "host": "192.168.1.100",
    "port": 8080,
    "username": "user1",
    "password": "pass1",
    "protocol": "http"
  }' \
  http://localhost:3000/api/v1/proxies
```

**Example Response**:
```json
{
  "message": "Proxy added successfully",
  "proxy": {
    "id": 1,
    "host": "192.168.1.100",
    "port": 8080,
    "username": "user1",
    "password": "pass1",
    "protocol": "http",
    "is_active": true,
    "last_checked": "2024-01-15T10:30:00Z",
    "response_time": 0,
    "fail_count": 0,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

#### Delete Proxy

Delete a specific proxy by ID.

**Endpoint**: `DELETE /api/v1/proxies/{id}`

**Parameters**:
- `id` (integer, required) - Proxy ID

**Example Request**:
```bash
curl -X DELETE http://localhost:3000/api/v1/proxies/1
```

**Example Response**:
```json
{
  "message": "Proxy deleted successfully"
}
```

---

#### Clear All Proxies

Remove all proxies from the database.

**Endpoint**: `DELETE /api/v1/proxies`

**Example Request**:
```bash
curl -X DELETE http://localhost:3000/api/v1/proxies
```

**Example Response**:
```json
{
  "message": "All proxies cleared successfully"
}
```

---

### Statistics

#### Get Proxy Statistics

Retrieve comprehensive proxy statistics.

**Endpoint**: `GET /api/v1/proxies/stats`

**Example Request**:
```bash
curl http://localhost:3000/api/v1/proxies/stats
```

**Example Response**:
```json
{
  "total_proxies": 100,
  "active_proxies": 85,
  "healthy_proxies": 75,
  "failed_proxies": 15
}
```

**Response Fields**:
- `total_proxies` - Total number of proxies in database
- `active_proxies` - Number of active proxies (is_active = true)
- `healthy_proxies` - Number of healthy proxies (active + fail_count < 5 + response_time < 10s)
- `failed_proxies` - Number of failed proxies (inactive or fail_count >= 5)

---

### Health Monitoring

#### Run Health Check

Perform health checks on all active proxies.

**Endpoint**: `POST /api/v1/proxies/health-check`

**Query Parameters**:
- `url` (string, optional) - Custom test URL. Default: "https://httpbin.org/ip"

**Example Request**:
```bash
curl -X POST http://localhost:3000/api/v1/proxies/health-check
```

**Example Request with Custom URL**:
```bash
curl -X POST "http://localhost:3000/api/v1/proxies/health-check?url=https://example.com"
```

**Example Response**:
```json
{
  "message": "Health check completed"
}
```

**Health Check Process**:
1. Tests connectivity to each active proxy
2. Measures response time
3. Updates proxy statistics
4. Deactivates proxies with 5+ consecutive failures
5. Records last check timestamp

---

### Application Health

#### Application Health Status

Check the overall health and status of the application.

**Endpoint**: `GET /health`

**Example Request**:
```bash
curl http://localhost:3000/health
```

**Example Response**:
```json
{
  "status": "ok",
  "time": "2024-01-15T10:30:00Z",
  "version": "v1.0.0",
  "build_time": "2024-01-15_09:00:00",
  "git_commit": "abc123f"
}
```

---

## Proxy Usage

### Using the Proxy Rotator

The application acts as a proxy server that automatically rotates through healthy proxies.

**Proxy Endpoint**: `http://localhost:3000`

**Example Usage**:
```bash
# Use as HTTP proxy
curl -x http://localhost:3000 https://httpbin.org/ip

# Use with wget
wget -e use_proxy=yes -e http_proxy=localhost:3000 https://httpbin.org/ip

# Use with Python requests
import requests

proxies = {
    'http': 'http://localhost:3000',
    'https': 'http://localhost:3000'
}

response = requests.get('https://httpbin.org/ip', proxies=proxies)
print(response.json())
```

**Proxy Selection Logic**:
1. Filters for active proxies (is_active = true)
2. Prioritizes healthy proxies (fail_count < 5, response_time < 10s)
3. Randomly selects from available healthy proxies
4. Falls back to any active proxy if no healthy proxies available
5. Updates proxy health statistics after each use

---

## Error Handling

### Common Error Responses

#### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

#### 404 Not Found
```json
{
  "error": "Proxy with id 123 not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to connect to database"
}
```

#### 503 Service Unavailable
```json
{
  "error": "No proxies available"
}
```

---

## Rate Limiting

Currently, there are no rate limits implemented. Consider implementing rate limiting for production use.

---

## Examples

### Complete Workflow Example

```bash
# 1. Check application health
curl http://localhost:3000/health

# 2. Upload proxy list
curl -X POST -F "file=@my_proxies.txt" http://localhost:3000/api/v1/proxies/upload

# 3. Check statistics
curl http://localhost:3000/api/v1/proxies/stats

# 4. Run health check
curl -X POST http://localhost:3000/api/v1/proxies/health-check

# 5. Get active proxies
curl http://localhost:3000/api/v1/proxies/active

# 6. Use proxy rotator
curl -x http://localhost:3000 https://httpbin.org/ip

# 7. Add individual proxy
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"host":"1.2.3.4","port":8080,"username":"user","password":"pass"}' \
  http://localhost:3000/api/v1/proxies

# 8. Delete specific proxy
curl -X DELETE http://localhost:3000/api/v1/proxies/1
```

### Python SDK Example

```python
import requests
import json

class ProxyRotatorAPI:
    def __init__(self, base_url="http://localhost:3000"):
        self.base_url = base_url
        self.api_url = f"{base_url}/api/v1"
    
    def upload_proxies(self, file_path):
        with open(file_path, 'rb') as f:
            files = {'file': f}
            response = requests.post(f"{self.api_url}/proxies/upload", files=files)
        return response.json()
    
    def get_stats(self):
        response = requests.get(f"{self.api_url}/proxies/stats")
        return response.json()
    
    def health_check(self):
        response = requests.post(f"{self.api_url}/proxies/health-check")
        return response.json()
    
    def add_proxy(self, host, port, username=None, password=None, protocol="http"):
        data = {
            "host": host,
            "port": port,
            "protocol": protocol
        }
        if username:
            data["username"] = username
        if password:
            data["password"] = password
        
        response = requests.post(
            f"{self.api_url}/proxies",
            headers={"Content-Type": "application/json"},
            data=json.dumps(data)
        )
        return response.json()
    
    def use_proxy(self, url):
        proxies = {
            'http': self.base_url,
            'https': self.base_url
        }
        response = requests.get(url, proxies=proxies)
        return response

# Usage
api = ProxyRotatorAPI()
stats = api.get_stats()
print(f"Total proxies: {stats['total_proxies']}")

# Use proxy
response = api.use_proxy("https://httpbin.org/ip")
print(response.json())
```

---

## WebSocket Support

Currently, WebSocket support is not implemented. Consider adding WebSocket endpoints for real-time updates in future versions.

---

## Versioning

The API uses URL versioning with the `/api/v1` prefix. Future versions will use `/api/v2`, etc.

---

## Support

For issues and questions:
- GitHub Issues: [Repository Issues](https://github.com/your-repo/go-proxy-rotator/issues)
- Documentation: [README.md](README.md)
- Health Check: `GET /health`