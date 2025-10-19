# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2025-10-20

### Added
- Initial release
- REST API endpoint `/ip/location` for IP geolocation
- Binary search implementation for 2.9M+ IP ranges
- In-memory repository with O(log n) lookups
- Comprehensive unit tests (85%+ coverage)
- K6 performance tests
- Docker support with multi-stage build
- Swagger/OpenAPI documentation
- Health check endpoint
- CI/CD with GitHub Actions
- Detailed README with architecture decisions

### Performance
- Average response time: 0.238ms
- P95 latency: 0.370ms
- Throughput: 1,833 requests/second
- Concurrent users tested: 1,000

### Technical Details
- Go 1.24
- Standard library HTTP server
- Zero external dependencies (except Swagger)
- Memory usage: ~450MB with full dataset
- Docker image size: ~25MB