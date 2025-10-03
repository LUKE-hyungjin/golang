# Lesson 22: 통합 테스트 (Integration Testing) 🔄

> 실제 데이터베이스와 전체 애플리케이션 스택을 테스트하는 통합 테스트 구현

## 📌 이번 레슨에서 배우는 내용

통합 테스트는 시스템의 여러 구성 요소가 함께 올바르게 작동하는지 검증합니다. 데이터베이스, 라우터, 미들웨어를 모두 포함한 실제 환경과 유사한 테스트 환경을 구축합니다.

### 핵심 학습 목표
- ✅ 테스트 데이터베이스 설정
- ✅ 트랜잭션 기반 테스트 격리
- ✅ 테스트 픽스처 관리
- ✅ E2E 시나리오 테스트
- ✅ 동시성 테스트
- ✅ Test Suite 활용

## 🏗 통합 테스트 아키텍처

### 테스트 환경 구조
```
┌──────────────────┐
│   Test Suite     │
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  Test Database   │ → In-Memory SQLite
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  Transaction     │ → Rollback after test
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│   Application    │ → Full stack
└──────────────────┘
```

### 테스트 격리 전략
```
1. Setup Phase
   - Create test database
   - Run migrations
   - Load fixtures

2. Test Execution
   - Begin transaction
   - Run test
   - Capture results

3. Cleanup Phase
   - Rollback transaction
   - Reset state
   - Close connections
```

## 🛠 구현된 기능

### 1. **테스트 데이터베이스**
- In-memory SQLite
- 자동 마이그레이션
- 트랜잭션 지원
- 테스트별 격리

### 2. **테스트 서버**
- 전체 애플리케이션 스택
- 실제 라우팅
- 미들웨어 포함
- 의존성 주입

### 3. **테스트 픽스처**
- 시드 데이터 관리
- 재사용 가능한 테스트 데이터
- 관계형 데이터 설정
- 일관된 테스트 환경

### 4. **시나리오 테스트**
- 사용자 플로우
- CRUD 작업 체인
- 트랜잭션 테스트
- 동시성 처리

### 5. **Test Suite**
- Setup/Teardown 훅
- 테스트 격리
- 공유 리소스 관리
- 병렬 실행 지원

## 💻 실습 가이드

### 1. 설치 및 설정
```bash
cd gin/22
go mod init integration-test
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u github.com/stretchr/testify
go get -u github.com/mattn/go-sqlite3

# 테스트 파일 생성
touch main_integration_test.go
```

### 2. 테스트 실행
```bash
# 모든 통합 테스트 실행
go test -v -tags=integration

# 특정 테스트만 실행
go test -v -run TestCreateUser_Integration

# Test Suite 실행
go test -v -run TestBlogIntegrationSuite

# 벤치마크 실행
go test -bench=Integration

# 레이스 컨디션 체크
go test -race -tags=integration

# 커버리지 측정
go test -cover -tags=integration
```

## 🎯 주요 테스트 예제

### 1. 테스트 서버 설정
```go
type TestServer struct {
    Router  *gin.Engine
    DB      *TestDatabase
    Service *BlogService
    Handler *BlogHandler
}

func NewTestServer() (*TestServer, error) {
    gin.SetMode(gin.TestMode)

    // In-memory database
    db, err := NewTestDatabase()
    if err != nil {
        return nil, err
    }

    service := NewBlogService(db.GetDB())
    handler := NewBlogHandler(service)
    router := SetupRouter(handler)

    return &TestServer{
        Router:  router,
        DB:      db,
        Service: service,
        Handler: handler,
    }, nil
}
```

### 2. 트랜잭션 기반 테스트
```go
func TestWithTransaction(t *testing.T) {
    server, err := NewTestServer()
    require.NoError(t, err)
    defer server.Cleanup()

    // Begin transaction
    server.DB.Begin()
    defer server.DB.Rollback()  // Always rollback

    // Test operations
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    err = server.DB.Create(user).Error
    assert.NoError(t, err)

    // Verify in same transaction
    var found User
    err = server.DB.First(&found, user.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, user.Username, found.Username)

    // Rollback happens automatically
}
```

