o# IP Location API

**High-performance REST API for IP geolocation** - Lightning-fast IP address lookup service using binary search over 2.9M+ location records.

## 📋 Description

Backend challenge implementation that resolves IPv4 addresses to geographic locations using the IP2Location dataset. Built with Go for maximum performance, handling 1,800+ requests/second with sub-millisecond latency.

## 🚀 Demo

- **API Endpoint**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`
- **Swagger Documentation**: `http://localhost:8080/swagger/index.html`

## ✨ Features

- ✅ **Ultra-fast lookups** - Binary search on 2.9M+ records
- ✅ **Sub-millisecond response time** - p95 latency: 0.37ms
- ✅ **High throughput** - 1,833 requests/second
- ✅ **Concurrent load** - Handles 1,000+ simultaneous users
- ✅ **RESTful API** with JSON responses
- ✅ **Swagger/OpenAPI documentation**
- ✅ **Health check endpoint** for monitoring
- ✅ **Docker ready** with multi-stage build
- ✅ **Production-ready** with comprehensive tests
- ✅ **Zero external dependencies** (except Swagger)

## 🛠️ Technologies

**Backend:**
- Go 1.24
- Standard library (net/http)
- Swagger/OpenAPI (github.com/swaggo/http-swagger)
- Docker & Docker Compose

**Testing:**
- Go testing framework
- K6 for load testing
- golangci-lint for code quality

## 📋 Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (optional)
- K6 for performance testing (optional)
- golangci-lint for linting (optional)

### Tool Installation

#### Go
```bash
# MacOS
brew install go

# Ubuntu/Debian
sudo apt install golang-1.24
```

### 📥 IP2Location Dataset

The project includes a **sample CSV file** (`data/sample.csv`) with 20 IP ranges for testing.

**For full dataset with 2.9M+ records:**

1Extract and place `IP2LOCATION-LITE-DB11.CSV` in the `data/` directory
2Update `.env`:
```bash
   CSV_FILE_PATH=data/IP2LOCATION-LITE-DB11.CSV
```

**Note:** The full dataset (330MB) is not included in the repository. See `data/README.md` for details.

#### Docker
```bash
# MacOS
brew install --cask docker

# Ubuntu
sudo apt install docker.io docker-compose
```

#### K6 (Performance Testing)
```bash
# MacOS
brew install k6

# Ubuntu/Debian
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

#### golangci-lint (Code Linting)
```bash
# MacOS
brew install golangci-lint

# Ubuntu/Debian
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### Swag (Swagger Documentation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🚀 Quick Start

### 📦 Option 1: Run with Docker (Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/moura95/arena-backend-challenge
cd arena-backend-challenge

# 2. Build and run with docker-compose
make docker-run

# The API will be available at:
# - API: http://localhost:8080
# - Health: http://localhost:8080/health
# - Swagger: http://localhost:8080/swagger/index.html
```

### 🔧 Option 2: Run Locally

```bash
# 1. Clone the repository
git clone https://github.com/moura95/arena-backend-challenge
cd arena-backend-challenge

# 2. Download dependencies
go mod download

# 3. Run the server
make run

# Or build and run the binary
make build
./bin/server
```

The API will be available at `http://localhost:8080`

## 🧪 Testing

### Unit Tests
```bash
# Run all tests with coverage
make test

# Output:
# Running tests...
# go test -v -race -coverprofile=coverage.out ./...
# === RUN   TestIPToID
# --- PASS: TestIPToID (0.00s)
# === RUN   TestMemoryRepository_FindByIPID
# --- PASS: TestMemoryRepository_FindByIPID (0.01s)
# ...
# PASS
# coverage: 85.2% of statements
```

### Performance Tests (K6)
```bash
# Start the server first
make run

# In another terminal, run K6 tests
make test-performance
```

