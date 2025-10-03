# 09. HTTP 상태코드와 에러 응답 규약

## 📌 개요
웹 API에서 일관성 있고 명확한 에러 처리는 매우 중요합니다. 이 장에서는 HTTP 상태 코드를 올바르게 사용하고, 표준화된 에러 응답 구조를 구현하는 방법을 학습합니다. RESTful API의 에러 처리 모범 사례를 따르는 체계적인 에러 관리 시스템을 구축합니다.

## 🎯 학습 목표
- HTTP 상태 코드의 올바른 사용법 이해
- 표준화된 에러 응답 구조 설계
- 상태 코드별 헬퍼 함수 구현
- 비즈니스 로직 에러 처리
- 검증 에러의 구조화된 응답
- 에러 추적을 위한 Request ID 활용

## 📂 파일 구조
```
09/
└── main.go     # HTTP 상태코드와 에러 응답 예제
```

## 💻 주요 구성 요소

### 1. 표준 에러 응답 구조
```go
type StandardError struct {
    Code      int         `json:"code"`        // HTTP 상태 코드
    Message   string      `json:"message"`     // 사용자 메시지
    ErrorCode string      `json:"error_code"`  // 내부 에러 코드
    Details   interface{} `json:"details"`     // 상세 정보
    Timestamp time.Time   `json:"timestamp"`   // 발생 시간
    Path      string      `json:"path"`        // 요청 경로
    RequestID string      `json:"request_id"`  // 추적 ID
}
```

### 2. HTTP 상태 코드 분류
- **2xx Success**: 요청 성공
- **4xx Client Error**: 클라이언트 오류
- **5xx Server Error**: 서버 오류

## 🚀 실행 방법

```bash
cd gin
go run ./09

# 서버 실행 확인
curl http://localhost:8080/api/users
```

## 📋 HTTP 상태 코드별 테스트

### 1️⃣ 성공 응답 (2xx)

**200 OK - 성공적인 조회:**
```bash
curl http://localhost:8080/api/users

# 응답:
{
  "success": true,
  "data": [
    {"id": 1, "name": "John", "email": "john@example.com"},
    {"id": 2, "name": "Jane", "email": "jane@example.com"}
  ],
  "meta": {
    "total": 2,
    "page": 1
  }
}
```

**201 Created - 리소스 생성:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

# 응답:
{
  "success": true,
  "data": {
    "id": 123,
    "name": "Alice",
    "email": "alice@example.com",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

**204 No Content - 삭제 성공:**
```bash
curl -X DELETE http://localhost:8080/api/users/123 -I

# 응답:
HTTP/1.1 204 No Content
```

### 2️⃣ 클라이언트 에러 (4xx)

**400 Bad Request - 잘못된 요청:**
```bash
curl http://localhost:8080/api/bad-request

# 응답:
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Missing required parameters",
    "error_code": "BAD_REQUEST",
    "details": {
      "required": ["name", "email"],
      "provided": []
    },
    "timestamp": "2024-01-01T10:00:00Z",
    "path": "/api/bad-request",
    "request_id": "req-1234567890"
  }
}
```

**401 Unauthorized - 인증 필요:**
```bash
curl http://localhost:8080/api/protected

# 응답:
{
  "success": false,
  "error": {
    "code": 401,
    "message": "Authentication required",
    "error_code": "UNAUTHORIZED",
    "timestamp": "2024-01-01T10:00:00Z",
    "path": "/api/protected",
    "request_id": "req-1234567890"
  }
}

# 유효한 토큰으로 요청:
curl http://localhost:8080/api/protected \
  -H "Authorization: Bearer valid-token"
```

**403 Forbidden - 권한 없음:**
```bash
curl -X DELETE http://localhost:8080/api/admin/users

# 응답:
{
  "success": false,
  "error": {
    "code": 403,
    "message": "Admin access required",
    "error_code": "FORBIDDEN"
  }
}
```

**404 Not Found - 리소스 없음:**
```bash
curl http://localhost:8080/api/users/999

# 응답:
{
  "success": false,
  "error": {
    "code": 404,
    "message": "User not found",
    "error_code": "NOT_FOUND"
  }
}
```

**409 Conflict - 충돌:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"duplicate@example.com"}'

# 응답:
{
  "success": false,
  "error": {
    "code": 409,
    "message": "Email already exists",
    "error_code": "CONFLICT"
  }
}
```

**422 Unprocessable Entity - 검증 실패:**
```bash
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid","password":"123","age":15}'

# 응답:
{
  "success": false,
  "error": {
    "code": 422,
    "message": "Validation failed",
    "error_code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format",
        "value": "invalid"
      },
      {
        "field": "password",
        "message": "Password must be at least 6 characters"
      },
      {
        "field": "age",
        "message": "Must be 18 or older",
        "value": "15"
      }
    ]
  }
}
```

**429 Too Many Requests - 요청 제한:**
```bash
curl http://localhost:8080/api/rate-limited