### 3. E2E 시나리오 테스트
```go
func TestUserPostFlow_Integration(t *testing.T) {
    server, err := NewTestServer()
    require.NoError(t, err)
    defer server.Cleanup()

    // Step 1: Create user
    userReq := map[string]string{
        "username": "flowuser",
        "email":    "flow@example.com",
        "password": "password123",
    }
    userBody, _ := json.Marshal(userReq)

    w := performRequest(server.Router, "POST", "/api/v1/users",
        bytes.NewBuffer(userBody))
    assert.Equal(t, http.StatusCreated, w.Code)

    var user User
    json.Unmarshal(w.Body.Bytes(), &user)

    // Step 2: Create post
    postReq := map[string]interface{}{
        "title":   "Test Post",
        "content": "Test Content",
        "user_id": user.ID,
    }
    postBody, _ := json.Marshal(postReq)

    w = performRequest(server.Router, "POST", "/api/v1/posts",
        bytes.NewBuffer(postBody))
    assert.Equal(t, http.StatusCreated, w.Code)

    var post Post
    json.Unmarshal(w.Body.Bytes(), &post)

    // Step 3: Add comment
    commentReq := map[string]interface{}{
        "content": "Great post!",
        "post_id": post.ID,
        "user_id": user.ID,
    }
    commentBody, _ := json.Marshal(commentReq)

    w = performRequest(server.Router, "POST", "/api/v1/comments",
        bytes.NewBuffer(commentBody))
    assert.Equal(t, http.StatusCreated, w.Code)

    // Step 4: Verify complete post
    w = performRequest(server.Router, "GET",
        fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
    assert.Equal(t, http.StatusOK, w.Code)

    var fullPost Post
    json.Unmarshal(w.Body.Bytes(), &fullPost)
    assert.Len(t, fullPost.Comments, 1)
    assert.NotNil(t, fullPost.User)
}
```

### 4. 테스트 픽스처 관리
```go
type TestFixtures struct {
    Users    []User
    Posts    []Post
    Comments []Comment
    Tags     []Tag
}

func LoadTestFixtures(db *gorm.DB) *TestFixtures {
    fixtures := &TestFixtures{
        Users: []User{
            {Username: "admin", Email: "admin@example.com"},
            {Username: "editor", Email: "editor@example.com"},
            {Username: "viewer", Email: "viewer@example.com"},
        },
        Posts: []Post{
            {Title: "Welcome", Content: "Hello", UserID: 1},
            {Title: "Tutorial", Content: "Learn", UserID: 1},
        },
    }

    // Load into database
    for _, user := range fixtures.Users {
        db.Create(&user)
    }
    for _, post := range fixtures.Posts {
        db.Create(&post)
    }

    return fixtures
}
```

### 5. Test Suite 구성
```go
type BlogIntegrationSuite struct {
    suite.Suite
    server *TestServer
}

func (suite *BlogIntegrationSuite) SetupSuite() {
    server, err := NewTestServer()
    suite.Require().NoError(err)
    suite.server = server
}

func (suite *BlogIntegrationSuite) TearDownSuite() {
    suite.server.Cleanup()
}

func (suite *BlogIntegrationSuite) SetupTest() {
    // Start transaction for each test
    suite.server.DB.Begin()
}

func (suite *BlogIntegrationSuite) TearDownTest() {
    // Rollback after each test
    suite.server.DB.Rollback()
}

func (suite *BlogIntegrationSuite) TestCompleteScenario() {
    // Seed data
    suite.server.SeedTestData()

    // Test operations
    req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
    w := httptest.NewRecorder()
    suite.server.Router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    posts := response["posts"].([]interface{})
    suite.NotEmpty(posts)
}
```

