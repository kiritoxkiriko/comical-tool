# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Comical Tool is a backend service for providing common web tools, built with Go and the Kitex RPC framework. The project is designed to support features like short URL generation, file sharing, and image hosting.

## Tech Stack

- **Language**: Go
- **RPC Framework**: Kitex (CloudWeGo)
- **Database**: MySQL
- **Cache**: Redis
- **Containerization**: Docker

## Project Structure

This is a Go backend service project. As the codebase develops, the typical structure should follow:
- Service definitions using Kitex IDL (Thrift/Protobuf)
- Handler implementations for RPC services
- Database models and repositories
- Redis caching layer
- Docker configuration for deployment

## Development Commands

Since this is a Go project using Kitex, common commands will include:

```bash
# Initialize Go module (if not already done)
go mod init github.com/[username]/comical-tool

# Install Kitex (if not already installed)
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# Generate code from IDL files (once IDL files are created)
kitex -module github.com/[username]/comical-tool [idl_file]

# Run tests
go test ./...

# Build the project
go build -o comical-tool

# Run the service
./comical-tool
```

## Architecture Notes

The project will implement three main features:
1. **Short URL Service**: URL shortening and redirection
2. **File Sharing Service**: Temporary file storage and sharing
3. **Image Hosting Service**: Image upload and serving

Each service should be implemented as a separate Kitex service handler with appropriate database models and Redis caching strategies.