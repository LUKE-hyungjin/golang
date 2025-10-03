# 📡 API Documentation

Complete API reference for the Blog Community Platform.

## Base URL

```
http://localhost:8080
```

## Authentication

Most protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <token>
```

## API Versioning

The API is versioned using URL path versioning:

- **v1**: `/api/v1/*` - Stable, production-ready
- **v2**: `/api/v2/*` - Enhanced features, improved responses

---

## 🔐 Authentication Endpoints

### Register User

Create a new user account.

**Endpoint**: `POST /api/v1/auth/register`

**Request Body**:
```json
{
  "username": "string (required, 3-20 chars)",
  "email": "string (required, valid email)",
  "password": "string (required, min 6 chars)"
}
```

**Response** (201 Created):
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 3,
    "username": "newuser",
    "email": "user@example.com",
    "role": "user"
  }
}
```

**Error Response** (409 Conflict):
```json
{
  "error": "Username already exists"
}
```

---

### Login

Authenticate and receive a token.

**Endpoint**: `POST /api/v1/auth/login`

**Request Body**:
```json
{
  "username": "string (required)",
  "password": "string (required, min 6 chars)"
}
```

**Response** (200 OK):
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

**Error Response** (401 Unauthorized):
```json
{
  "error": "Invalid username or password"
}
```

---

## 📝 Posts Endpoints

### Get All Posts

Retrieve a list of all published posts with pagination.

**Endpoint**: `GET /api/v1/posts`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `category` (optional): Filter by category

**Example Request**:
```bash
GET /api/v1/posts?page=1&limit=5&category=Tutorial
```

**Response** (200 OK):
```json
{
  "posts": [
    {
      "id": 1,
      "title": "Welcome to Our Blog Platform",
      "content": "This is our first blog post...",
      "author_id": 1,
      "author": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin"
      },
      "category": "Announcement",
      "tags": ["welcome", "first-post"],
      "image_url": "/static/images/blog1.jpg",
      "view_count": 42,
      "is_published": true,
      "created_at": "2025-09-03T20:00:00Z",
      "updated_at": "2025-09-03T20:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 5,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### Get Single Post

Retrieve a specific post by ID.

**Endpoint**: `GET /api/v1/posts/:id`

**Path Parameters**:
- `id` (required): Post ID

**Response** (200 OK):
```json
{
  "post": {
    "id": 1,
    "title": "Welcome to Our Blog Platform",
    "content": "This is our first blog post...",
    "author_id": 1,
    "author": {
      "id": 1,
      "username": "admin"
    },
    "category": "Announcement",
    "tags": ["welcome", "first-post"],
    "image_url": "/static/images/blog1.jpg",
    "view_count": 43,
    "is_published": true,
    "created_at": "2025-09-03T20:00:00Z",
    "updated_at": "2025-09-03T20:00:00Z"
  }
}
```

**Error Response** (404 Not Found):
```json
{
  "error": "Post not found"
}
```

---

### Create Post

Create a new blog post. **Requires authentication**.

**Endpoint**: `POST /api/v1/posts`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "title": "string (required, 3-200 chars)",
  "content": "string (required)",
  "category": "string (optional)",
  "tags": ["string"] (optional)
}
```

**Response** (201 Created):
```json
{
  "message": "Post created successfully",
  "post": {
    "id": 3,
    "title": "My New Post",
    "content": "Post content here...",
    "author_id": 1,
    "category": "Tutorial",
    "tags": ["golang", "gin"],
    "view_count": 0,
    "is_published": true,
    "created_at": "2025-10-03T20:00:00Z",
    "updated_at": "2025-10-03T20:00:00Z"
  }
}
```

---

### Update Post

Update an existing post. **Requires authentication and ownership**.

**Endpoint**: `PUT /api/v1/posts/:id`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body** (all fields optional):
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

**Response** (200 OK):
```json
{
  "message": "Post updated successfully",
  "post": {
    "id": 3,
    "title": "Updated Title",
    "content": "Updated content...",
    "updated_at": "2025-10-03T21:00:00Z"
  }
}
```

**Error Response** (403 Forbidden):
```json
{
  "error": "You don't have permission to update this post"
}
```

---

### Delete Post

Delete a post. **Requires authentication and ownership**.

**Endpoint**: `DELETE /api/v1/posts/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "message": "Post deleted successfully"
}
```

---

## 💬 Comments Endpoints

### Get Comments

Get all comments for a specific post.

**Endpoint**: `GET /api/v1/posts/:id/comments`

**Path Parameters**:
- `id` (required): Post ID

**Response** (200 OK):
```json
{
  "comments": [
    {
      "id": 1,
      "post_id": 1,
      "author_id": 2,
      "author": {
        "id": 2,
        "username": "user"
      },
      "content": "Great post!",
      "parent_id": null,
      "created_at": "2025-09-28T20:00:00Z",
      "updated_at": "2025-09-28T20:00:00Z"
    }
  ],
  "total": 1
}
```

---

### Create Comment

Add a comment to a post. **Requires authentication**.

**Endpoint**: `POST /api/v1/posts/comments`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "post_id": 1,
  "content": "string (required, 1-1000 chars)",
  "parent_id": null (optional, for nested comments)
}
```

**Response** (201 Created):
```json
{
  "message": "Comment created successfully",
  "comment": {
    "id": 3,
    "post_id": 1,
    "author_id": 2,
    "content": "This is a great tutorial!",
    "parent_id": null,
    "created_at": "2025-10-03T20:00:00Z",
    "updated_at": "2025-10-03T20:00:00Z"
  }
}
```

---

### Delete Comment

Delete a comment. **Requires authentication and ownership**.

**Endpoint**: `DELETE /api/v1/posts/comments/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "message": "Comment deleted successfully"
}
```

---

## 📁 Upload Endpoints

### Upload Image

Upload an image file. **Requires authentication**.

**Endpoint**: `POST /api/v1/upload/image`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data**:
- `image`: File (required)
  - Allowed types: JPEG, PNG, GIF, WebP
  - Max size: 5MB

**Response** (200 OK):
```json
{
  "message": "File uploaded successfully",
  "filename": "1696358400_myimage.jpg",
  "size": 245832,
  "url": "/uploads/1696358400_myimage.jpg"
}
```

**Error Responses**:

Invalid file type (400):
```json
{
  "error": "Invalid file type. Only images are allowed"
}
```

File too large (400):
```json
{
  "error": "File too large. Maximum size is 5MB"
}
```

---

## 👤 User Endpoints

### Get Profile

Get current user's profile. **Requires authentication**.

**Endpoint**: `GET /api/v1/user/profile`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

## 🏥 Health Check

### API v1 Health

**Endpoint**: `GET /api/v1/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "version": "1.0",
  "time": "2025-10-03T20:00:00Z",
  "request_id": "req-1696358400000000000"
}
```

---

### API v2 Health

**Endpoint**: `GET /api/v2/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "version": "2.0",
  "features": [
    "enhanced-pagination",
    "advanced-filtering",
    "real-time-updates"
  ],
  "time": "2025-10-03T20:00:00Z"
}
```

---

## 🔧 Admin Endpoints

### Admin Dashboard Stats

Get platform statistics. **Requires admin role**.

**Endpoint**: `GET /admin/api/stats`

**Headers**:
```
Authorization: Bearer admin-token-456
```

**Response** (200 OK):
```json
{
  "total_users": 2,
  "total_posts": 2,
  "total_comments": 2,
  "uptime": "100%"
}
```

---

## ⚠️ Error Responses

### Common Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

### Error Response Format

```json
{
  "error": "Error message",
  "details": "Detailed error information (optional)"
}
```

### Validation Error Example

```json
{
  "error": "Invalid request data",
  "details": "Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

---

## 🔒 Rate Limiting

- **API v1**: 100 requests per minute per IP
- **API v2**: 150 requests per minute per IP

When rate limit is exceeded:

**Response** (429 Too Many Requests):
```json
{
  "error": "Rate limit exceeded",
  "retry_after": "60 seconds"
}
```

---

## 📝 Notes

1. All timestamps are in RFC3339 format (ISO 8601)
2. The API uses in-memory storage - data resets on server restart
3. Tokens are simplified for demo purposes - use JWT in production
4. CORS is enabled for all origins - restrict in production
5. File uploads are stored locally - use cloud storage in production

---

## 🧪 Testing with cURL

See [QUICKSTART.md](QUICKSTART.md) for detailed cURL examples.

---

For more information, see the [README.md](README.md).