### 6. 동시성 테스트
```go
func TestConcurrentRequests(t *testing.T) {
    server, err := NewTestServer()
    require.NoError(t, err)
    defer server.Cleanup()

    done := make(chan bool, 10)
    errors := make(chan error, 10)

    // Create 10 concurrent requests
    for i := 0; i < 10; i++ {
        go func(index int) {
            user := map[string]string{
                "username": fmt.Sprintf("user%d", index),
                "email":    fmt.Sprintf("user%d@example.com", index),
                "password": "password123",
            }
            jsonBody, _ := json.Marshal(user)

            w := performRequest(server.Router, "POST", "/api/v1/users",
                bytes.NewBuffer(jsonBody))

            if w.Code != http.StatusCreated {
                errors <- fmt.Errorf("Request %d failed", index)
            }
            done <- true
        }(i)
    }

    // Wait for completion
    for i := 0; i < 10; i++ {
        <-done
    }

    close(errors)
    for err := range errors {
        t.Error(err)
    }

    // Verify all users created
    var count int64
    server.DB.Model(&User{}).Count(&count)
    assert.Equal(t, int64(10), count)
}
```

## 🔍 데이터베이스 테스트 패턴

### 테스트 데이터베이스 설정
```go
type TestDatabase struct {
    *Database
    tx *gorm.DB
}

func NewTestDatabase() (*TestDatabase, error) {
    // In-memory database
    config := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    }

    db, err := gorm.Open(sqlite.Open(":memory:"), config)
    if err != nil {
        return nil, err
    }

    // Run migrations
    db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Tag{})

    return &TestDatabase{Database: &Database{DB: db}}, nil
}

func (tdb *TestDatabase) Begin() {
    tdb.tx = tdb.DB.Begin()
}

func (tdb *TestDatabase) Rollback() {
    if tdb.tx != nil {
        tdb.tx.Rollback()
    }
}
```

### 마이그레이션 테스트
```go
func TestMigrations(t *testing.T) {
    db, err := NewTestDatabase()
    require.NoError(t, err)

    // Check tables exist
    var tableNames []string
    db.Raw("SELECT name FROM sqlite_master WHERE type='table'").
        Pluck("name", &tableNames)

    assert.Contains(t, tableNames, "users")
    assert.Contains(t, tableNames, "posts")
    assert.Contains(t, tableNames, "comments")
    assert.Contains(t, tableNames, "tags")
}
```

## 📊 테스트 환경 관리

### 환경 변수 설정
```bash
# Test environment
export GIN_MODE=test
export DB_FILE=:memory:
export LOG_LEVEL=error

# Run tests
go test -v
```

### Docker를 사용한 테스트
```dockerfile
# Dockerfile.test
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go test -v ./...
```

```bash
# Docker로 테스트 실행
docker build -f Dockerfile.test -t test-runner .
docker run --rm test-runner
```

## 🎨 테스트 데이터 생성

### Faker 사용
```go
import "github.com/bxcodec/faker/v3"

func GenerateTestUser() *User {
    user := &User{}
    faker.FakeData(user)
    return user
}

func GenerateTestPosts(count int) []Post {
    var posts []Post
    for i := 0; i < count; i++ {
        var post Post
        faker.FakeData(&post)
        posts = append(posts, post)
    }
    return posts
}
```

### Factory 패턴
```go
type UserFactory struct {
    defaults User
}

func NewUserFactory() *UserFactory {
    return &UserFactory{
        defaults: User{
            Username: "testuser",
            Email:    "test@example.com",
            Password: "password123",
        },
    }
}

func (f *UserFactory) Build() *User {
    user := f.defaults
    return &user
}

func (f *UserFactory) WithUsername(username string) *UserFactory {
    f.defaults.Username = username
    return f
}
```