# 응답:
{
  "success": false,
  "error": {
    "code": 429,
    "message": "Rate limit exceeded",
    "error_code": "TOO_MANY_REQUESTS",
    "details": {
      "limit": 100,
      "remaining": 0,
      "reset_after": "60 seconds"
    }
  }
}
```

### 3️⃣ 서버 에러 (5xx)

**500 Internal Server Error:**
```bash
curl http://localhost:8080/api/error

# 응답:
{
  "success": false,
  "error": {
    "code": 500,
    "message": "An unexpected error occurred",
    "error_code": "INTERNAL_SERVER_ERROR"
  }
}

# 데이터베이스 에러:
curl "http://localhost:8080/api/error?type=db"

# 응답:
{
  "success": false,
  "error": {
    "code": 500,
    "message": "Database connection failed",
    "error_code": "DATABASE_ERROR",
    "details": {
      "retry_after": "30 seconds"
    }
  }
}
```

**502 Bad Gateway:**
```bash
curl http://localhost:8080/api/external

# 응답:
{
  "success": false,
  "error": {
    "code": 502,
    "message": "External service is not responding",
    "error_code": "EXTERNAL_SERVICE_ERROR",
    "details": {
      "service": "payment-gateway",
      "timeout": "30s"
    }
  }
}
```

**503 Service Unavailable:**
```bash
curl http://localhost:8080/api/maintenance

# 응답:
{
  "success": false,
  "error": {
    "code": 503,
    "message": "Service is under maintenance",
    "error_code": "SERVICE_UNAVAILABLE",
    "details": {
      "retry_after": "2024-01-01T11:00:00Z"
    }
  }
}
```

### 4️⃣ 비즈니스 로직 에러

**비즈니스 규칙 위반:**
```bash
# 금액이 음수인 경우
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"account1","to":"account2","amount":-100}'

# 응답:
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Transfer amount must be positive",
    "error_code": "INVALID_AMOUNT",
    "details": {
      "amount": -100
    }
  }
}

# 한도 초과
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"account1","to":"account2","amount":20000}'

# 응답:
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Transfer amount exceeds daily limit",
    "error_code": "AMOUNT_LIMIT_EXCEEDED",
    "details": {
      "amount": 20000,
      "limit": 10000
    }
  }
}

# 잔액 부족
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"poor-account","to":"account2","amount":500}'

# 응답:
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Insufficient funds in source account",
    "error_code": "INSUFFICIENT_FUNDS",
    "details": {
      "available": 100,
      "requested": 500
    }
  }
}
```

### 5️⃣ 파일 업로드 에러

**파일 크기 초과:**
```bash
# 5MB 이상 파일 업로드 시뮬레이션
curl -X POST http://localhost:8080/api/upload \
  -F "file=@large-file.jpg"

# 응답:
{
  "success": false,
  "error": {
    "code": 413,
    "message": "File size exceeds maximum allowed size",
    "error_code": "FILE_TOO_LARGE",
    "details": {
      "max_size": "5MB",
      "uploaded_size": "10.5MB"
    }
  }
}
```

## 📝 핵심 포인트

### 1. HTTP 상태 코드 선택 가이드

| 상태 코드 | 사용 시점 | 예시 |
|----------|----------|------|
| 200 OK | 성공적인 GET, PUT | 사용자 조회, 수정 완료 |
| 201 Created | POST로 리소스 생성 | 새 사용자 생성 |
| 204 No Content | 성공했지만 응답 본문 없음 | DELETE 성공 |
| 400 Bad Request | 잘못된 요청 구문 | 잘못된 JSON 형식 |
| 401 Unauthorized | 인증 필요 | 토큰 없음/만료 |
| 403 Forbidden | 권한 없음 | 관리자 기능 접근 |
| 404 Not Found | 리소스 없음 | 존재하지 않는 사용자 |
| 409 Conflict | 상태 충돌 | 중복된 이메일 |
| 422 Unprocessable Entity | 검증 실패 | 유효하지 않은 필드 |
| 429 Too Many Requests | 요청 제한 초과 | API 호출 제한 |
| 500 Internal Server Error | 서버 오류 | 예기치 않은 오류 |
| 502 Bad Gateway | 외부 서비스 오류 | 결제 게이트웨이 오류 |
| 503 Service Unavailable | 서비스 일시 중단 | 유지보수 중 |

### 2. 에러 응답 구조 설계 원칙

```go
// 일관된 구조 사용
type ErrorResponse struct {
    Success bool           `json:"success"`  // 항상 false
    Error   *StandardError `json:"error"`    // 에러 상세 정보
}

