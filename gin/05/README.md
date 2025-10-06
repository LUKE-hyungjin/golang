# 미들웨어로 공통 기능 관리하기 🎯

안녕하세요! 이번 챕터에서는 **미들웨어**라는 강력한 기능을 배워보겠습니다. 미들웨어는 요청이 들어올 때마다 자동으로 실행되는 함수로, 인증, 로깅, 보안 등 공통 기능을 한 곳에서 관리할 수 있게 해줍니다.

## 미들웨어가 뭔가요?

미들웨어는 **요청과 응답 사이에서 실행되는 함수**입니다. 마치 공항 보안 검색대처럼, 모든 요청이 최종 목적지(핸들러)에 도달하기 전에 거쳐야 하는 체크포인트라고 생각하면 됩니다!

### 실생활 비유
- **공항 보안검색**: 탑승객(요청)이 게이트(핸들러)에 도착하기 전 보안검색(미들웨어)을 거쳐야 함
- **식당 입구**: 손님(요청)이 테이블(핸들러)에 앉기 전 웨이팅(미들웨어)을 거침
- **아파트 경비실**: 방문자(요청)가 집(핸들러)에 가기 전 경비실(미들웨어)에서 확인

## 이번 챕터에서 배울 내용
- 미들웨어가 무엇이고 왜 필요한지 이해하기
- 전역, 그룹, 개별 라우트에 미들웨어 적용하기
- `c.Next()`와 `c.Abort()`로 흐름 제어하기
- 나만의 미들웨어 만들어보기
- 실전 패턴: 인증, 로깅, 속도 제한 등

## 📂 파일 구조
```
05/
└── main.go     # 미들웨어 예제
```

## 핵심 개념 이해하기

### 1. 전역 미들웨어 - 모든 요청에 적용
```go
r.Use(Middleware())  // 모든 라우트에 적용
```

**언제 사용할까요?**
- 모든 요청에 Request ID를 부여하고 싶을 때
- 모든 요청을 로그에 기록하고 싶을 때
- 모든 요청의 응답 시간을 측정하고 싶을 때

### 2. 그룹 미들웨어 - 특정 그룹에만 적용
```go
v1 := r.Group("/api/v1")
v1.Use(Middleware())  // 그룹 내 모든 라우트에 적용
```

**언제 사용할까요?**
- `/api/v1` 아래 모든 경로에만 인증을 적용하고 싶을 때
- 관리자 경로(`/admin/*`)에만 권한 체크를 하고 싶을 때

### 3. 개별 라우트 미들웨어 - 딱 한 경로에만 적용
```go
r.GET("/path", Middleware(), Handler)  // 특정 라우트에만 적용
```

**언제 사용할까요?**
- 파일 업로드 경로에만 크기 제한을 적용하고 싶을 때
- 특정 API에만 속도 제한을 걸고 싶을 때

## 미들웨어는 어떻게 실행될까요?

미들웨어는 **양파 껍질처럼** 겹겹이 실행됩니다. 요청이 들어오면 바깥쪽 미들웨어부터 차례대로 실행되고, 응답은 반대로 안쪽에서 바깥쪽으로 나갑니다.

```
요청 들어옴 → 미들웨어1 시작 → 미들웨어2 시작 → 미들웨어3 시작 → 핸들러
                ↓                ↓                ↓                ↓
            c.Next()         c.Next()         c.Next()         응답 생성
                ↓                ↓                ↓                ↑
        미들웨어1 끝 ← 미들웨어2 끝 ← 미들웨어3 끝 ←────────────┘
```

**실생활 비유**: 러시아 인형(마트료시카)처럼, 작은 인형(핸들러)을 큰 인형들(미들웨어)이 감싸고 있는 구조입니다!

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

## 💡 꼭 알아야 할 핵심 개념!

### 1. 나만의 미들웨어 만들기

미들웨어는 이렇게 간단하게 만들 수 있어요!

```go
func MyMiddleware() gin.HandlerFunc {
    // 여기는 서버 시작할 때 딱 한 번만 실행돼요 (초기화)

    return func(c *gin.Context) {
        // Before: 요청이 핸들러에 가기 전에 실행되는 부분
        fmt.Println("요청이 들어왔어요!")

        c.Next()  // 다음 미들웨어나 핸들러로 넘어가요

        // After: 응답이 끝나고 돌아와서 실행되는 부분
        fmt.Println("응답을 보냈어요!")
    }
}
```

**실생활 비유**: 식당에서 음식을 주문하면, 웨이터(미들웨어)가 주문을 받고(Before) → 주방에 전달하고(c.Next()) → 음식이 나오면 서빙합니다(After)

### 2. c.Next() vs c.Abort() - 계속 vs 중단

**c.Next() - "다음 단계로 가세요!"**
```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("요청 시작")
        c.Next()  // 다음으로 진행
        fmt.Println("요청 끝")  // 이 줄도 실행됨!
    }
}
```

**c.Abort() - "여기서 멈춰!"**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "로그인이 필요해요"})
            return  // return도 꼭 써야 해요!
        }

        c.Next()  // 토큰이 있으면 계속 진행
    }
}
```

**실생활 비유**:
- `c.Next()`: 검문소에서 "통과!"하고 다음 단계로 진행
- `c.Abort()`: 검문소에서 "정지! 여기서 끝!"하고 뒤로 돌려보냄

### 3. 미들웨어끼리 정보 주고받기

미들웨어에서 저장한 정보를 핸들러에서 쓸 수 있어요!

```go
// 미들웨어에서 정보 저장
c.Set("user_id", 123)
c.Set("user_name", "홍길동")

// 핸들러에서 정보 가져오기
userID, exists := c.Get("user_id")
if exists {
    fmt.Println("사용자 ID:", userID)
}
```

**실생활 비유**: 은행에서 번호표를 받고(Set), 창구에서 번호표를 제시하는(Get) 것과 비슷해요!

### 4. 미들웨어 실행 순서

미들웨어는 **등록한 순서대로** 실행됩니다!

```go
r.Use(로깅())      // 1번째: 모든 요청 로깅
r.Use(인증())      // 2번째: 인증 체크

admin := r.Group("/admin")
admin.Use(권한체크())  // 3번째: 관리자 권한 체크

admin.GET("/", 특별처리(), 핸들러)  // 4번째: 특별 처리 → 핸들러
```

**순서가 중요한 이유**: 로깅을 먼저 해야 인증 실패도 기록되고, 인증을 먼저 확인해야 권한도 체크할 수 있어요!

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