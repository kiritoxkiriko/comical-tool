# comical-tool

Comical Tool is a powerful tool for providing common web tools. This project contains multiple services including short URL, file sharing, and image hosting.

## Features

### Short URL Service ✅
- ✅ Custom alias or auto-generated codes
- ✅ Expiration time (time-based and click-based)
- ✅ Analytics (click count, referrer, location, user agent)
- ✅ Configuration options (URL length, allowed characters)
- ✅ Redis caching for high performance
- ✅ RESTful API with complete CRUD operations

### Planned Features
- [ ] File Sharing
- [ ] Image Hosting

## Tech Stack
* **Backend:** Golang
* **Web Framework:** [Hertz](https://github.com/cloudwego/hertz)
* **Database:** MySQL
* **Cache:** Redis
* **Configuration:** Viper
* **CLI:** Cobra
* **Containerization:** Docker

## Quick Start

### Prerequisites
- Go 1.25+
- MySQL 8.0+
- Redis 6.0+

### Installation

```bash
# Clone the repository
git clone https://github.com/kiritoxkiriko/comical-tool.git
cd comical-tool

# Install dependencies
go mod tidy

# Build the application
make build
```

### Configuration

The application uses Viper for configuration management. You can configure it via:

1. **Config file** (recommended): `config.yaml`
2. **Environment variables**: `COMICAL_*`
3. **Command line flags**

```bash
# Initialize default config file
./bin/comical-tool config init

# Show current configuration
./bin/comical-tool config show
```

### Running the Server

```bash
# Start server with default configuration
./bin/comical-tool server

# Start server with custom host and port
./bin/comical-tool server --host 0.0.0.0 --port 8080

# Start server with custom config file
./bin/comical-tool server --config /path/to/config.yaml
```

### Using Docker

```bash
# Start all services with Docker Compose
make docker-run

# Or manually
docker-compose up -d
```

## API Endpoints

### Short URL Management
- `POST /api/v1/urls` - Create short URL
- `GET /api/v1/urls/{code}` - Get short URL details
- `PUT /api/v1/urls/{code}` - Update short URL
- `DELETE /api/v1/urls/{code}` - Delete short URL

### Analytics
- `GET /api/v1/urls/{code}/analytics` - Get analytics data
- `GET /api/v1/urls/{code}/clicks` - Get click history

### Redirection
- `GET /{code}` - Redirect to original URL

## CLI Commands

```bash
# Show help
./bin/comical-tool --help

# Server commands
./bin/comical-tool server --help

# Configuration management
./bin/comical-tool config show
./bin/comical-tool config init

# Version information
./bin/comical-tool version
```

## Development

```bash
# Run tests
make test

# Build application
make build

# Run server in development mode
make run-server

# Clean build artifacts
make clean
```

