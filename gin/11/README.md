# 11. 로깅 미들웨어 (요청/응답 로깅)

## 📌 개요
프로덕션 환경에서 애플리케이션의 동작을 추적하고 문제를 진단하기 위한 포괄적인 로깅 시스템을 구현합니다. 구조화된 로그, 접근 로그, 감사 로그, 성능 모니터링 등 다양한 로깅 전략을 학습합니다.

## 🎯 학습 목표
- 구조화된 로깅 (JSON 형식) 구현
- 요청/응답 정보 캡처 및 로깅
- 로그 레벨 관리 (DEBUG, INFO, WARN, ERROR)
- 파일 로테이션 구현
- 접근 로그 및 감사 로그 작성
- 느린 요청 감지 및 로깅
- 민감한 정보 마스킹

## 📂 파일 구조
```
11/
├── main.go     # 로깅 미들웨어 예제
└── logs/       # 로그 파일 디렉토리 (자동 생성)
    ├── app.log     # 애플리케이션 로그
    └── access.log  # 접근 로그
```

## 💻 로깅 미들웨어 종류

### 1. StructuredLoggingMiddleware
JSON 형식의 구조화된 로그 생성

### 2. AccessLoggingMiddleware
Apache Combined Log Format 형식의 접근 로그

### 3. SlowRequestLoggingMiddleware
임계값을 초과하는 느린 요청 감지

### 4. ErrorLoggingMiddleware
에러 전용 상세 로깅

### 5. AuditLoggingMiddleware
중요한 작업에 대한 감사 로그

## 🚀 실행 방법

```bash
cd gin
go run ./11

# 로그 파일 확인
tail -f 11/logs/app.log
tail -f 11/logs/access.log
```

## 📋 로깅 테스트

### 1️⃣ 기본 요청 로깅

```bash
curl http://localhost:8080/api/health

# 콘솔 출력 (JSON 형식):
{
  "level": "INFO",
  "timestamp": "2024-01-01T10:00:00Z",
  "request_id": "req-1234567890",
  "method": "GET",
  "path": "/api/health",
  "status_code": 200,
  "latency": "1.234ms",
  "client_ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}

# access.log:
127.0.0.1 - - [01/Jan/2024:10:00:00 +0900] "GET /api/health HTTP/1.1" 200 45 "" "curl/7.68.0" 1.234ms
```

### 2️⃣ 느린 요청 감지

```bash
curl http://localhost:8080/api/slow

# app.log에 경고 로그:
{
  "level": "WARN",
  "timestamp": "2024-01-01T10:00:01Z",
  "request_id": "req-1234567891",
  "method": "GET",
  "path": "/api/slow",
  "status_code": 200,
  "latency": "200ms",
  "client_ip": "127.0.0.1",
  "extra": {
    "threshold": "100ms",
    "exceeded": "100ms"
  }
}
```

### 3️⃣ 에러 로깅

```bash
curl http://localhost:8080/api/error

# 에러 로그:
{
  "level": "ERROR",
  "timestamp": "2024-01-01T10:00:02Z",
  "request_id": "req-1234567892",
  "method": "GET",
  "path": "/api/error",
  "status_code": 500,
  "client_ip": "127.0.0.1",
  "error": "this is a test error",
  "extra": {
    "type": 1,
    "meta": null
  }
}
```

### 4️⃣ 감사 로그 (사용자 생성)

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token" \
  -d '{"name":"John","email":"john@example.com"}'

# 감사 로그:
{
  "level": "INFO",
  "timestamp": "2024-01-01T10:00:03Z",
  "request_id": "req-1234567893",
  "method": "POST",
  "path": "/api/users",
  "status_code": 201,
  "latency": "5ms",
  "client_ip": "127.0.0.1",
  "extra": {
    "user_id": "user123",
    "user_role": "admin",
    "action": "CREATE",
    "resource": "users",
    "request_body": "{\"name\":\"John\",\"email\":\"john@example.com\"}"
  }
}
```

### 5️⃣ 민감한 정보 마스킹

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret123"}'

# 로그 (password 마스킹됨):
{
  "level": "INFO",
  "timestamp": "2024-01-01T10:00:04Z",
  "request_id": "req-1234567894",
  "method": "POST",
  "path": "/api/login",
  "status_code": 200,
  "extra": {
    "action": "CREATE",
    "resource": "login",
    "request_body": "{\"username\":\"admin\",\"password\":\"***MASKED***\"}"
  }
}
```

### 6️⃣ 디버그 모드 상세 로깅

```bash
# 디버그 모드에서는 헤더와 body 정보 포함
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "X-Custom-Header: custom-value" \
  -d '{"name":"Jane"}'

# 상세 로그:
{
  "level": "INFO",
  "timestamp": "2024-01-01T10:00:05Z",
  "request_id": "req-1234567895",
  "method": "POST",
  "path": "/api/users",
  "status_code": 201,
  "request": {
    "headers": {
      "Content-Type": "application/json",
      "X-Custom-Header": "custom-value"
    },
    "body": "{\"name\":\"Jane\"}",
    "query": ""
  },
  "response": {
    "headers": {
      "Content-Type": "application/json; charset=utf-8"
    },
    "body": "{\"id\":\"user-123\",\"message\":\"User created successfully\"}",
    "size": 58
  }
}
```

### 7️⃣ 다양한 상태 코드 로깅

