# 📑 Project Index - Quick Navigation

Welcome to the Blog Community Platform project! This index will help you navigate through all the documentation and code.

## 📖 Documentation Files

### Getting Started (Read in this order)
1. **[README.md](README.md)** - Complete project documentation
   - Project overview
   - Features list
   - Installation instructions
   - API reference
   - Architecture details

2. **[QUICKSTART.md](QUICKSTART.md)** - Get running in 3 steps
   - Fastest way to run the project
   - Testing examples with cURL
   - Common troubleshooting
   - Tips and tricks

3. **[API.md](API.md)** - Complete API documentation
   - All endpoints documented
   - Request/response examples
   - Error codes and handling
   - Authentication details

4. **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Project statistics
   - Code metrics
   - Feature checklist
   - Learning outcomes
   - Future enhancements

## 🗂️ Code Structure

### Application Entry Point
- **[cmd/main.go](cmd/main.go)** - Main application file
  - Server configuration
  - Route setup
  - Middleware registration
  - Template functions

### Handlers (Request Processing)
- **[internal/handlers/auth.go](internal/handlers/auth.go)**
  - User registration
  - User login
  - Profile retrieval

- **[internal/handlers/posts.go](internal/handlers/posts.go)**
  - Get all posts (with pagination)
  - Get single post
  - Create post
  - Update post
  - Delete post

- **[internal/handlers/comments.go](internal/handlers/comments.go)**
  - Get comments for post
  - Create comment
  - Delete comment

- **[internal/handlers/upload.go](internal/handlers/upload.go)**
  - Image upload handling
  - File validation
  - Storage management

- **[internal/handlers/web.go](internal/handlers/web.go)**
  - Home page rendering
  - Posts listing page
  - Single post page
  - Login/Register pages
  - Admin dashboard

### Middleware (Request Processing)
- **[internal/middleware/auth.go](internal/middleware/auth.go)**
  - Token validation
  - User authentication
  - Role-based authorization

- **[internal/middleware/cors.go](internal/middleware/cors.go)**
  - Cross-origin handling
  - Preflight requests

- **[internal/middleware/logger.go](internal/middleware/logger.go)**
  - Request logging
  - Response time tracking

- **[internal/middleware/ratelimit.go](internal/middleware/ratelimit.go)**
  - Per-IP rate limiting
  - Configurable limits

- **[internal/middleware/requestid.go](internal/middleware/requestid.go)**
  - Unique request ID generation
  - Request tracking

### Data Models
- **[internal/models/user.go](internal/models/user.go)**
  - User struct
  - Login/Register requests

- **[internal/models/post.go](internal/models/post.go)**
  - Post struct
  - Create/Update requests

- **[internal/models/comment.go](internal/models/comment.go)**
  - Comment struct
  - Comment requests

## 🎨 Frontend Files

### HTML Templates
- **[web/templates/index.html](web/templates/index.html)** - Home page
- **[web/templates/posts.html](web/templates/posts.html)** - Posts listing
- **[web/templates/post.html](web/templates/post.html)** - Single post view
- **[web/templates/login.html](web/templates/login.html)** - Login page
- **[web/templates/register.html](web/templates/register.html)** - Registration
- **[web/templates/admin-dashboard.html](web/templates/admin-dashboard.html)** - Admin panel
- **[web/templates/404.html](web/templates/404.html)** - Error page

### CSS Styles
- **[web/static/css/style.css](web/static/css/style.css)** - Complete styling
  - Responsive design
  - Component styles
  - Animations
  - Theme colors

### JavaScript
- **[web/static/js/main.js](web/static/js/main.js)** - Core functionality
  - Notification system
  - API helpers
  - Authentication utilities
  - DOM manipulation

- **[web/static/js/auth.js](web/static/js/auth.js)** - Authentication
  - Login form handling
  - Registration form handling
  - Token management

## ⚙️ Configuration

- **[go.mod](go.mod)** - Go module definition
- **[Makefile](Makefile)** - Build automation
- **[.gitignore](.gitignore)** - Git ignore rules

## 📁 Directories

