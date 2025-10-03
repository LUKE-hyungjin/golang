# 10. 에러 핸들링 미들웨어

## 📌 개요
Gin 애플리케이션에서 발생하는 모든 에러를 중앙에서 일관성 있게 처리하는 미들웨어를 구현합니다. 패닉 복구, 에러 타입별 처리, 통일된 에러 응답 형식 등 프로덕션 환경에서 필요한 강력한 에러 처리 시스템을 구축합니다.

## 🎯 학습 목표
- 전역 에러 핸들링 미들웨어 구현
- 패닉 복구 미들웨어 작성
- 에러 타입별 처리 로직 구현
- 검증 에러 자동 처리
- 통일된 에러 응답 구조 적용
- 에러 체인과 래핑 처리

## 📂 파일 구조
```
10/
└── main.go     # 에러 핸들링 미들웨어 예제
```

## 💻 주요 미들웨어

### 1. RecoveryMiddleware
패닉 발생 시 애플리케이션 크래시를 방지하고 적절한 에러 응답 반환

### 2. ErrorHandlingMiddleware
모든 에러를 캡처하고 타입별로 적절히 처리

### 3. ValidationErrorMiddleware
바인딩 및 검증 에러를 자동으로 처리

### 4. CustomErrorMiddleware
HTTP 상태 코드 기반 자동 에러 응답 생성

## 🚀 실행 방법

```bash
cd gin
go run ./10

# 서버 실행 확인
curl http://localhost:8080/api/success
```

## 📋 에러 처리 테스트

### 1️⃣ 정상 응답

```bash
curl http://localhost:8080/api/success

# 응답:
{
  "success": true,
  "data": "Operation successful"
}
```

### 2️⃣ HTTP 에러

```bash
curl http://localhost:8080/api/http-error

# 응답:
{
  "success": false,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "The requested user does not exist"
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 3️⃣ 검증 에러

```bash
curl -X POST http://localhost:8080/api/validation-error \
  -H "Content-Type: application/json" \
  -d '{"id":0,"name":"Jo","email":"invalid","age":15}'

# 응답:
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format",
        "value": "invalid"
      },
      {
        "field": "age",
        "message": "Must be 18 or older",
        "value": "15"
      }
    ]
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 4️⃣ 비즈니스 에러

```bash
curl -X POST "http://localhost:8080/api/business-error?amount=0"

# 응답:
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "Transaction amount must be greater than zero"
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 5️⃣ 패닉 복구

```bash
curl http://localhost:8080/api/panic

# 응답:
{
  "success": false,
  "error": {
    "code": "PANIC_ERROR",
    "message": "This is a deliberate panic!",
    "stack": "goroutine 1 [running]:..." // 개발 모드에서만 표시
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}

# 서버 로그:
# PANIC RECOVERED: This is a deliberate panic!
# goroutine 1 [running]:
# main.RecoveryMiddleware.func1.1()
# ...
```

### 6️⃣ 일반 에러

```bash
curl http://localhost:8080/api/generic-error

# 응답:
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred"
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 7️⃣ 체인된 에러

```bash
curl http://localhost:8080/api/chained-error

# 응답:
{
  "success": false,
  "error": {
    "code": "DATABASE_ERROR",
    "message": "Unable to connect to database",
    "details": "database error: connection refused"
  },
  "request_id": "req-1234567890",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

### 8️⃣ 권한 에러

```bash
# 토큰 없음
curl http://localhost:8080/api/admin

# 응답:
{
  "success": false,
  "error": {
    "code": "AUTH_REQUIRED",
    "message": "Authentication token is required"
  }
}

# 잘못된 토큰
curl http://localhost:8080/api/admin \
  -H "Authorization: Bearer wrong-token"

# 응답:
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_PRIVILEGES",
    "message": "You don't have permission to access this resource"
  }
}

# 올바른 토큰
curl http://localhost:8080/api/admin \
  -H "Authorization: Bearer admin-token"

# 응답:
{
  "success": true,
  "message": "Welcome admin!"
}
```

### 9️⃣ 파일 업로드 에러

```bash
# 파일 없음
curl -X POST http://localhost:8080/api/upload

# 응답:
{
  "success": false,
  "error": {
    "code": "FILE_REQUIRED",
    "message": "No file was uploaded"
  }
}

# 파일 크기 초과 (5MB 이상)
curl -X POST http://localhost:8080/api/upload \
  -F "file=@large-file.jpg"

