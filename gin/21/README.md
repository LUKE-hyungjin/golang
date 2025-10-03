# Lesson 21: 핸들러 유닛 테스트 (httptest) 🧪

> Go의 httptest 패키지로 Gin 핸들러를 체계적으로 테스트하기

## 📌 이번 레슨에서 배우는 내용

테스트는 안정적인 웹 애플리케이션 개발의 핵심입니다. Go의 `httptest` 패키지와 `testify`를 활용하면 Gin 핸들러를 효과적으로 테스트할 수 있습니다. 이번 레슨에서는 유닛 테스트의 모든 패턴을 다룹니다.

### 핵심 학습 목표
- ✅ httptest를 사용한 HTTP 요청/응답 테스트
- ✅ Mock 객체를 활용한 의존성 격리
- ✅ 테스트 헬퍼 함수 작성
- ✅ Table-driven 테스트
- ✅ Test Suite 구성
- ✅ 인증이 필요한 엔드포인트 테스트

## 🏗 테스트 아키텍처

### 테스트 레이어
```
┌──────────────────┐
│   Test Cases     │
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  httptest.Recorder │ → 응답 캡처
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│   Gin Router     │ → 라우팅 처리
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│   Mock Repository │ → 데이터 격리
└──────────────────┘
```

### 테스트 구성 요소
```
1. httptest.ResponseRecorder: HTTP 응답 캡처
2. http.NewRequest: 테스트 요청 생성
3. Mock Objects: 외부 의존성 격리
4. Assertions: 결과 검증
5. Test Fixtures: 테스트 데이터 준비
```

## 🛠 구현된 테스트 패턴

### 1. **기본 유닛 테스트**
- GET 요청 테스트
- POST 요청 테스트
- PUT/DELETE 테스트
- 상태 코드 검증
- 응답 바디 검증

### 2. **Mock 객체 활용**
- Repository 인터페이스 정의
- MockRepository 구현
- 의존성 주입
- 테스트 데이터 격리

### 3. **인증 테스트**
- 미들웨어 테스트
- 토큰 검증
- 권한 체크
- 인증 실패 시나리오

### 4. **파일 업로드 테스트**
- Multipart form 생성
- 파일 크기 검증
- 파일 타입 체크
- 업로드 성공/실패

### 5. **Table-driven 테스트**
- 다양한 입력 케이스
- 경계값 테스트
- 유효성 검사
- 에러 케이스

### 6. **Test Suite**
- Setup/Teardown
- 테스트 그룹화
- 공통 설정 재사용
- 통합 플로우 테스트

## 💻 실습 가이드

### 1. 설치 및 설정
```bash
cd gin/21
go mod init test-example
go get -u github.com/gin-gonic/gin
go get -u github.com/stretchr/testify

# main_test.go 파일 생성
touch main_test.go
```

### 2. 테스트 파일 구조
```go
// main_test.go
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
    // Arrange (준비)
    router := SetupTestRouter()

    // Act (실행)
    w := performRequest(router, "GET", "/users/1", nil)

    // Assert (검증)
    assert.Equal(t, 200, w.Code)
}
```

### 3. 테스트 실행
```bash
# 모든 테스트 실행
go test

# 상세 출력
go test -v

# 특정 테스트만 실행
go test -run TestGetUser

# 커버리지 측정
go test -cover

# 커버리지 리포트 생성
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# 벤치마크 실행
go test -bench=.

# 레이스 컨디션 체크
go test -race
```

## 🎯 주요 테스트 예제

### 1. GET 요청 테스트
```go
func TestGetUser_Success(t *testing.T) {
    // Setup
    repo := NewMockUserRepository()
    service := NewUserService(repo)
    handler := NewUserHandler(service)
    router := SetupRouter(handler)

    // Test data
    testUser := &User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
    repo.Create(testUser)

    // Perform request
    w := performRequest(router, "GET", "/users/1", nil)

    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)

    var response User
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, testUser.Username, response.Username)
}
```

### 2. POST 요청 테스트
```go
func TestCreateUser_Success(t *testing.T) {
    router := SetupTestRouter()

    user := User{
        Username: "newuser",
        Email:    "new@example.com",
    }
    jsonBody, _ := json.Marshal(user)

    w := performRequest(router, "POST", "/users",
        bytes.NewBuffer(jsonBody))

    assert.Equal(t, http.StatusCreated, w.Code)

    var response User
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotZero(t, response.ID)
}
```

