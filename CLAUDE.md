# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go HTTP API server (`github.com/zapi-sh/api`) that provides a simple REST API. The codebase uses Go 1.25.3 and the standard library's `net/http` package.

## Development Philosophy

**Prefer standard library**: Always use Go's standard library over third-party dependencies unless explicitly stated otherwise. This project prioritizes simplicity and minimizes external dependencies.

## Development Commands

### Running the server
```bash
go run main.go
```

The server starts on port 8080 with these configured timeouts:
- ReadTimeout: 10s
- WriteTimeout: 10s
- MaxHeaderBytes: 1MB

### Building
```bash
go build -o api main.go
```

### Testing
```bash
go test ./...
```

## Architecture

**Single-file architecture**: The entire API is currently in `main.go` with:
- HTTP server setup using `http.Server` with explicit timeouts
- Route-based handlers using `http.ServeMux` with method prefixes (e.g., `GET /status`)
- Standard response structs: `ResponseSuccess` and `ResponseError` for consistent JSON responses

**Response structure**: All JSON responses follow a structured format with `status` field and either `data` (success) or `error` (failure) fields.

## API Endpoints

- `GET /status` - Health check endpoint returning service status
