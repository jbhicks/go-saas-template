# Go SaaS Template

A ready-to-use template for building SaaS applications with Go.

## Project Overview
This is a general purpose Go SaaS template that provides a solid foundation for building web applications. It includes authentication, database integration, and a simple web interface to get you started quickly.

## Tech Stack
- **Backend**: Go with [Gorilla Mux](https://github.com/gorilla/mux) router
- **Frontend**: HTMX/CSS/JavaScript for interactive UI components
- **Database**: [PocketBase](https://pocketbase.io/) for data storage and authentication

## Core Features
- [x] PocketBase backend integration
  - [x] Built-in authentication
  - [ ] Define your custom collections
  - [ ] Set up proper permissions for your data model
- [x] User authentication flows
  - [x] Login/Register
  - [x] Password reset
- [ ] Ready-to-customize templates
- [ ] API endpoints for your business logic
- [ ] Dockerized development environment

## Getting Started

```bash
# Clone the repository
git clone https://github.com/yourusername/go-saas-template.git

# Navigate to the project directory
cd go-saas-template

# Install dependencies
go mod download

# Run the application
go run cmd/server/main.go

# Access the application
# Open http://localhost:8080 in your browser
```

## Troubleshooting
### Complete Reset for PocketBase Docker

If experiencing caching issues or the volume might not be fully removed, try a more thorough reset:

#### 1. Force Stop and Remove All Containers Using the Volume

```bash
# Find and force remove all containers using the volume
docker ps -a -q --filter volume=pocketbase-data | xargs docker rm -f
```

#### 2. Force Remove the Volume

```bash
docker volume rm -f pocketbase-data
```

#### 3. Create a New Container with a Different Volume Name

```bash
docker run -p 8090:8090 -v pb_data_new:/pb_data ghcr.io/muchobien/pocketbase:latest
```

#### 4. Clear Browser Cache

After running the above commands:
1. Open a private/incognito window in your browser
2. Visit http://localhost:8090/_/
3. You should see the initial setup screen

#### 5. If Still Not Working

Try running without a named volume to verify the issue:

```bash
# Run with an anonymous volume
docker run -p 8090:8090 ghcr.io/muchobien/pocketbase:latest
```

This should definitely give you a fresh instance without any existing data.

#### Browser Cache Issue?

If you still see a login screen, it's likely your browser is caching the page. Try:
- Using an incognito/private window
- Different browser
- Adding a cache-busting parameter: http://localhost:8090/_/?nocache=1