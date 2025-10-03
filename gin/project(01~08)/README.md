# 📝 Blog Community Platform - Comprehensive Gin Framework Project

A full-featured blog and community platform built with the Gin web framework, demonstrating all essential features learned from lessons 01-08. This project serves as a production-ready example showcasing best practices for building web applications with Go and Gin.

## 🎯 Project Overview

This comprehensive project combines all the concepts covered in the Gin framework tutorial series:

- **Lesson 01**: Basic server setup and routing
- **Lesson 02**: Full CRUD operations
- **Lesson 03**: Parameter binding (JSON, form, query)
- **Lesson 04**: Context management and usage
- **Lesson 05**: Middleware (authentication, logging, CORS, rate limiting)
- **Lesson 06**: API versioning with route groups
- **Lesson 07**: Static file serving
- **Lesson 08**: HTML template rendering

## ✨ Features

### Core Features

- 🔐 **User Authentication**: Register, login, and session management
- 📚 **Blog Posts**: Full CRUD operations for blog posts
- 💬 **Comments System**: Nested comments with author information
- 📁 **File Upload**: Image upload functionality with validation
- 🔒 **Authorization**: Role-based access control (Admin/User)
- 🎨 **Responsive UI**: Modern, mobile-friendly interface
- 📡 **RESTful API**: Versioned API endpoints (v1 & v2)
- 🛡️ **Middleware Stack**: Comprehensive middleware implementation

### Technical Features

- **Parameter Binding**: JSON, form, and query parameter validation
- **Context Management**: Proper context usage throughout handlers
- **Rate Limiting**: Protect APIs from abuse
- **CORS Support**: Cross-origin resource sharing
- **Request Logging**: Detailed request/response logging
- **Error Handling**: Consistent error responses
- **Template Rendering**: Dynamic HTML pages
- **Static File Serving**: CSS, JavaScript, and uploaded files

## 📁 Project Structure

