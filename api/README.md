# API Definitions

This directory contains the Interface Definition Language (IDL) files for the Comical Tool Short URL service.

## Files

### 1. `short_url.thrift`
Apache Thrift definition file containing:
- Data structures for all API requests and responses
- Service interface definition
- Type definitions for analytics and pagination

### 2. `short_url.proto`
Protocol Buffers definition file containing:
- gRPC service definitions
- Message types for all API operations
- Field definitions with proper types and validation

### 3. `short_url.yaml`
OpenAPI 3.0 specification containing:
- Complete REST API documentation
- Request/response schemas
- Parameter definitions
- Example requests and responses
- Error handling documentation

## Code Generation

To generate Go code from these IDL files, run:

```bash
# Make the script executable
chmod +x scripts/generate_idl.sh

# Run the generation script
./scripts/generate_idl.sh
```

This will generate:
- Thrift-generated Go code in `api/generated/thrift/`
- Protocol Buffer-generated Go code in `api/generated/proto/`

## Prerequisites

Before running the generation script, ensure you have:

1. **Apache Thrift**: Install from [thrift.apache.org](https://thrift.apache.org/)
2. **Protocol Buffers**: Install from [protobuf.dev](https://protobuf.dev/)
3. **Go plugins**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Usage in Hertz

The generated code can be used in your Hertz application for:

1. **Type Safety**: Use generated structs for request/response handling
2. **Validation**: Leverage built-in validation from IDL definitions
3. **Documentation**: Auto-generate API documentation from OpenAPI spec
4. **Client Generation**: Generate client SDKs for other languages

## API Endpoints

The service provides the following endpoints:

- `POST /api/v1/urls` - Create short URL
- `GET /api/v1/urls/{code}` - Get short URL details
- `PUT /api/v1/urls/{code}` - Update short URL
- `DELETE /api/v1/urls/{code}` - Delete short URL
- `GET /api/v1/urls/{code}/analytics` - Get analytics
- `GET /api/v1/urls/{code}/clicks` - Get click history
- `GET /{code}` - Redirect to original URL

## Configuration

The API supports various configuration options through environment variables:

- `SHORT_URL_DOMAIN` - Domain for short URLs
- `SHORT_URL_CODE_LENGTH` - Length of generated codes
- `SHORT_URL_ALLOWED_CHARS` - Characters allowed in codes
- `SHORT_URL_DEFAULT_EXPIRY_HOURS` - Default expiration time
- `SHORT_URL_ANALYTICS_RETENTION_DAYS` - Analytics data retention period