# 응답:
{
  "success": false,
  "error": {
    "code": "FILE_TOO_LARGE",
    "message": "File size exceeds maximum allowed",
    "details": {
      "max_size": "5MB",
      "uploaded_size": "10.5MB"
    }
  }
}
```

## 📝 핵심 포인트

### 1. 에러 인터페이스 설계

```go
type AppError interface {
    error
    Status() int        // HTTP 상태 코드
    Code() string       // 에러 코드
    Details() interface{} // 상세 정보
}
```

### 2. 미들웨어 실행 순서

```go
// 올바른 순서
r.Use(RecoveryMiddleware())        // 1. 패닉 복구 (가장 먼저)
r.Use(ErrorHandlingMiddleware())   // 2. 에러 처리
r.Use(ValidationErrorMiddleware()) // 3. 검증 에러

// 핸들러에서 에러 발생 시
c.Error(err)    // 에러 추가
c.Abort()       // 체인 중단
```

### 3. 에러 타입별 처리

```go
func handleError(c *gin.Context, err error) {
    // 타입 체크 순서가 중요
    switch e := err.(type) {
    case AppError:
        // 커스텀 애플리케이션 에러
    case ValidationError:
        // 검증 에러
    case BusinessError:
        // 비즈니스 로직 에러
    default:
        // 기본 에러
    }
}
```

### 4. 패닉 복구 패턴

```go
defer func() {
    if err := recover(); err != nil {
        // 스택 트레이스 로깅
        stack := debug.Stack()
        log.Printf("PANIC: %v\n%s", err, stack)

        // 클라이언트 응답
        respondWithError(c, 500, "PANIC", "Internal error")
        c.Abort()
    }
}()
```

## 🔍 트러블슈팅

### 에러가 처리되지 않는 경우

```go
// ❌ 잘못된 예: Abort() 호출 누락
c.Error(err)
// 다음 핸들러가 계속 실행됨

// ✅ 올바른 예: Abort() 호출
c.Error(err)
c.Abort()  // 체인 중단
```

### 중복 응답 방지

```go
// Response Writer 상태 체크
if !c.Writer.Written() {
    c.JSON(status, response)
}
```

### 에러 스택 트레이스

```go
// 개발/프로덕션 환경 구분
if gin.Mode() == gin.DebugMode {
    errorInfo.Stack = string(debug.Stack())
}
```

## 🏗️ 실전 활용 팁

### 1. 에러 로깅 통합

```go
func ErrorLoggingMiddleware(logger *log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        for _, err := range c.Errors {
            logger.Printf("[%s] %s %s - Error: %v",
                c.GetString("RequestID"),
                c.Request.Method,
                c.Request.URL.Path,
                err.Err,
            )
        }
    }
}
```

### 2. 에러 메트릭 수집

```go
func ErrorMetricsMiddleware(counter *prometheus.CounterVec) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            counter.WithLabelValues(
                c.Request.Method,
                c.Request.URL.Path,
                fmt.Sprintf("%d", c.Writer.Status()),
            ).Inc()
        }
    }
}
```

### 3. 에러 알림

```go
func CriticalErrorAlert(notifier Notifier) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        for _, err := range c.Errors {
            if isCritical(err.Err) {
                notifier.Send(fmt.Sprintf(
                    "Critical error: %v\nPath: %s\nRequestID: %s",
                    err.Err,
                    c.Request.URL.Path,
                    c.GetString("RequestID"),
                ))
            }
        }
    }
}
```

### 4. 에러 재시도 헤더

```go
func RetryableErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if c.Writer.Status() == http.StatusServiceUnavailable {
            c.Header("Retry-After", "30")
        }
    }
}
```

### 5. 상세 에러 모드

```go
type DetailedError struct {
    HTTPError
    File       string `json:"file,omitempty"`
    Line       int    `json:"line,omitempty"`
    Function   string `json:"function,omitempty"`
}

func NewDetailedError(status int, code, message string) DetailedError {
    // runtime.Caller를 사용하여 에러 발생 위치 추적
    pc, file, line, _ := runtime.Caller(1)
    fn := runtime.FuncForPC(pc)

    return DetailedError{
        HTTPError: HTTPError{status, code, message, nil},
        File:      file,
        Line:      line,
        Function:  fn.Name(),
    }
}
```

## 📚 다음 단계
- [11. 로깅 미들웨어](../11/README.md)

## 🔗 참고 자료
- [Gin 에러 처리 문서](https://gin-gonic.com/docs/examples/error-handling/)
- [Go 에러 처리 베스트 프랙티스](https://blog.golang.org/error-handling-and-go)
- [에러 래핑과 errors.Is/As](https://blog.golang.org/go1.13-errors)