```
project/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/              # HTTP request handlers
│   │   ├── auth.go           # Authentication handlers
│   │   ├── posts.go          # Blog post handlers
│   │   ├── comments.go       # Comment handlers
│   │   ├── upload.go         # File upload handler
│   │   └── web.go            # Web page handlers
│   ├── middleware/            # Custom middleware
│   │   ├── auth.go           # Authentication middleware
│   │   ├── cors.go           # CORS middleware
│   │   ├── logger.go         # Logging middleware
│   │   ├── ratelimit.go      # Rate limiting middleware
│   │   └── requestid.go      # Request ID middleware
│   └── models/                # Data models
│       ├── user.go           # User model
│       ├── post.go           # Post model
│       └── comment.go        # Comment model
├── pkg/
│   └── utils/                 # Utility functions
├── web/
│   ├── templates/             # HTML templates
│   │   ├── index.html
│   │   ├── posts.html
│   │   ├── post.html
│   │   ├── login.html
│   │   ├── register.html
│   │   ├── admin-dashboard.html
│   │   └── 404.html
│   └── static/                # Static assets
│       ├── css/
│       │   └── style.css
│       └── js/
│           ├── main.js
│           └── auth.js
├── uploads/                   # Uploaded files directory
├── config/                    # Configuration files
├── go.mod                     # Go module file
├── go.sum                     # Go dependencies checksum
├── Makefile                   # Build automation
├── .gitignore                 # Git ignore file
└── README.md                  # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone or navigate to the project**:
   ```bash
   cd /Users/ihyeongjin/dev/golang/gin/project
   ```

2. **Install dependencies**:
   ```bash
   make install
   # or
   go mod download
   ```

3. **Run the application**:
   ```bash
   make run
   # or
   go run cmd/main.go
   ```

4. **Open your browser**:
   ```
   http://localhost:8080
   ```

## 🔧 Available Commands

```bash
make install    # Install dependencies
make run        # Run the application
make build      # Build the binary
make test       # Run tests
make clean      # Clean build artifacts
make dev        # Run in development mode (with hot reload)
```

## 🌐 API Endpoints

### Public Endpoints

#### Health Check
```bash
GET /api/v1/health
```

Response:
```json
{
  "status": "healthy",
  "version": "1.0",
  "time": "2025-10-03T20:00:00Z"
}
```

#### Register User
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

Response:
```json
{
  "message": "Login successful",
  "token": "admin-token-456",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

#### Get Posts
```bash
GET /api/v1/posts?page=1&limit=10&category=Tutorial
```

Response:
```json
{
  "posts": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

#### Get Single Post
```bash
GET /api/v1/posts/:id
```

### Protected Endpoints (Require Authentication)

Add header: `Authorization: Bearer <token>`

#### Create Post
```bash
POST /api/v1/posts
Authorization: Bearer admin-token-456
Content-Type: application/json

{
  "title": "My New Post",
  "content": "This is the content of my post...",
  "category": "Tutorial",
  "tags": ["golang", "gin"]
}
```

#### Update Post
```bash
PUT /api/v1/posts/:id
Authorization: Bearer admin-token-456
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

#### Delete Post
```bash
DELETE /api/v1/posts/:id
Authorization: Bearer admin-token-456
```

#### Get Comments
```bash
GET /api/v1/posts/:id/comments
```

#### Create Comment
```bash
POST /api/v1/posts/comments
Authorization: Bearer valid-token-123
Content-Type: application/json

{
  "post_id": 1,
  "content": "Great post!",
  "parent_id": null
}
```

#### Upload Image
```bash
POST /api/v1/upload/image
Authorization: Bearer valid-token-123
Content-Type: multipart/form-data

image: <file>
```

### API Version 2 Endpoints

```bash
GET /api/v2/health        # Enhanced health check with features list
GET /api/v2/posts         # Enhanced post listing
```

## 🔐 Test Credentials

### Admin Account
- **Username**: admin
- **Password**: password123
- **Token**: admin-token-456
- **Permissions**: Full access to all resources

### User Account
- **Username**: user
- **Password**: password123
- **Token**: valid-token-123
- **Permissions**: Can create, update, and delete own posts/comments

## 📖 Usage Examples

### Example 1: Register and Login

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

### Example 2: Create a Blog Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token-456" \
  -d '{
    "title": "Introduction to Gin Framework",
    "content": "Gin is a high-performance HTTP web framework...",
    "category": "Tutorial",
    "tags": ["golang", "gin", "web"]
  }'
```

### Example 3: Upload an Image

```bash
curl -X POST http://localhost:8080/api/v1/upload/image \
  -H "Authorization: Bearer admin-token-456" \
  -F "image=@/path/to/image.jpg"
```

### Example 4: Get Posts with Filtering

```bash
# Get all posts
curl http://localhost:8080/api/v1/posts

# Get posts by category
curl http://localhost:8080/api/v1/posts?category=Tutorial

# Get posts with pagination
curl http://localhost:8080/api/v1/posts?page=1&limit=5
```

## 🎨 Web Interface

The project includes a fully functional web interface accessible through:

- **Home Page**: `http://localhost:8080/`
- **Posts Listing**: `http://localhost:8080/posts`
- **Single Post**: `http://localhost:8080/posts/:id`
- **Login**: `http://localhost:8080/login`
- **Register**: `http://localhost:8080/register`
- **Admin Dashboard**: `http://localhost:8080/admin/dashboard` (requires admin role)

## 🛡️ Middleware Stack

### 1. Recovery Middleware
Recovers from panics and returns 500 error

### 2. Logger Middleware
Logs all requests with:
- Client IP
- HTTP method
- Request path
- Status code
- Request duration

### 3. CORS Middleware
Handles Cross-Origin Resource Sharing:
- Allows all origins (configurable)
- Supports credentials
- Handles preflight requests

