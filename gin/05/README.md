# 05. 미들웨어 (전역/그룹/개별)와 next() 흐름

## 📌 개요
Gin의 미들웨어 시스템을 완벽하게 이해하고 활용하는 방법을 학습합니다. 미들웨어는 요청 처리 파이프라인에서 핵심 역할을 하며, 인증, 로깅, CORS, 에러 처리 등 횡단 관심사를 처리하는 데 사용됩니다.

## 🎯 학습 목표
- 미들웨어의 개념과 동작 원리 이해
- 전역, 그룹, 라우트별 미들웨어 적용
- c.Next()와 c.Abort()의 흐름 제어
- 커스텀 미들웨어 작성 방법
- 실전 미들웨어 패턴 (인증, 로깅, CORS, Rate Limiting)

## 📂 파일 구조
```
05/
└── main.go     # 미들웨어 예제
```

## 💻 미들웨어 타입과 적용 범위

### 1. 전역 미들웨어
```go
r.Use(Middleware())  // 모든 라우트에 적용
```

### 2. 그룹 미들웨어
```go
v1 := r.Group("/api/v1")
v1.Use(Middleware())  // 그룹 내 모든 라우트에 적용
```

### 3. 라우트별 미들웨어
```go
r.GET("/path", Middleware(), Handler)  // 특정 라우트에만 적용
```

## 🔄 미들웨어 실행 흐름

```
Request → MW1 Before → MW2 Before → MW3 Before → Handler
           ↓            ↓            ↓            ↓
        c.Next()     c.Next()     c.Next()    Response
           ↓            ↓            ↓            ↑
      MW1 After ← MW2 After ← MW3 After ←────────┘
```

## 🚀 실행 방법

```bash
cd gin
go run ./05
```

## 📋 API 테스트 예제

### 1️⃣ 기본 엔드포인트 (전역 미들웨어 적용)

```bash
# Request ID가 자동으로 추가됨
curl http://localhost:8080/

# 응답:
# {
#   "message": "Welcome to Gin Middleware Example",
#   "request_id": "req-1234567890"
# }

# 헬스 체크
curl http://localhost:8080/health
```

### 2️⃣ Rate Limiting 테스트

```bash
# Rate limit: 분당 10회
for i in {1..12}; do
  echo "Request $i:"
  curl http://localhost:8080/api/v1/public
  echo ""
done

# 11번째 요청부터:
# {"error":"Rate limit exceeded","retry_after":"60 seconds"}
```

### 3️⃣ 인증 미들웨어

**인증 없이 접근 (실패):**
```bash
curl http://localhost:8080/api/v1/protected/profile

# 응답:
# {"error":"Invalid authorization header"}
```

**일반 사용자 토큰으로 접근:**
```bash
curl http://localhost:8080/api/v1/protected/profile \
  -H "Authorization: Bearer valid-token-123"

# 응답:
# {
#   "message": "Protected resource",
#   "user": {
#     "id": "1",
#     "username": "testuser",
#     "role": "user"
#   }
# }
```

### 4️⃣ 역할 기반 접근 제어

**일반 사용자가 관리자 영역 접근 (실패):**
```bash
curl http://localhost:8080/api/v1/protected/admin/users \
  -H "Authorization: Bearer valid-token-123"

# 응답:
# {"error":"Requires admin role"}
```

**관리자 토큰으로 접근:**
```bash
curl http://localhost:8080/api/v1/protected/admin/users \
  -H "Authorization: Bearer admin-token-456"

# 응답:
# {
#   "message": "Admin only resource",
#   "users": [
#     {"id":"1","username":"admin","role":"admin"},
#     {"id":"2","username":"user1","role":"user"}
#   ]
# }

# 사용자 삭제
curl -X DELETE http://localhost:8080/api/v1/protected/admin/users/123 \
  -H "Authorization: Bearer admin-token-456"
```

### 5️⃣ 타임아웃 미들웨어

```bash
# 2초 타임아웃, 1초 처리 (성공)
curl http://localhost:8080/slow

# 응답:
# {"message":"Slow operation completed"}
```

### 6️⃣ 미들웨어 체인 흐름 확인

```bash
curl http://localhost:8080/middleware-chain

# 응답:
# {
#   "message": "Handler executed",
#   "flow": [
#     "first-before",
#     "second-before",
#     "third-before",
#     "first-after",
#     "second-after",
#     "third-after"
#   ]
# }

# 서버 로그:
# 1. First Middleware - Before
# 2. Second Middleware - Before
# 3. Third Middleware - Before
# 4. Main Handler
# 5. Main Handler - After response
# 6. First Middleware - After
# 7. Second Middleware - After
# 8. Third Middleware - After
```

