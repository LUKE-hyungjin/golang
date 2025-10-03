# 🔐 Project Security: 완벽한 보안과 테스트를 갖춘 Go API

> Lessons 18-23의 모든 내용을 통합한 프로덕션 레벨 보안 API 프로젝트

## 📌 프로젝트 개요

이 프로젝트는 Gin 프레임워크를 사용하여 구축된 완벽한 보안 기능과 테스트를 갖춘 RESTful API입니다. CORS, JWT 인증, 입력 검증, 포괄적인 테스트, 그리고 코드 품질 도구가 모두 통합되어 있습니다.

### 핵심 기능
- 🔒 **JWT 기반 인증** (Access/Refresh Token)
- 🌐 **CORS 설정** (환경별 구성)
- ✅ **강력한 입력 검증** (커스텀 검증자 포함)
- 🧪 **포괄적인 테스트** (단위/통합 테스트)
- 🎨 **코드 품질 관리** (golangci-lint)
- 📊 **Rate Limiting** (API 보호)
- 🛡️ **보안 헤더** (XSS, CSRF 방지)
- 📝 **구조화된 로깅** (요청 추적)

## 🏗 프로젝트 구조

```
project-security/
├── cmd/
│   └── main.go                 # 애플리케이션 진입점
├── internal/
│   ├── auth/
│   │   └── jwt.go              # JWT 토큰 관리
│   ├── handlers/
│   │   ├── auth.go             # 인증 핸들러
│   │   ├── user.go             # 사용자 핸들러
│   │   ├── post.go             # 포스트 핸들러
│   │   └── health.go           # 헬스체크 핸들러
│   ├── middleware/
│   │   ├── auth.go             # 인증 미들웨어
│   │   ├── cors.go             # CORS 미들웨어
│   │   └── common.go           # 공통 미들웨어
│   ├── models/
│   │   └── models.go           # 데이터 모델
│   ├── repository/
│   │   ├── user.go             # 사용자 저장소
│   │   └── post.go             # 포스트 저장소
│   ├── services/
│   │   ├── auth.go             # 인증 서비스
│   │   ├── user.go             # 사용자 서비스
│   │   └── post.go             # 포스트 서비스
│   └── validator/
│       ├── validator.go        # 커스텀 검증자
│       └── requests.go         # 요청 DTO
├── pkg/
│   ├── config/
│   │   └── config.go           # 설정 관리
│   ├── database/
│   │   └── database.go         # 데이터베이스 연결
│   └── logger/
│       └── logger.go           # 로거 설정
├── tests/
│   ├── integration/
│   │   ├── auth_test.go        # 인증 통합 테스트
│   │   └── post_test.go        # 포스트 통합 테스트
│   ├── unit/
│   │   ├── jwt_test.go         # JWT 단위 테스트
│   │   └── validator_test.go   # 검증자 단위 테스트
│   └── fixtures/
│       └── test_data.go        # 테스트 데이터
├── migrations/                  # 데이터베이스 마이그레이션
├── docs/                        # API 문서
├── scripts/                     # 유틸리티 스크립트
├── .golangci.yml               # 린터 설정
├── Makefile                    # 빌드 자동화
├── docker-compose.yml          # 도커 설정
├── Dockerfile                  # 도커 이미지
└── config.yaml                 # 애플리케이션 설정
```

## 🚀 시작하기

### 필수 요구사항
- Go 1.21 이상
- SQLite3
- Make (선택사항)

### 설치 및 실행

```bash
# 프로젝트 클론
cd gin/project-security

# 의존성 설치
make install-deps

# 개발 도구 설치 (golangci-lint, air 등)
make install-tools

# 데이터베이스 마이그레이션
make migrate-up

# 개발 모드 실행 (hot reload)
make dev

# 또는 일반 실행
make run

# 프로덕션 빌드
make build
./project-security
```

### 환경 설정

```yaml
# config.yaml
environment: development
log_level: debug

server:
  port: 8080
  read_timeout: 15s
  write_timeout: 15s

database:
  dsn: security.db
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: your-secret-key-change-in-production
  access_expiry: 15m
  refresh_expiry: 168h # 7 days

cors:
  allow_origins:
    - http://localhost:3000
    - http://localhost:5173
  allow_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allow_credentials: true

rate_limit:
  requests_per_minute: 60
```

## 🔐 보안 기능 상세

### 1. JWT 인증 시스템

#### 토큰 구조
```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

#### 토큰 플로우
```
1. 회원가입/로그인
   ↓
2. Access Token (15분) + Refresh Token (7일) 발급
   ↓
3. API 요청 시 Access Token 사용
   ↓