### 4. Request ID Middleware
Adds unique ID to each request for tracking

### 5. Rate Limiting Middleware
Prevents API abuse:
- 100 requests/minute for API v1
- 150 requests/minute for API v2

### 6. Authentication Middleware
Validates JWT tokens and sets user context

### 7. Authorization Middleware
Checks user roles for protected resources

## 📚 Key Concepts Demonstrated

### 1. Basic Server Setup (Lesson 01)
```go
r := gin.New()
r.Use(gin.Recovery())
r.Run(":8080")
```

### 2. CRUD Operations (Lesson 02)
```go
posts.GET("", handlers.GetPosts)
posts.GET("/:id", handlers.GetPost)
posts.POST("", handlers.CreatePost)
posts.PUT("/:id", handlers.UpdatePost)
posts.DELETE("/:id", handlers.DeletePost)
```

### 3. Parameter Binding (Lesson 03)
```go
var req models.CreatePostRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 4. Context Usage (Lesson 04)
```go
userID, _ := c.Get("user_id")
c.Set("username", "admin")
c.JSON(200, gin.H{"data": result})
```

### 5. Middleware (Lesson 05)
```go
r.Use(middleware.LoggerMiddleware())
r.Use(middleware.AuthMiddleware())
r.Use(middleware.RateLimitMiddleware(100))
```

### 6. Route Groups & Versioning (Lesson 06)
```go
v1 := r.Group("/api/v1")
v2 := r.Group("/api/v2")
admin := r.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
```

### 7. Static Files (Lesson 07)
```go
r.Static("/static", "./web/static")
r.Static("/uploads", "./uploads")
r.StaticFile("/favicon.ico", "./web/static/favicon.ico")
```

### 8. Template Rendering (Lesson 08)
```go
r.LoadHTMLGlob("web/templates/*")
c.HTML(200, "index.html", gin.H{
    "title": "Home",
})
```

## 🔒 Security Features

1. **Authentication**: Token-based authentication
2. **Authorization**: Role-based access control
3. **Input Validation**: Request data validation
4. **Rate Limiting**: Prevent abuse
5. **CORS**: Controlled cross-origin access
6. **File Upload Validation**: Type and size checks
7. **SQL Injection Prevention**: Parameterized queries (when using DB)
8. **XSS Prevention**: Template auto-escaping

## 🚧 Future Enhancements

- [ ] Database integration (PostgreSQL/MySQL)
- [ ] Real JWT token implementation
- [ ] Redis for session storage
- [ ] WebSocket support for real-time features
- [ ] Email notifications
- [ ] Search functionality
- [ ] Pagination improvements
- [ ] File upload to cloud storage (S3)
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] API documentation with Swagger

## 📝 Development Notes

### In-Memory Storage
Currently, the application uses in-memory maps for data storage. This is for demonstration purposes only. In production:
- Use a proper database (PostgreSQL, MySQL, MongoDB)
- Implement data persistence
- Add transaction support
- Use proper session management

### Authentication
The current authentication is simplified. For production:
- Implement proper JWT tokens
- Add token refresh mechanism
- Use bcrypt for password hashing
- Implement session management
- Add OAuth2 support

### File Uploads
File uploads are stored locally. For production:
- Use cloud storage (AWS S3, Google Cloud Storage)
- Implement CDN for serving files
- Add image processing (resize, optimize)
- Implement virus scanning

## 🤝 Contributing

This is a learning project. Feel free to:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

This project is created for educational purposes as part of the Gin Framework tutorial series.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Go Programming Language](https://golang.org/)
- All contributors to the Gin framework

## 📞 Support

For questions or issues:
- Check the code comments
- Review the Gin documentation
- Refer to lessons 01-08 for specific concepts

---

**Built with ❤️ using Gin Framework**

This project demonstrates a production-ready structure while maintaining educational clarity. Each component is well-documented and follows Go best practices.
