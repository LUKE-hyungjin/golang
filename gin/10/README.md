# 에러를 한 곳에서 관리하기 🛡️

안녕하세요! 코드 여기저기서 에러가 발생하면 관리하기 힘들어요. 이번에는 **미들웨어**를 사용해서 모든 에러를 한 곳에서 깔끔하게 처리하는 방법을 배워봅시다!

## 에러 핸들링 미들웨어가 뭔가요?

에러 핸들링 미들웨어는 **에러 전담 팀**이라고 생각하면 돼요. 앱 어디서든 에러가 발생하면, 이 미들웨어가 자동으로 잡아서 적절하게 처리해줍니다.

### 실생활 비유
- **119 종합상황실**: 화재, 응급환자, 사고 등 모든 긴급상황을 한 곳에서 받아서 처리
- **고객센터**: 모든 고객 불만을 한 곳에서 접수받고 해결
- **공항 보안팀**: 모든 보안 문제를 전담으로 처리

### 왜 필요할까요?
```go
// ❌ 나쁜 예: 에러 처리가 여기저기 흩어져 있음
func Handler1(c *gin.Context) {
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})  // 형식이 제각각
    }
}

func Handler2(c *gin.Context) {
    if err != nil {
        c.JSON(400, gin.H{"message": "error"})  // 또 다른 형식
    }
}

// ✅ 좋은 예: 에러 처리를 한 곳에 모음
func Handler(c *gin.Context) {
    if err != nil {
        c.Error(err)  // 미들웨어가 알아서 처리!
        return
    }
}
```

## 이번 챕터에서 배울 내용
- 모든 에러를 한 곳에서 처리하는 미들웨어 만들기
- 패닉(서버 크래시)을 막고 복구하기
- 에러 종류별로 다르게 처리하기
- 검증 에러를 예쁘게 정리해서 보여주기
- 일관된 에러 응답 형식 만들기

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

## 💡 꼭 알아야 할 핵심 개념!

### 1. 에러를 체계적으로 분류하기

에러에도 **종류**가 있어요! 각 종류마다 다르게 처리해야 해요.

```go
type AppError interface {
    error
    Status() int        // HTTP 상태 코드 (400, 404, 500 등)
    Code() string       // 구체적인 에러 코드 ("USER_NOT_FOUND")
    Details() interface{} // 추가 정보
}
```

**에러 종류 예시**:
- **인증 에러**: 로그인 안 했거나, 권한 없음
- **검증 에러**: 이메일 형식 틀림, 비밀번호 짧음
- **비즈니스 에러**: 잔액 부족, 재고 없음
- **서버 에러**: 데이터베이스 연결 실패

### 2. 미들웨어 순서가 중요해요!

에러 처리 미들웨어는 **가장 먼저** 등록해야 모든 에러를 잡을 수 있어요!

```go
// ✅ 올바른 순서
r.Use(RecoveryMiddleware())        // 1순위: 패닉 복구 (가장 바깥)
r.Use(ErrorHandlingMiddleware())   // 2순위: 에러 처리
r.Use(ValidationErrorMiddleware()) // 3순위: 검증 에러
r.Use(LoggingMiddleware())         // 4순위: 로깅

// 핸들러에서 에러가 나면
c.Error(err)    // 에러를 미들웨어에 전달
c.Abort()       // 더 이상 진행하지 말고 멈춤
```

**실생활 비유**: 안전망을 설치할 때 가장 아래부터 설치하는 것처럼!

### 3. 에러 타입별로 다르게 응답하기

같은 에러라도 **종류에 따라** 다른 메시지를 보내야 해요!

```go
func handleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case AppError:
        // 우리가 만든 에러 → 상세하게 알려줌
        c.JSON(e.Status(), gin.H{
            "error": e.Code(),
            "message": e.Message(),
        })

    case ValidationError:
        // 검증 실패 → 어떤 필드가 틀렸는지 알려줌
        c.JSON(400, gin.H{
            "error": "VALIDATION_ERROR",
            "fields": e.Fields(),
        })

    default:
        // 모르는 에러 → 일반적인 메시지만
        c.JSON(500, gin.H{
            "error": "INTERNAL_ERROR",
            "message": "문제가 발생했습니다",
        })
    }
}
```

**실생활 비유**: 병원에서 감기, 골절, 화상을 각각 다른 방법으로 치료하는 것!

### 4. 패닉 복구 - 서버가 죽지 않게 하기

Go 프로그램에서 **패닉**이 발생하면 서버가 죽어버려요! 미들웨어로 이를 막을 수 있어요.

```go
defer func() {
    if err := recover(); err != nil {
        // 패닉 발생! 하지만 서버는 계속 돌아가요

        // 1. 로그에 자세히 기록 (개발자가 나중에 확인)
        stack := debug.Stack()
        log.Printf("🚨 PANIC: %v\n%s", err, stack)

        // 2. 사용자에게는 간단한 메시지
        c.JSON(500, gin.H{
            "error": "INTERNAL_ERROR",
            "message": "일시적인 문제가 발생했습니다",
        })

        c.Abort()  // 요청 처리 중단
    }
}()
```

**실생활 비유**:
- **패닉 복구 없음**: 공장 기계 하나가 고장나면 전체 공장이 멈춤
- **패닉 복구 있음**: 한 기계가 고장나도 다른 기계들은 계속 돌아감

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