### 7️⃣ 조건부 미들웨어

```bash
# 인증 필요 (실패)
curl http://localhost:8080/conditional

# 응답:
# {"error":"Authorization required (use ?skip_auth=true to bypass)"}

# 인증 스킵
curl "http://localhost:8080/conditional?skip_auth=true"

# 응답:
# {"message":"Conditional middleware passed"}
```

### 8️⃣ CORS 테스트

```bash
# OPTIONS 요청 (Preflight)
curl -X OPTIONS http://localhost:8080/api/v1/public \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -v

# 헤더 확인:
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE
```

### 9️⃣ 에러 처리 미들웨어

```bash
curl http://localhost:8080/error

# 응답:
# {"error":"Internal server error"}

# 서버 로그:
# Error occurred: something went wrong
```

### 🔟 데이터 변환 미들웨어

```bash
curl -X POST http://localhost:8080/transform \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john",
    "city": "seoul",
    "country": "korea"
  }'

# 응답:
# {
#   "original": {
#     "name": "john",
#     "city": "seoul",
#     "country": "korea"
#   },
#   "transformed": {
#     "name": "JOHN",
#     "city": "SEOUL",
#     "country": "KOREA"
#   }
# }
```

## 📝 핵심 포인트

### 1. 미들웨어 작성 패턴

```go
func MyMiddleware() gin.HandlerFunc {
    // 초기화 코드 (한 번만 실행)

    return func(c *gin.Context) {
        // Before: 요청 처리 전

        c.Next()  // 다음 미들웨어/핸들러 실행

        // After: 응답 후
    }
}
```

### 2. c.Next() vs c.Abort()

```go
// c.Next(): 다음 핸들러 실행 후 돌아옴
func Middleware1() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("Before")
        c.Next()
        fmt.Println("After")  // 실행됨
    }
}

// c.Abort(): 체인 중단
func Middleware2() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !authorized {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return  // return 필수!
        }
        c.Next()
    }
}
```

### 3. 미들웨어 간 데이터 전달

```go
// 설정
c.Set("key", value)

// 가져오기
value, exists := c.Get("key")
if exists {
    // 사용
}
```

### 4. 미들웨어 실행 순서

```go
r.Use(MW1())  // 1번째
r.Use(MW2())  // 2번째

group := r.Group("/")
group.Use(MW3())  // 3번째

group.GET("/", MW4(), Handler)  // 4번째 → Handler
```

## 🔍 트러블슈팅

### Abort 후에도 코드가 실행되는 경우

```go
// 잘못된 예
c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
fmt.Println("This will still execute!")  // 실행됨!

// 올바른 예
c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
return  // 명시적 return 필요
```

### 미들웨어에서 Body 읽기

```go
// Body는 한 번만 읽을 수 있음
// 재사용이 필요한 경우:
bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
```

### 고루틴에서 Context 사용

```go
// 미들웨어에서 고루틴 사용 시
cCp := c.Copy()  // 복사 필수!
go func() {
    // cCp 사용
}()
```

## 🏗️ 실전 미들웨어 패턴

### 1. JWT 인증 미들웨어

```go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        claims, err := validateJWT(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

### 2. 요청/응답 로깅

```go
func DetailedLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        // Request 로깅
        log.Printf("→ %s %s", c.Request.Method, path)

        c.Next()

        // Response 로깅
        latency := time.Since(start)
        status := c.Writer.Status()
        log.Printf("← %s %s %d %v", c.Request.Method, path, status, latency)
    }
}
```

### 3. 에러 복구 미들웨어

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                c.AbortWithStatusJSON(500, gin.H{
                    "error": "Internal server error",
                })
            }
        }()
        c.Next()
    }
}
```

## 📚 다음 단계
- [06. 라우트 그룹과 버저닝](../06/README.md)
- [07. 정적 파일 서빙](../07/README.md)
- [08. 템플릿 렌더링](../08/README.md)

## 🔗 참고 자료
- [Gin 미들웨어 문서](https://gin-gonic.com/docs/examples/using-middleware/)
- [Gin Contrib 미들웨어 모음](https://github.com/gin-contrib)
- [HTTP 미들웨어 패턴](https://www.alexedwards.net/blog/making-and-using-middleware)