**K6 Results:**
```
✓ p(95)<100ms - 95% of requests under 100ms
  
Checks:
  ✓ status is correct........: 96.95% (220,201 passed)
  ✓ response time < 100ms....: 100%
  ✓ content-type is JSON.....: 100%

Performance:
  http_req_duration..........: avg=238µs  p(95)=370µs  max=55ms
  http_reqs..................: 257,092 (1,833 req/s)
  vus........................: 1,000 concurrent users
  iteration_duration.........: avg=200ms
  
Network:
  data_received..............: 48 MB (342 kB/s)
  data_sent..................: 24 MB (171 kB/s)
```

### Code Linting
```bash
# Run linter
make lint

# Run linter with auto-fix
make lint-fix
```

## 📚 API Endpoints

### 🌍 IP Location Lookup
```http
GET /ip/location?ip={ipv4_address}
```

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ip` | string | Yes | IPv4 address in dotted decimal notation (e.g., "8.8.8.8") |

**Success Response (200 OK):**
```json
{
  "country": "United States",
  "countryCode": "US",
  "city": "Mountain View"
}
```

**Error Responses:**

*400 Bad Request - Invalid IP format:*
```json
{
  "error": "invalid IP address: invalid octet 1: strconv.ParseUint: parsing \"256\": value out of range"
}
```

*404 Not Found - IP not in database:*
```json
{
  "error": "Location not found for the given IP"
}
```

### ❤️ Health Check
```http
GET /health
```

**Success Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-18T14:30:45.123456Z",
}
```

### 📖 Swagger Documentation
```http
GET /swagger/index.html
```

Interactive API documentation with request/response examples and the ability to test endpoints directly from the browser.

## 💡 Usage Examples

### cURL Examples

**Get location for Google DNS (8.8.8.8):**
```bash
curl "http://localhost:8080/ip/location?ip=8.8.8.8"
```

**Get location for Cloudflare DNS (1.1.1.1):**
```bash
curl "http://localhost:8080/ip/location?ip=1.1.1.1"
```

**Health check:**
```bash
curl "http://localhost:8080/health"
```

## 🏗️ Architecture

The project follows a **layered architecture** with clear separation of concerns:

```
arena-backend-challenge/
├── cmd/
│   └── main.go                 # Application entry point
│
├── internal/
│   ├── handler/               # HTTP handlers (presentation layer)
│   │   ├── location_handler.go
│   │   └── location_handler_test.go
│   │
│   ├── service/               # Business logic
│   │   ├── location_service.go
│   │   └── location_service_test.go
│   │
│   ├── repository/            # Data access layer
│   │   ├── memory.go          # In-memory repository with binary search
│   │   ├── memory_test.go
│   │   └── mock_repository.go # Mock for testing
│   │
│   ├── domain/                # Domain entities and interfaces
│   │   └── location.go
│   │
│   └── server.go              # HTTP server setup and routing
│
├── api/v1/                    # API contracts (DTOs)
│   ├── location.go            # Location response
│   ├── health.go              # Health response
│   └── v1.go                  # Error response
│
├── pkg/
│   └── iputil/                # Utility packages
│       ├── converter.go       # IP to numeric ID conversion
│       └── converter_test.go
│
├── config/                    # Configuration management
│   ├── config.go
│   └── dotenv.go              # .env file parser
│
├── data/
│   └── IP2LOCATION-LITE-DB11.CSV  # IP geolocation dataset (2.9M+ records)
│
├── docs/                      # Generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yml         # Docker Compose configuration
├── Makefile                   # Build automation
├── go.mod                     # Go dependencies
└── test-k6.js                # K6 performance tests
```

### Layer Responsibilities

**Handler Layer** (`internal/handler/`)
- HTTP request/response handling
- Input validation
- Error formatting
- No business logic

**Service Layer** (`internal/service/`)
- Business logic orchestration
- IP validation and conversion
- Coordinates repository calls

**Repository Layer** (`internal/repository/`)
- Data access abstraction
- Binary search implementation
- CSV loading and parsing

**Domain Layer** (`internal/domain/`)
- Core business entities
- Repository interfaces
- Domain errors

