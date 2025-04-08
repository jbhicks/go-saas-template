# PocketBase Documentation Reference Guide

## Overview
This project includes a local copy of the PocketBase documentation for reference. The documentation is stored in `docs/pocketbase-docs/` and covers both Go and JavaScript implementations.

## How to Use This Documentation
When working with PocketBase in this project:

1. **For Go Implementation**: Refer to the folders prefixed with `go-` in the docs directory
2. **For JavaScript/TypeScript Implementation**: Refer to the folders prefixed with `js-` in the docs directory
3. **For API Documentation**: Refer to the folders prefixed with `api-` in the docs directory
4. ***DO NOT!!!*** leave dumb ass obviously AI comments that explain the changes you are making. ***ONLY*** explain code!!!!    

## Key Documentation Sections

### Go Implementation
- `go-overview` - Basic introduction to PocketBase Go SDK
- `go-collections` - Working with collections in Go
- `go-records` - CRUD operations on records
- `go-database` - Database operations and queries
- `go-realtime` - Realtime subscriptions and events
- `go-event-hooks` - Event hooks and callbacks
- `go-authentication` - Authentication and users

### JavaScript Implementation
- `js-overview` - Basic introduction to PocketBase JS SDK
- `js-collections` - Working with collections in JavaScript
- `js-records` - CRUD operations on records
- `js-database` - Database operations and queries
- `js-realtime` - Realtime subscriptions and events
- `js-event-hooks` - Event hooks and callbacks

### API and Other Topics
- `authentication` - Authentication methods
- `collections` - Collection management
- `files-handling` - File uploads and storage
- `working-with-relations` - Managing relationships between collections

## Tips for Implementing Authentication
When implementing authentication with PocketBase, refer to:
- `authentication/` for general auth concepts
- `go-overview/` and other Go-related folders for Go implementation details
- The current implementation in this project uses PocketBase for authentication as seen in `internal/auth/middleware.go`

## Updating PocketBase
When updating PocketBase in this project, always check the updated documentation to ensure API compatibility.

## Versioning
This documentation was sourced from the official PocketBase site repository and represents the version as of April 2025. Always check for potential API changes if using newer versions of PocketBase.