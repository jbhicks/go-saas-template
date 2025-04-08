# city-tiers

Helping keep track of all your favorite things.

## Project Overview
City Tiers is a web application that helps users categorize and rank their favorite places, activities, and things within a city. Users can create lists organized by category (restaurants, parks, attractions) and assign tier rankings (S, A, B, C, etc.).

## Tech Stack
- **Backend**: Go with [Gorilla Mux](https://github.com/gorilla/mux) router
- **Frontend**: HTMX/CSS/JavaScript when no other choice is available
- **Database**: Pocketbase

## Core Features
- [x] Pocketbase backend setup
  - [ ] Define your collections in Pocketbase for Users, Categories, and Places
  - [ ] Use Pocketbase's built-in auth for user management
  - [ ] Set up proper permissions on your collections
- [ ] Create your HTMX frontend to interact with Pocketbase API
- [ ] User authentication
- [ ] Create/edit/delete categories 
- [ ] Add places to categories with tier rankings
- [ ] View all places by category or by tier
- [ ] Map integration to visualize places
- [ ] Share lists with others

## Data Models

```go
// User represents a registered user
type User struct {
    ID       string
    Username string
    Email    string
    Password string // Stored securely
}

// Category represents a grouping of places (Restaurants, Parks, etc.)
type Category struct {
    ID     string
    Name   string
    UserID string
}

// Place represents a location within a city
type Place struct {
    ID          string
    Name        string
    Description string
    Address     string
    CategoryID  string
    Tier        string // S, A, B, C, etc.
    UserID      string
    Latitude    float64
    Longitude   float64
}
```

## API Endpoints

- `GET /api/categories` - List all categories
- `POST /api/categories` - Create a new category
- `GET /api/categories/{id}` - Get a specific category
- `PUT /api/categories/{id}` - Update a category
- `DELETE /api/categories/{id}` - Delete a category
- `GET /api/categories/{id}/places` - Get places in a category
- `POST /api/places` - Create a new place
- `GET /api/places/{id}` - Get a specific place
- `PUT /api/places/{id}` - Update a place
- `DELETE /api/places/{id}` - Delete a place

## Implementation Roadmap

1. **Phase 1**: Basic CRUD operations
   - Set up database
   - Implement API endpoints for categories and places
   - Create basic UI for managing categories and places

2. **Phase 2**: User Authentication
   - Implement user registration and login
   - Add user-specific data isolation

3. **Phase 3**: Enhanced Features
   - Add map visualization
   - Implement sharing functionality
   - Add search capabilities

## Development Setup

```bash
# Clone the repository
git clone https://github.com/jbhicks/city-tiers.git

# Navigate to the project directory
cd city-tiers

# Install dependencies
go mod download

# Run the application
go run main.go

# Access the application
# Open http://localhost:8080 in your browser
```

## Troubleshooting
### Complete Reset for Pocketbase Docker

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