## 🧠 Technical Decisions

### 1. **In-Memory Storage with Binary Search**

**Why:**
- Dataset is static (2.9M records, 330MB CSV)
- Fits comfortably in RAM (~400MB loaded)
- Binary search provides O(log n) lookup time
- Eliminates database overhead and network latency

**Trade-offs:**
- ✅ **Pros:**
  - Ultra-fast lookups (< 1ms)
  - Simple deployment (no database needed)
  - High throughput (1,800+ req/s)
  - No database maintenance

- ❌ **Cons:**
  - Requires server restart to update data
  - Memory usage: ~400MB per instance
  - Not suitable if dataset updates frequently

**Alternative considered:** PostgreSQL with indexed ranges
- Would add ~10-50ms latency per query
- Requires database management
- Better for frequently updated data
- Chosen approach is optimal for this static dataset use case

### 2. **Standard Library HTTP Server**

**Why:**
- Challenge requirement: minimize third-party dependencies
- Standard library is production-ready and well-tested
- Sufficient for simple REST API

**Trade-offs:**
- ✅ **Pros:**
  - Zero dependencies (except Swagger)
  - Lightweight and fast
  - Well-documented

- ❌ **Cons:**
  - More verbose routing
  - No built-in middleware chain
  - Manual request parsing

**Alternative considered:** Gin/Echo frameworks
- Would provide cleaner routing syntax
- Built-in middleware and validation
- Not used to respect challenge constraints

### 3. **Layered Architecture (Not Full Clean Architecture/DDD)**