### 3. 인증이 필요한 요청 테스트
```go
func TestUpdateUser_WithAuth(t *testing.T) {
    router := SetupTestRouter()

    updateData := User{
        Username: "updated",
        Email:    "updated@example.com",
    }
    jsonBody, _ := json.Marshal(updateData)

    // With valid token
    headers := map[string]string{
        "Authorization": "Bearer valid-token",
    }
    w := performRequestWithHeaders(router, "PUT", "/users/1",
        bytes.NewBuffer(jsonBody), headers)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
    router := SetupTestRouter()

    // Without token
    w := performRequest(router, "PUT", "/users/1", nil)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

### 4. 파일 업로드 테스트
```go
func TestUploadAvatar_Success(t *testing.T) {
    router := SetupTestRouter()

    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // Add file field
    part, _ := writer.CreateFormFile("avatar", "test.jpg")
    part.Write([]byte("fake-image-data"))
    writer.Close()

    // Perform request
    req, _ := http.NewRequest("POST", "/users/1/avatar", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Bearer valid-token")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 5. Table-driven 테스트
```go
func TestUserValidation_TableDriven(t *testing.T) {
    tests := []struct {
        name         string
        user         User
        expectedCode int
    }{
        {
            name: "Valid user",
            user: User{
                Username: "validuser",
                Email:    "valid@example.com",
            },
            expectedCode: http.StatusCreated,
        },
        {
            name: "Username too short",
            user: User{
                Username: "ab",
                Email:    "valid@example.com",
            },
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Invalid email",
            user: User{
                Username: "validuser",
                Email:    "invalid-email",
            },
            expectedCode: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jsonBody, _ := json.Marshal(tt.user)
            w := performRequest(router, "POST", "/users",
                bytes.NewBuffer(jsonBody))
            assert.Equal(t, tt.expectedCode, w.Code)
        })
    }
}
```

### 6. Test Suite 사용
```go
type UserHandlerTestSuite struct {
    suite.Suite
    router  *gin.Engine
    repo    *MockUserRepository
}

func (suite *UserHandlerTestSuite) SetupTest() {
    suite.repo = NewMockUserRepository()
    service := NewUserService(suite.repo)
    handler := NewUserHandler(service)
    suite.router = SetupRouter(handler)
}

func (suite *UserHandlerTestSuite) TestUserCRUDFlow() {
    // Create
    user := User{Username: "test", Email: "test@example.com"}
    jsonBody, _ := json.Marshal(user)

    w := performRequest(suite.router, "POST", "/users",
        bytes.NewBuffer(jsonBody))
    suite.Equal(http.StatusCreated, w.Code)

    var created User
    json.Unmarshal(w.Body.Bytes(), &created)

    // Read
    w = performRequest(suite.router, "GET",
        fmt.Sprintf("/users/%d", created.ID), nil)
    suite.Equal(http.StatusOK, w.Code)

    // Update
    // ... update logic

    // Delete
    // ... delete logic
}

func TestUserHandlerSuite(t *testing.T) {
    suite.Run(t, new(UserHandlerTestSuite))
}
```

## 🔍 Mock 객체 패턴

### Repository 인터페이스
```go
type UserRepository interface {
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}
```

### Mock Repository 구현
```go
type MockUserRepository struct {
    users map[uint]*User
}

func (r *MockUserRepository) FindByID(id uint) (*User, error) {
    user, exists := r.users[id]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}

func (r *MockUserRepository) Create(user *User) error {
    user.ID = uint(len(r.users) + 1)
    user.CreatedAt = time.Now()
    r.users[user.ID] = user
    return nil
}
```

## 📊 테스트 커버리지

### 커버리지 목표
- 핵심 비즈니스 로직: 90%+
- HTTP 핸들러: 80%+
- 유틸리티 함수: 70%+
- 전체 평균: 80%+

### 커버리지 측정
```bash
# 커버리지 측정
go test -cover

# 상세 리포트
go test -coverprofile=coverage.out
go tool cover -func=coverage.out

# HTML 리포트
go tool cover -html=coverage.out -o coverage.html

# 특정 패키지만
go test -cover ./handlers/...
```

## 🎨 테스트 헬퍼 함수

### HTTP 요청 헬퍼
```go
func performRequest(r http.Handler, method, path string,
    body io.Reader) *httptest.ResponseRecorder {
    req, _ := http.NewRequest(method, path, body)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    return w
}

func performRequestWithHeaders(r http.Handler, method, path string,
    body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
    req, _ := http.NewRequest(method, path, body)
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    return w
}
```

### JSON 헬퍼
```go
func toJSON(v interface{}) []byte {
    data, _ := json.Marshal(v)
    return data
}

func fromJSON(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}
```

## 🚀 벤치마크 테스트

### 성능 측정
```go
func BenchmarkGetUser(b *testing.B) {
    router := SetupTestRouter()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := performRequest(router, "GET", "/users/1", nil)
        if w.Code != http.StatusOK {
            b.Errorf("Expected 200, got %d", w.Code)
        }
    }
}

func BenchmarkCreateUser(b *testing.B) {
    router := SetupTestRouter()
    user := User{Username: "bench", Email: "bench@example.com"}
    jsonBody, _ := json.Marshal(user)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := performRequest(router, "POST", "/users",
            bytes.NewBuffer(jsonBody))
    }
}
```

### 벤치마크 실행
```bash
# 모든 벤치마크
go test -bench=.

# 특정 벤치마크
go test -bench=BenchmarkGetUser

# 메모리 할당 측정
go test -bench=. -benchmem

# 시간 지정
go test -bench=. -benchtime=10s
```

## 📝 베스트 프랙티스

### 1. **테스트 격리**
```go
// ❌ Bad: 테스트 간 의존성
func TestCreateUser(t *testing.T) {
    // Creates user with ID 1
}

func TestGetUser(t *testing.T) {
    // Assumes user with ID 1 exists
}

// ✅ Good: 독립적인 테스트
func TestGetUser(t *testing.T) {
    // Setup own test data
    user := createTestUser()
    // Test with that user
}
```

### 2. **명확한 테스트 이름**
```go
// ❌ Bad
func TestUser(t *testing.T) {}

// ✅ Good
func TestGetUser_WhenUserExists_ReturnsUser(t *testing.T) {}
func TestGetUser_WhenUserNotFound_Returns404(t *testing.T) {}
```

### 3. **AAA 패턴 사용**
```go
func TestCreateUser(t *testing.T) {
    // Arrange (준비)
    router := SetupRouter()
    user := User{Username: "test"}

    // Act (실행)
    w := performRequest(router, "POST", "/users", toJSON(user))

    // Assert (검증)
    assert.Equal(t, 201, w.Code)
}
```

### 4. **테스트 데이터 빌더**
```go
type UserBuilder struct {
    user User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: User{
            Username: "default",
            Email:    "default@example.com",
        },
    }
}

func (b *UserBuilder) WithUsername(name string) *UserBuilder {
    b.user.Username = name
    return b
}

func (b *UserBuilder) Build() User {
    return b.user
}

// Usage
user := NewUserBuilder().
    WithUsername("custom").
    WithEmail("custom@example.com").
    Build()
```

## 🔒 테스트 체크리스트

- [ ] 모든 엔드포인트 테스트
- [ ] 성공 케이스 테스트
- [ ] 실패 케이스 테스트
- [ ] 경계값 테스트
- [ ] 인증/권한 테스트
- [ ] 유효성 검사 테스트
- [ ] 에러 핸들링 테스트
- [ ] 동시성 테스트
- [ ] 성능 벤치마크

## 📚 추가 학습 자료

- [Go testing package](https://golang.org/pkg/testing/)
- [httptest package](https://golang.org/pkg/net/http/httptest/)
- [testify](https://github.com/stretchr/testify)
- [gomock](https://github.com/golang/mock)
- [Go Test Patterns](https://github.com/gotestyourself/gotest.tools)

## 🎯 다음 레슨 예고

**Lesson 22: 통합 테스트**
- 테스트 데이터베이스 설정
- 트랜잭션 롤백
- E2E 테스트
- API 시나리오 테스트
- 테스트 환경 구성

체계적인 테스트로 안정적인 애플리케이션을 만드세요! 🧪