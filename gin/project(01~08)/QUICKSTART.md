# 🚀 Quick Start Guide

Get the Blog Community Platform up and running in minutes!

## ⚡ 3-Step Quick Start

### 1. Install Dependencies
```bash
cd /Users/ihyeongjin/dev/golang/gin/project
make install
```

### 2. Run the Application
```bash
make run
```

### 3. Open Your Browser
```
http://localhost:8080
```

You should see the welcome page! 🎉

## 🧪 Testing the Application

### Web Interface Testing

1. **Home Page**: Navigate to `http://localhost:8080`
   - View feature overview
   - See API documentation
   - Check test credentials

2. **Browse Posts**: Go to `http://localhost:8080/posts`
   - See existing blog posts
   - Click on a post to read it

3. **Login**: Visit `http://localhost:8080/login`
   - Use test credentials:
     - Username: `admin`
     - Password: `password123`

4. **Register**: Try `http://localhost:8080/register`
   - Create a new account
   - Test the registration flow

### API Testing with cURL

#### 1. Health Check
```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0",
  "time": "2025-10-03T...",
  "request_id": "req-..."
}
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

Expected response:
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

#### 3. Get Posts
```bash
curl http://localhost:8080/api/v1/posts
```

#### 4. Create a Post (Requires Authentication)
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token-456" \
  -d '{
    "title": "My First Post via API",
    "content": "This is a test post created through the API!",
    "category": "Test",
    "tags": ["test", "api"]
  }'
```

#### 5. Add a Comment
```bash
curl -X POST http://localhost:8080/api/v1/posts/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer valid-token-123" \
  -d '{
    "post_id": 1,
    "content": "Great post! Thanks for sharing."
  }'
```

#### 6. Upload an Image
```bash
curl -X POST http://localhost:8080/api/v1/upload/image \
  -H "Authorization: Bearer admin-token-456" \
  -F "image=@/path/to/your/image.jpg"
```

### API Testing with Postman

1. Import these endpoints into Postman:
   - Base URL: `http://localhost:8080`

2. Create a Postman Collection with:

**Authentication**
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

**Posts**
- `GET /api/v1/posts`
- `GET /api/v1/posts/:id`
- `POST /api/v1/posts` (with Bearer token)
- `PUT /api/v1/posts/:id` (with Bearer token)
- `DELETE /api/v1/posts/:id` (with Bearer token)

**Comments**
- `GET /api/v1/posts/:id/comments`
- `POST /api/v1/posts/comments` (with Bearer token)

**Upload**
- `POST /api/v1/upload/image` (with Bearer token)

3. Set up environment variables:
   - `base_url`: `http://localhost:8080`
   - `admin_token`: `admin-token-456`
   - `user_token`: `valid-token-123`

## 📊 Feature Checklist

Test all these features:

- [ ] View home page
- [ ] Browse posts listing
- [ ] View single post with comments
- [ ] Register new user account
- [ ] Login with test credentials
- [ ] Access admin dashboard (admin role required)
- [ ] Create new post via API
- [ ] Update existing post
- [ ] Delete post
- [ ] Add comment to post
- [ ] Upload image file
- [ ] Test rate limiting (make 100+ requests quickly)
- [ ] Test CORS (from different origin)
- [ ] View 404 page (visit non-existent route)

## 🎯 Key Endpoints Summary

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/` | GET | No | Home page |
| `/posts` | GET | No | Posts listing |
| `/posts/:id` | GET | No | Single post |
| `/login` | GET | No | Login page |
| `/register` | GET | No | Registration page |
| `/admin/dashboard` | GET | Yes (Admin) | Admin dashboard |
| `/api/v1/health` | GET | No | Health check |
| `/api/v1/auth/register` | POST | No | Register user |
| `/api/v1/auth/login` | POST | No | Login user |
| `/api/v1/posts` | GET | No | Get all posts |
| `/api/v1/posts/:id` | GET | No | Get single post |
| `/api/v1/posts` | POST | Yes | Create post |
| `/api/v1/posts/:id` | PUT | Yes | Update post |
| `/api/v1/posts/:id` | DELETE | Yes | Delete post |
| `/api/v1/posts/:id/comments` | GET | No | Get comments |
| `/api/v1/posts/comments` | POST | Yes | Add comment |
| `/api/v1/upload/image` | POST | Yes | Upload image |
| `/api/v2/health` | GET | No | Enhanced health check |

## 🔑 Test Credentials

### Admin User
```
Username: admin
Password: password123
Token: admin-token-456
```

### Regular User
```
Username: user
Password: password123
Token: valid-token-123
```

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Dependencies Not Found
```bash
# Clean and reinstall
go clean -modcache
make install
```

### Templates Not Loading
Make sure you're running from the project root:
```bash
cd /Users/ihyeongjin/dev/golang/gin/project
make run
```

## 📖 Next Steps

1. ✅ Complete the Quick Start
2. 📚 Read the [full README.md](README.md)
3. 🔍 Explore the code in `internal/` directory
4. 🎨 Customize the templates in `web/templates/`
5. 💻 Modify the styles in `web/static/css/style.css`
6. 🚀 Add your own features!

## 💡 Tips

- Check the console output for request logs
- Use browser DevTools to inspect API calls
- The database is in-memory, so data resets on restart
- All uploaded files go to `uploads/` directory
- Static files are served from `web/static/`

---

**Happy Coding! 🎉**

For detailed information, see [README.md](README.md)