## 🚀 성능 테스트

### 부하 테스트
```go
func TestHighLoad(t *testing.T) {
    server, _ := NewTestServer()
    defer server.Cleanup()

    start := time.Now()
    requests := 1000
    done := make(chan bool, requests)

    for i := 0; i < requests; i++ {
        go func(index int) {
            w := performRequest(server.Router, "GET", "/health", nil)
            assert.Equal(t, 200, w.Code)
            done <- true
        }(i)
    }

    for i := 0; i < requests; i++ {
        <-done
    }

    elapsed := time.Since(start)
    rps := float64(requests) / elapsed.Seconds()

    t.Logf("Processed %d requests in %v", requests, elapsed)
    t.Logf("Requests per second: %.2f", rps)
    assert.Greater(t, rps, 100.0) // At least 100 RPS
}
```

### 메모리 사용량 테스트
```go
func TestMemoryUsage(t *testing.T) {
    server, _ := NewTestServer()
    defer server.Cleanup()

    // Initial memory
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    initialMem := m.Alloc

    // Perform operations
    for i := 0; i < 100; i++ {
        user := &User{
            Username: fmt.Sprintf("user%d", i),
            Email:    fmt.Sprintf("user%d@example.com", i),
        }
        server.DB.Create(user)
    }

    // Final memory
    runtime.ReadMemStats(&m)
    finalMem := m.Alloc

    memUsed := finalMem - initialMem
    t.Logf("Memory used: %d bytes", memUsed)
    assert.Less(t, memUsed, uint64(10*1024*1024)) // Less than 10MB
}
```

## 📝 베스트 프랙티스

### 1. **테스트 격리**
```go
// ✅ Good: Each test is independent
func TestUserCreate(t *testing.T) {
    server := NewTestServer()
    defer server.Cleanup()
    // Test logic
}

// ❌ Bad: Tests depend on shared state
var globalServer *TestServer

func TestUserCreate(t *testing.T) {
    // Uses shared server
}
```

### 2. **명확한 Assertion**
```go
// ✅ Good: Specific assertions
assert.Equal(t, http.StatusCreated, w.Code)
assert.Equal(t, "testuser", user.Username)
assert.NotZero(t, user.ID)

// ❌ Bad: Generic assertion
assert.True(t, w.Code == 201 && user.Username == "testuser")
```

### 3. **테스트 데이터 정리**
```go
// ✅ Good: Cleanup after test
func TestSomething(t *testing.T) {
    server := NewTestServer()
    defer server.Cleanup()

    // Or use transaction
    server.DB.Begin()
    defer server.DB.Rollback()
}
```

### 4. **타임아웃 설정**
```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Test with timeout
    select {
    case <-performOperation(ctx):
        // Success
    case <-ctx.Done():
        t.Error("Test timed out")
    }
}
```

## 🔒 체크리스트

- [ ] 테스트 데이터베이스 설정
- [ ] 트랜잭션 격리
- [ ] 픽스처 관리
- [ ] E2E 시나리오
- [ ] 동시성 테스트
- [ ] 성능 벤치마크
- [ ] 메모리 프로파일링
- [ ] 타임아웃 처리
- [ ] 정리(Cleanup)

## 📚 추가 학습 자료

- [GORM Testing](https://gorm.io/docs/connecting_to_the_database.html#SQLite)
- [testify Suite](https://github.com/stretchr/testify#suite-package)
- [Go Integration Testing](https://go.dev/doc/tutorial/add-a-test)
- [Database Testing Best Practices](https://www.alexedwards.net/blog/organising-database-access)

## 🎯 다음 레슨 예고

**Lesson 23: 린팅과 포맷팅**
- golangci-lint 설정
- 코드 품질 검사
- 자동 포맷팅
- CI/CD 통합
- 커스텀 린터 규칙

통합 테스트로 시스템 전체의 안정성을 보장하세요! 🔄