- **cmd/** - Application entry points
- **internal/** - Internal application code
  - **handlers/** - HTTP request handlers
  - **middleware/** - Custom middleware
  - **models/** - Data models
- **pkg/** - Public packages (utilities)
- **web/** - Web assets
  - **templates/** - HTML templates
  - **static/** - Static files (CSS, JS, images)
- **uploads/** - User uploaded files
- **config/** - Configuration files

## 🚀 Quick Commands

```bash
# Install dependencies
make install

# Run the application
make run

# Build binary
make build

# Run tests
make test

# Clean artifacts
make clean
```

## 🎯 Key Endpoints

### Web Pages
- `http://localhost:8080/` - Home
- `http://localhost:8080/posts` - Posts listing
- `http://localhost:8080/posts/:id` - Single post
- `http://localhost:8080/login` - Login
- `http://localhost:8080/register` - Register
- `http://localhost:8080/admin/dashboard` - Admin (requires auth)

### API Endpoints
- `GET /api/v1/health` - Health check
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/posts` - Get posts
- `POST /api/v1/posts` - Create post (auth required)
- `GET /api/v1/posts/:id/comments` - Get comments
- `POST /api/v1/posts/comments` - Add comment (auth required)
- `POST /api/v1/upload/image` - Upload image (auth required)

## 📚 Learning Path

### Beginner
1. Read [QUICKSTART.md](QUICKSTART.md)
2. Run the application
3. Explore the web interface
4. Try the API endpoints

### Intermediate
1. Read [README.md](README.md)
2. Study [cmd/main.go](cmd/main.go)
3. Explore handlers in [internal/handlers/](internal/handlers/)
4. Understand middleware in [internal/middleware/](internal/middleware/)

### Advanced
1. Read [API.md](API.md)
2. Study all code files
3. Modify and extend features
4. Implement suggested enhancements from [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

## 🎓 Lesson Mapping

### Lesson 01: Basic Server
- See: [cmd/main.go](cmd/main.go) - Server setup
- Concepts: Gin initialization, basic routing

### Lesson 02: CRUD
- See: [internal/handlers/posts.go](internal/handlers/posts.go)
- Concepts: RESTful operations, HTTP methods

### Lesson 03: Parameter Binding
- See: All handler files
- Concepts: JSON binding, validation, query params

### Lesson 04: Context
- See: [internal/handlers/](internal/handlers/), [internal/middleware/](internal/middleware/)
- Concepts: Context storage, retrieval, passing

### Lesson 05: Middleware
- See: [internal/middleware/](internal/middleware/)
- Concepts: Request/response processing, chain

### Lesson 06: Versioning & Groups
- See: [cmd/main.go](cmd/main.go) - Route groups
- Concepts: API versioning, route organization

### Lesson 07: Static Files
- See: [cmd/main.go](cmd/main.go) - Static routes
- See: [web/static/](web/static/)
- Concepts: File serving, uploads

### Lesson 08: Templates
- See: [internal/handlers/web.go](internal/handlers/web.go)
- See: [web/templates/](web/templates/)
- Concepts: HTML rendering, template functions

## 🔍 Search Tips

### Find Specific Features
- **Authentication**: Search for "auth" in handlers and middleware
- **CRUD Operations**: Look in [internal/handlers/posts.go](internal/handlers/posts.go)
- **Validation**: Search for "binding" tags in models
- **Error Handling**: Search for "JSON" responses in handlers
- **Middleware**: Check [internal/middleware/](internal/middleware/)
- **Templates**: Browse [web/templates/](web/templates/)
- **Styles**: See [web/static/css/style.css](web/static/css/style.css)

## 📞 Help & Support

If you're stuck:
1. Check [QUICKSTART.md](QUICKSTART.md) for common issues
2. Read [README.md](README.md) for detailed explanations
3. Review [API.md](API.md) for endpoint details
4. Check code comments in source files
5. Refer to Gin documentation: https://gin-gonic.com/

## ✅ Checklist for Understanding

- [ ] Ran the application successfully
- [ ] Explored the web interface
- [ ] Tested API endpoints with cURL
- [ ] Read through main.go
- [ ] Understood the middleware flow
- [ ] Explored handler implementations
- [ ] Reviewed data models
- [ ] Examined templates
- [ ] Checked CSS styling
- [ ] Understood JavaScript functionality

## 🎉 Next Steps

1. ✅ Complete the checklist above
2. 📝 Modify existing features
3. ➕ Add new features
4. 🧪 Write tests
5. 🚀 Deploy to production (with proper database)

---

**Happy Learning! 📚**

This project contains everything you need to understand and build production-ready Go web applications with Gin framework.

---

Last Updated: October 3, 2025