```bash
# 400 Bad Request (WARN 레벨)
curl http://localhost:8080/api/status/400

# 401 Unauthorized (WARN 레벨)
curl http://localhost:8080/api/status/401

# 500 Internal Error (ERROR 레벨)
curl http://localhost:8080/api/status/500

# 로그 레벨이 상태 코드에 따라 자동 결정됨
```

### 8️⃣ 사용자 삭제 (감사 로그)

```bash
curl -X DELETE http://localhost:8080/api/users/123 \
  -H "Authorization: Bearer admin-token"

# 감사 로그:
{
  "level": "INFO",
  "timestamp": "2024-01-01T10:00:06Z",
  "request_id": "req-1234567896",
  "method": "DELETE",
  "path": "/api/users/123",
  "status_code": 200,
  "extra": {
    "user_id": "user123",
    "user_role": "admin",
    "action": "DELETE",
    "resource": "users"
  }
}
```

## 📝 핵심 포인트

### 1. 로그 레벨 관리

```go
type LogLevel int

const (
    DEBUG LogLevel = iota  // 개발 디버깅
    INFO                   // 일반 정보
    WARN                   // 경고
    ERROR                  // 에러
    FATAL                  // 치명적 에러
)

// 환경별 로그 레벨 설정
// 개발: DEBUG
// 스테이징: INFO
// 프로덕션: WARN
```

### 2. 구조화된 로그 형식

```go
type LogEntry struct {
    Level      string    `json:"level"`
    Timestamp  string    `json:"timestamp"`
    RequestID  string    `json:"request_id"`
    Method     string    `json:"method"`
    Path       string    `json:"path"`
    StatusCode int       `json:"status_code"`
    Latency    string    `json:"latency"`
    ClientIP   string    `json:"client_ip"`
    Error      string    `json:"error,omitempty"`
}
```

### 3. 민감한 정보 처리

```go
// 민감한 헤더 제외
sensitiveHeaders := []string{
    "Authorization",
    "Cookie",
    "X-Api-Key",
}

// 민감한 필드 마스킹
sensitiveFields := []string{
    "password",
    "token",
    "secret",
    "credit_card",
}
```

### 4. 로그 로테이션

```go
// 파일 크기 기반 로테이션
if fileSize >= maxSize {
    // 기존 파일 백업
    os.Rename("app.log", "app.log.20240101")
    // 새 파일 생성
    createNewLogFile()
}
```

## 🔍 트러블슈팅

### 메모리 사용량 증가

```go
// ❌ 문제: 모든 body를 메모리에 저장
body, _ := ioutil.ReadAll(c.Request.Body)

// ✅ 해결: 크기 제한
if c.Request.ContentLength < 1024 { // 1KB 이하만
    body, _ := ioutil.ReadAll(c.Request.Body)
}
```

### 성능 영향 최소화

```go
// 비동기 로깅
go func() {
    logger.Info(entry)
}()

// 버퍼링된 채널 사용
logChannel := make(chan LogEntry, 1000)
```

### 로그 파일 권한

```bash
# 로그 디렉토리 권한 설정
chmod 755 logs/
chmod 644 logs/*.log
```

## 🏗️ 실전 활용 팁

### 1. 분산 추적

```go
// OpenTelemetry 통합
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, span := tracer.Start(c.Request.Context(), "http.request")
        defer span.End()

        c.Request = c.Request.WithContext(ctx)
        c.Set("TraceID", span.SpanContext().TraceID().String())
        c.Next()
    }
}
```

### 2. 로그 집계

```go
// ELK Stack으로 전송
func SendToElasticsearch(entry LogEntry) {
    client.Index().
        Index("app-logs").
        Type("_doc").
        BodyJson(entry).
        Do(context.Background())
}
```

### 3. 메트릭 수집

```go
// Prometheus 메트릭
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
        },
        []string{"method", "path", "status"},
    )
)
```

### 4. 알림 통합

```go
// 에러 임계값 초과 시 알림
func AlertOnErrors(threshold int) {
    if errorCount > threshold {
        slack.PostMessage(channel, fmt.Sprintf(
            "Error rate exceeded: %d errors in last minute",
            errorCount,
        ))
    }
}
```

### 5. 컨텍스트 로깅

```go
// 요청 전체에서 사용할 로거
func WithRequestLogger(c *gin.Context) {
    logger := log.WithFields(log.Fields{
        "request_id": c.GetString("RequestID"),
        "method":     c.Request.Method,
        "path":       c.Request.URL.Path,
    })
    c.Set("Logger", logger)
}

// 핸들러에서 사용
func handler(c *gin.Context) {
    logger := c.MustGet("Logger").(*log.Entry)
    logger.Info("Processing request")
}
```

## 📚 학습 완료

에러 처리 & 로깅 섹션의 3개 항목을 모두 완료했습니다:
- ✅ 09. HTTP 상태코드와 에러 응답 규약
- ✅ 10. 에러 핸들링 미들웨어
- ✅ 11. 로깅 미들웨어

## 🔗 참고 자료
- [Structured Logging](https://www.honeybadger.io/blog/golang-logging/)
- [Gin Logger Middleware](https://gin-gonic.com/docs/examples/custom-log-format/)
- [Logrus - Structured Logger for Go](https://github.com/sirupsen/logrus)
- [Zap - Uber's Blazing Fast Logger](https://github.com/uber-go/zap)