// 성공 응답도 일관된 구조
type SuccessResponse struct {
    Success bool        `json:"success"`  // 항상 true
    Data    interface{} `json:"data"`     // 응답 데이터
    Meta    interface{} `json:"meta"`     // 메타 정보
}
```

### 3. 에러 코드 체계

```go
// 도메인별 에러 코드
const (
    // 인증 관련
    ErrAuthTokenExpired = "AUTH_TOKEN_EXPIRED"
    ErrAuthInvalidToken = "AUTH_INVALID_TOKEN"

    // 사용자 관련
    ErrUserNotFound     = "USER_NOT_FOUND"
    ErrUserDuplicate    = "USER_DUPLICATE"

    // 비즈니스 로직
    ErrInsufficientFunds = "INSUFFICIENT_FUNDS"
    ErrLimitExceeded     = "LIMIT_EXCEEDED"
)
```

### 4. Request ID를 통한 추적

```go
// 모든 요청에 고유 ID 할당
c.Set("RequestID", generateRequestID())

// 에러 응답에 포함
error.RequestID = c.GetString("RequestID")

// 로그에도 기록
log.Printf("[%s] Error: %s", requestID, error.Message)
```

## 🔍 트러블슈팅

### 적절한 상태 코드 선택

```go
// ❌ 잘못된 예: 모든 에러에 500 사용
c.JSON(500, gin.H{"error": "User not found"})

// ✅ 올바른 예: 적절한 상태 코드 사용
c.JSON(404, gin.H{"error": "User not found"})
```

### 민감한 정보 노출 방지

```go
// ❌ 위험: 내부 정보 노출
c.JSON(500, gin.H{
    "error": err.Error(),  // 스택 트레이스 노출
    "query": sqlQuery,     // SQL 쿼리 노출
})

// ✅ 안전: 일반적인 메시지
c.JSON(500, gin.H{
    "error": "Internal server error",
    "request_id": requestID,  // 추적용 ID만 제공
})
```

## 🏗️ 실전 활용 팁

### 1. 에러 핸들러 중앙화

```go
func HandleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case BusinessError:
        NewErrorResponse(c, e.Status, e.Code, e.Message, nil)
    case ValidationError:
        ValidationFailed(c, e.Errors)
    default:
        InternalServerError(c, "An error occurred")
    }
}
```

### 2. 에러 로깅

```go
func LogError(c *gin.Context, err error) {
    log.Printf(
        "Error: [%s] %s %s - %v",
        c.GetString("RequestID"),
        c.Request.Method,
        c.Request.URL.Path,
        err,
    )
}
```

### 3. 환경별 에러 상세 정보

```go
func GetErrorDetails(err error) interface{} {
    if gin.Mode() == gin.DebugMode {
        return err.Error()  // 개발 환경: 상세 정보
    }
    return nil  // 프로덕션: 상세 정보 숨김
}
```

### 4. 재시도 가능 여부 표시

```go
type ErrorResponse struct {
    // ... 기존 필드들
    Retryable    bool   `json:"retryable"`
    RetryAfter   int    `json:"retry_after,omitempty"`
}
```

## 📚 다음 단계
- [10. 에러 핸들링 미들웨어](../10/README.md)
- [11. 로깅 미들웨어](../11/README.md)

## 🔗 참고 자료
- [HTTP 상태 코드 MDN](https://developer.mozilla.org/ko/docs/Web/HTTP/Status)
- [REST API 에러 처리 가이드](https://www.baeldung.com/rest-api-error-handling-best-practices)
- [RFC 7807 - Problem Details for HTTP APIs](https://tools.ietf.org/html/rfc7807)