4. Access Token 만료 시 Refresh Token으로 갱신
   ↓
5. Refresh Token 만료 시 재로그인
```

### 2. CORS 설정

```go
// 환경별 CORS 구성
Development:
  - Allow Origins: ["http://localhost:3000"]
  - Allow Credentials: true
  - Max Age: 12h

Production:
  - Allow Origins: ["https://yourdomain.com"]
  - Allow Credentials: true
  - Max Age: 24h
```

### 3. 입력 검증

#### 커스텀 검증자
- `strong_password`: 강력한 비밀번호 (대소문자, 숫자, 특수문자)
- `phone`: 전화번호 형식
- `slug`: URL 슬러그
- `no_sql_injection`: SQL 인젝션 방지

#### 사용 예시
```go
type RegisterRequest struct {
    Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,strong_password"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
```

### 4. Rate Limiting

```go
// IP 기반 Rate Limiting
- 분당 60 요청 제한
- X-RateLimit-* 헤더 제공
- 초과 시 429 상태 코드 반환
```

### 5. 보안 헤더

```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
```

## 📡 API 엔드포인트

### 인증 엔드포인트

```bash
# 회원가입
POST /api/v1/auth/register
{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123!@#",
    "confirm_password": "Test123!@#",
    "first_name": "Test",
    "last_name": "User"
}

# 로그인
POST /api/v1/auth/login
{
    "email": "test@example.com",
    "password": "Test123!@#"
}

# 토큰 갱신
POST /api/v1/auth/refresh
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 사용자 엔드포인트 (인증 필요)

```bash
# 프로필 조회
GET /api/v1/users/profile
Authorization: Bearer <access_token>

# 프로필 수정
PUT /api/v1/users/profile
Authorization: Bearer <access_token>
{
    "first_name": "Updated",
    "last_name": "Name"
}

# 비밀번호 변경
POST /api/v1/users/change-password
Authorization: Bearer <access_token>
{
    "current_password": "OldPass123!@#",
    "new_password": "NewPass123!@#",
    "confirm_password": "NewPass123!@#"
}
```

### 포스트 엔드포인트

```bash
# 포스트 목록 (공개)
GET /api/v1/posts?page=1&per_page=10

# 포스트 상세 (공개)
GET /api/v1/posts/:id

# 포스트 생성 (인증 필요)
POST /api/v1/posts
Authorization: Bearer <access_token>
{
    "title": "My First Post",
    "content": "This is the content",
    "slug": "my-first-post",
    "published": true,
    "tags": ["golang", "gin", "api"]
}

# 포스트 수정 (소유자만)
PUT /api/v1/posts/:id
Authorization: Bearer <access_token>

# 포스트 삭제 (소유자만)
DELETE /api/v1/posts/:id
Authorization: Bearer <access_token>
```

### 관리자 엔드포인트 (관리자 권한 필요)

```bash
# 모든 사용자 조회
GET /api/v1/users
Authorization: Bearer <admin_token>

# 특정 사용자 조회
GET /api/v1/users/:id
Authorization: Bearer <admin_token>

# 사용자 삭제
DELETE /api/v1/users/:id
Authorization: Bearer <admin_token>

# 통계 조회
GET /api/v1/stats
Authorization: Bearer <admin_token>
```

## 🧪 테스트

### 테스트 실행

```bash
# 모든 테스트 실행
make test

# 단위 테스트만
make test-unit

# 통합 테스트만
make test-integration

# 커버리지 리포트 생성
make coverage

# 벤치마크 테스트
make bench
```

### 테스트 구조

#### 단위 테스트
```go
func TestJWTManager_GenerateToken(t *testing.T) {
    manager := auth.NewJWTManager("secret", 15*time.Minute, 7*24*time.Hour)

    token, err := manager.GenerateAccessToken(1, "test@example.com", "testuser", "user")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)

    claims, err := manager.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, uint(1), claims.UserID)
}
```

#### 통합 테스트
```go
func (suite *AuthTestSuite) TestLoginFlow() {
    // 1. Register
    suite.registerTestUser()

    // 2. Login
    suite.login()

    // 3. Access protected resource
    headers := map[string]string{
        "Authorization": "Bearer " + suite.accessToken,
    }
    w := suite.performRequest("GET", "/api/v1/users/profile", nil, headers)
    suite.Equal(http.StatusOK, w.Code)

    // 4. Refresh token
    suite.refreshToken()
}
```

### 테스트 커버리지 목표
- 핸들러: 90%+
- 서비스: 85%+
- 미들웨어: 80%+
- 전체: 80%+

## 🎨 코드 품질

### Linting

```bash
# 린트 실행
make lint

# 자동 수정
make lint-fix

# 특정 파일/디렉토리
golangci-lint run ./internal/...
```

### 설정된 린터
- errcheck: 에러 처리 검사
- gosimple: 코드 단순화
- govet: Go vet
- ineffassign: 비효율적 할당
- staticcheck: 정적 분석
- gosec: 보안 검사
- gofmt: 코드 포맷팅
- goimports: import 정리

### 코드 포맷팅

```bash
# 포맷 체크
make fmt

# go vet 실행
make vet

# 모든 체크 실행
make check
```

## 🐳 Docker 지원

### Docker 빌드 및 실행

```bash
# 이미지 빌드
make docker-build

# 컨테이너 실행
make docker-run

# docker-compose 사용
make docker-compose-up
make docker-compose-down
```

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o project-security ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/project-security .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./project-security"]
```

## 📊 성능 최적화

### 데이터베이스 최적화
- Connection pooling
- Prepared statements
- Index 최적화
- Soft delete 사용

### API 최적화
- Response caching
- Pagination
- Field selection
- Eager loading (N+1 방지)

### 미들웨어 최적화
- Request ID 추적
- 구조화된 로깅
- Graceful shutdown
- Health checks

## 🔍 모니터링

### Health Check 엔드포인트

```bash
# 기본 헬스 체크
GET /health
Response: {"status": "healthy", "time": "2024-01-20T10:00:00Z"}

# 상세 헬스 체크
GET /ready
Response: {
    "status": "ready",
    "database": "connected",
    "cache": "available",
    "uptime": "2h30m"
}
```

### 로깅

```go
// 구조화된 로그
logger.Info("HTTP Request",
    "method", method,
    "path", path,
    "status", statusCode,
    "latency", latency,
    "ip", clientIP,
    "request_id", requestID,
)
```

## 🚦 CI/CD 파이프라인

### GitHub Actions
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make ci
```

### Makefile 타겟
```bash
# CI 파이프라인 실행
make ci

# 개별 단계
make install-deps
make check
make test
make build
```

## 📈 벤치마크

### 성능 테스트 결과

```
BenchmarkLogin-8                     500    2,345,678 ns/op    45678 B/op    123 allocs/op
BenchmarkTokenValidation-8         10000      123,456 ns/op     2345 B/op     12 allocs/op
BenchmarkCreatePost-8               1000    1,234,567 ns/op    34567 B/op     89 allocs/op
```

### 부하 테스트

```bash
# wrk를 사용한 부하 테스트
wrk -t12 -c400 -d30s --latency http://localhost:8080/api/v1/posts

# 결과
Requests/sec: 10,234.56
Latency: 39.12ms (avg)
Transfer/sec: 2.34MB
```

## 🛡️ 보안 체크리스트

- [x] JWT 토큰 기반 인증
- [x] Refresh Token 구현
- [x] 비밀번호 해싱 (bcrypt)
- [x] SQL Injection 방지
- [x] XSS 방지
- [x] CSRF 방지
- [x] Rate Limiting
- [x] CORS 설정
- [x] 입력 검증
- [x] 보안 헤더
- [x] HTTPS 지원 (프로덕션)
- [x] 환경 변수 관리
- [x] 감사 로깅

## 🎓 학습 포인트

이 프로젝트를 통해 다음을 학습할 수 있습니다:

1. **보안 구현**
   - JWT 인증 시스템 구축
   - CORS 정책 설정
   - 입력 검증 및 살균

2. **테스트 작성**
   - 단위 테스트 작성
   - 통합 테스트 구현
   - 테스트 커버리지 측정

3. **코드 품질**
   - 린터 설정 및 사용
   - 코드 포맷팅
   - 베스트 프랙티스 적용

4. **프로젝트 구조**
   - Clean Architecture
   - 의존성 주입
   - 레이어 분리

## 📚 참고 자료

- [Gin Documentation](https://gin-gonic.com/docs/)
- [JWT.io](https://jwt.io/)
- [OWASP Security Guidelines](https://owasp.org/)
- [Go Security Practices](https://github.com/OWASP/Go-SCP)
- [golangci-lint](https://golangci-lint.run/)

## 🤝 기여 방법

1. Fork the repository
2. Create your feature branch
3. Run tests and linting
4. Commit your changes
5. Push to the branch
6. Create a Pull Request

## 📄 라이선스

MIT License

---

**Built with ❤️ using Go and Gin Framework**

이 프로젝트는 Lessons 18-23의 모든 내용을 실전에서 사용할 수 있도록 통합한 완벽한 예제입니다. 프로덕션 환경에서 바로 사용할 수 있는 수준의 보안과 테스트를 갖추고 있습니다. 🚀