**Why:**
- Project is simple and doesn't require complex domain logic
- Pragmatic approach: YAGNI (You Aren't Gonna Need It)
- Clear separation of concerns without over-engineering

**What I used:**
- **Handler → Service → Repository** pattern
- Domain entities without complex business rules
- Dependency injection for testability
- Interface-based repository for mocking

**What I didn't use (and why):**
- **Use Cases layer:** Business logic is simple enough for service layer
- **Domain events:** No complex state transitions
- **Value objects:** IP and Location are simple structs
- **Aggregates:** No complex entity relationships

**Trade-offs:**
- ✅ **Pros:**
  - Simpler codebase, easier to understand
  - Faster development for this scope
  - Still testable and maintainable
  - Easy to refactor to Clean Arch if needed

- ❌ **Cons:**
  - Less guidance for complex business rules
  - Domain logic could leak into service layer
  - Not ideal for large, complex domains

**When to use Clean Architecture:**
- Complex business rules and workflows
- Multiple bounded contexts
- Large team with many developers
- Long-term project evolution expected

**Why layered architecture is enough here:**
- Single, simple use case (IP lookup)
- No complex domain logic
- Small codebase (~1000 LOC)
- Performance is more important than domain modeling

### 4. **CSV Loading on Startup**

**Why:**
- Static dataset doesn't change at runtime
- One-time cost on server start (~2-3 seconds)
- All data available in memory for fast lookups

**Alternative considered:** Lazy loading
- Would delay first requests
- No benefit for static data
- Current approach preferred

### 5. **Docker Multi-Stage Build**

**Why:**
- Production image is only ~25MB
- Build tools not included in runtime
- Faster deployments and lower memory usage

**Build stages:**
1. **Builder stage:** Go 1.24 + git + swag CLI (~400MB)
2. **Runtime stage:** Alpine + binary + CSV (~25MB)

## 📊 Performance Characteristics

### Lookup Performance
- **Average latency:** 238µs (0.238ms)
- **P95 latency:** 370µs (0.370ms)
- **P99 latency:** < 1ms
- **Throughput:** 1,833 requests/second
- **Concurrency:** Tested up to 1,000 simultaneous users

### Memory Usage
- **CSV data loaded:** ~400MB
- **Total application:** ~450MB
- **Docker container:** ~25MB image + 450MB runtime

### Binary Search Complexity
- **Time complexity:** O(log n) where n = 2,979,950
- **Worst case lookups:** ~22 comparisons
- **Average lookups:** ~18 comparisons

### Startup Time
- **CSV loading:** 2-3 seconds
- **Server ready:** < 5 seconds total

## 🐳 Docker

### Building the Image
```bash
make docker-build
```

This creates a multi-stage Docker image:
- **Stage 1 (builder):** Compiles the Go binary with Swagger docs
- **Stage 2 (runtime):** Lightweight Alpine Linux with binary and CSV

**Final image size:** ~25MB

### Running with Docker Compose
```bash
# Start
make docker-run

# Stop
make docker-stop

# Clean up (remove images and containers)
make docker-clean
```

### Docker Compose Configuration
```yaml
version: '3.8'
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - HTTP_SERVER_ADDRESS=0.0.0.0:8080
      - CSV_FILE_PATH=data/IP2LOCATION-LITE-DB11.CSV
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

## 🔧 Makefile Commands

### Development
```bash
make build              # Build the binary
make run                # Run the application
make test               # Run all tests with coverage
make swagger            # Generate Swagger documentation
make lint               # Run golangci-lint
make lint-fix           # Run linter with auto-fix
make clean              # Clean build artifacts
```

### Testing
```bash
make test               # Run all tests
make test-performance   # Run K6 load tests
```

### Docker
```bash
make docker-build       # Build Docker image
make docker-run         # Start with docker-compose
make docker-stop        # Stop containers
make docker-clean       # Remove all Docker artifacts
```

## 🔐 Configuration

Configuration is managed through environment variables, loaded from `.env` file:

```bash
# .env
HTTP_SERVER_ADDRESS=0.0.0.0:8080
CSV_FILE_PATH=data/IP2LOCATION-LITE-DB11.CSV
```

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_SERVER_ADDRESS` | `0.0.0.0:8080` | Server bind address |
| `CSV_FILE_PATH` | `data/IP2LOCATION-LITE-DB11.CSV` | Path to IP location dataset |

## 🚦 CI/CD

The project includes GitHub Actions workflow (`.github/workflows/ci.yml`):

### Pipeline Stages
1. **Lint** - Code quality checks with golangci-lint
2. **Test** - Run all unit tests with race detection
3. **Build** - Compile the application

### Triggers
- Push to `main` or `master` branch
- Pull requests to `main` or `master` branch

## 🤝 Challenge Requirements

| Requirement | Status | Implementation |
|------------|--------|----------------|
| REST API with `/ip/location` endpoint | ✅ | Implemented with standard library |
| IP to ID conversion | ✅ | `pkg/iputil/converter.go` |
| Binary range search | ✅ | `internal/repository/memory.go` |
| Response time < 100ms | ✅ | **Average: 0.238ms (417x faster)** |
| Handle 100+ concurrent users | ✅ | **Tested: 1,000 concurrent users** |
| Unit tests | ✅ | 85%+ test coverage |
| No third-party libs (except REST) | ✅ | Only Swagger for docs |
| Layered architecture | ✅ | Handler → Service → Repository |
| K6 load testing | ✅ | Included in repository |

## 📈 Possible Improvements

1. **Caching Layer**
- Redis/Memcached for frequently queried IPs
- Would reduce repeated lookups
- Trade-off: Added complexity and infrastructure

2. **Rate Limiting**
- Prevent abuse
- Token bucket or sliding window
- Trade-off: Added middleware

3. **Metrics & Monitoring**
- Prometheus metrics endpoint
- Request latency histograms
- Error rate tracking
- Trade-off: Additional dependencies

4. **Database Option**
- PostgreSQL with GiST index on IP ranges
- Better for frequently updated datasets
- Trade-off: Higher latency, more infrastructure

## 👨‍💻 Author

**Guilherme Moura** - *Software Engineer*
- GitHub: [@moura95](https://github.com/moura95)
- LinkedIn: [Guilherme Moura](https://linkedin.com/in/guilherme-moura95)
- Email: junior.moura19@hotmail.com