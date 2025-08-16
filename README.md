# Clans - Reddit Clone Microservices Backend

A Reddit-style community platform built with Go microservices architecture, featuring clans (communities), threaded comments, voting systems, and user management.

## Quick Deploy

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/clans-reddit-clone)

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │────│  API Gateway    │
│   (React)       │    │   (Port 8000)   │
└─────────────────┘    └─────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼────┐  ┌───────▼────┐  ┌──────▼─────┐
        │ User Service│  │Post Service│  │Comment Svc │
        │ (Port 8080) │  │(Port 8081) │  │(Port 8082) │
        └─────────────┘  └─────────────┘  └────────────┘
                │               │               │
                └───────────────┼───────────────┘
                                │
                        ┌───────▼────┐  ┌──────────────┐
                        │Clan Service│  │  PostgreSQL  │
                        │(Port 8083) │  │  (Port 5432) │
                        └────────────┘  └──────────────┘
```

## Services

### 🔒 User Service (Port 8080)
- User registration and authentication
- JWT token management
- Password hashing with bcrypt
- User profile management

### 📝 Post Service (Port 8081)
- CRUD operations for posts
- Post voting system (upvote/downvote)
- Clan-based post organization
- Vote count aggregation

### 💬 Comment Service (Port 8082)
- Threaded comment system with unlimited nesting
- Comment voting (upvote/downvote)
- Reply depth tracking
- Nested replies in API responses

### 🏰 Clan Service (Port 8083)
- Clan (community) creation and management
- Membership system with roles (owner, moderator, member)
- Public/private clan visibility
- Clan statistics (member count, post count)

### 🌐 API Gateway (Port 8000)
- Centralized request routing
- JWT authentication middleware
- CORS handling
- Request logging
- Service health monitoring

## Database Schema

### Core Tables
- **users** - User accounts and authentication
- **clans** - Communities/subreddits
- **clan_memberships** - User-clan relationships with roles
- **posts** - Content posts linked to clans
- **comments** - Threaded comments with parent-child relationships
- **post_votes** - Post voting records
- **comment_votes** - Comment voting records

## API Endpoints

### Authentication
```http
POST /api/auth/signup     # Register new user
POST /api/auth/login      # User login
```

### Clans
```http
GET    /api/clans                    # List public clans
POST   /api/clans                    # Create clan (auth required)
GET    /api/clans/{id}               # Get clan by ID
GET    /api/clans/name/{name}        # Get clan by name
PUT    /api/clans/{id}               # Update clan (auth required)
DELETE /api/clans/{id}               # Delete clan (auth required)
POST   /api/clans/{id}/join          # Join clan (auth required)
POST   /api/clans/{id}/leave         # Leave clan (auth required)
GET    /api/clans/{id}/members       # Get clan members
GET    /api/users/clans              # Get user's clans (auth required)
```

### Posts
```http
GET    /api/posts           # List all posts
POST   /api/posts           # Create post (auth required)
GET    /api/posts/{id}      # Get specific post
PUT    /api/posts/{id}      # Update post (auth required)
DELETE /api/posts/{id}      # Delete post (auth required)
POST   /api/posts/{id}/vote # Vote on post (auth required)
```

### Comments
```http
GET    /api/comments/post/{postId}  # Get threaded comments for post
POST   /api/comments               # Create comment/reply (auth required)
GET    /api/comments/{id}          # Get specific comment
PUT    /api/comments/{id}          # Update comment (auth required)
DELETE /api/comments/{id}          # Delete comment (auth required)
POST   /api/comments/{id}/vote     # Vote on comment (auth required)
```

## Key Features

### 🧵 Threaded Comments
- Unlimited reply nesting with depth tracking
- Nested JSON responses for easy frontend rendering
- Reply count tracking per comment

### 🗳️ Voting System
- Upvote/downvote for posts and comments
- Vote count aggregation
- One vote per user per item

### 🏘️ Clan System
- Public clans (discoverable) vs private clans (invite-only)
- Role-based permissions (owner, moderator, member)
- Clan statistics and member management

### 🔐 Authentication & Security
- JWT-based authentication
- Centralized auth at API Gateway
- User context forwarded to services
- Public endpoints for reading, auth required for writing

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for development)

### Running the System
```bash
# Clone repository
git clone <repository-url>
cd clans

# Start all services
docker-compose up --build

# Access the API
curl http://localhost:8000/api/clans
```

### Database Setup
The database schema is automatically created through migration files in each service's `migrations/` directory.

## Development

### Adding New Services
1. Create service directory structure following existing pattern
2. Add to `docker-compose.yml`
3. Update API Gateway routing in `internal/proxy/proxy.go`
4. Update auth middleware for public endpoints

### Database Migrations
Each service contains SQL migration files in `migrations/init.sql`. Run them manually if needed:
```bash
docker exec -i postgres-clans psql -U admin -d clans < clans/service-name/migrations/init.sql
```

### Service Communication
- Services communicate via HTTP through the API Gateway
- No direct service-to-service communication
- User context passed via `X-User-ID` header

## Frontend Integration

### Authentication Flow
1. User logs in via `/api/auth/login`
2. Store JWT token from response
3. Include token in `Authorization: Bearer <token>` header
4. Use `/api/users/clans` for user's personal clans

### Creating Content
- **Posts**: Include `clan_id` in request body
- **Comments**: Include `post_id` and optional `parent_id` for replies
- **Clans**: Include `is_public` boolean (defaults to true)

### Threaded Comments
Comments are returned with nested `replies` arrays. Render recursively:
```javascript
const CommentThread = ({ comments }) => (
  comments.map(comment => (
    <div key={comment.id}>
      <CommentItem comment={comment} />
      {comment.replies && <CommentThread comments={comment.replies} />}
    </div>
  ))
);
```

## Configuration

### Environment Variables
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `PORT` - Service port (defaults: gateway 8000, services 8080-8083)
- `JWT_SECRET` - JWT signing secret
- Service URLs for API Gateway routing

### Docker Compose
All services are orchestrated via `docker-compose.yml` with health checks and dependency management.

## Troubleshooting

### Common Issues
1. **Service Unavailable**: Check if `go.sum` exists, run `go mod tidy`
2. **Auth Errors**: Verify JWT token and API Gateway routing
3. **CORS Issues**: Ensure only API Gateway sets CORS headers
4. **Database Errors**: Check foreign key constraints and migrations

### Logs
```bash
# View specific service logs
docker-compose logs service-name

# View all logs
docker-compose logs
```

## Contributing

1. Follow existing service patterns
2. Add proper error handling and validation
3. Update API documentation
4. Test with Docker Compose
5. Ensure database migrations are included