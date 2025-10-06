# 서버의 모든 것을 기록하기 📝

안녕하세요! 서버가 잘 돌아가는지, 문제는 없는지 확인하려면 **로그**가 필요해요. 로그는 서버의 일기장이라고 생각하면 됩니다. 무슨 일이 있었는지 모두 기록해두는 거죠!

## 로깅이 뭔가요?

로깅은 **서버에서 일어나는 모든 일을 기록**하는 거예요. 누가 언제 접속했는지, 어떤 요청을 보냈는지, 에러는 없었는지 등을 파일에 남겨두는 거죠.

### 실생활 비유
- **CCTV**: 가게에 누가 왔는지, 무엇을 샀는지 기록
- **비행기 블랙박스**: 비행 중 모든 정보를 기록해서 사고 시 원인 파악
- **병원 진료 기록**: 환자의 증상, 처방, 치료 과정을 모두 기록

### 왜 로그가 중요할까요?
```
로그가 없으면: "어? 서버가 느린데 왜 그러지?" 🤷
로그가 있으면: "아! 10시 30분에 /api/users 요청이 3초 걸렸네!" 💡
```

## 로그 레벨 이해하기

로그에는 **중요도**가 있어요!

- **DEBUG**: 개발할 때 디버깅용 (변수 값 확인 등)
- **INFO**: 일반적인 정보 ("사용자 로그인 성공")
- **WARN**: 경고! 문제가 될 수 있음 ("응답 시간 느림")
- **ERROR**: 에러 발생! ("데이터베이스 연결 실패")

**실생활 비유**: 병원의 환자 상태
- DEBUG: 혈압, 체온 (의사만 확인)
- INFO: 정상 상태
- WARN: 주의 필요 (경미한 이상)
- ERROR: 긴급 치료 필요

## 이번 챕터에서 배울 내용
- 모든 요청과 응답을 로그로 남기기
- 로그를 JSON 형식으로 깔끔하게 정리하기
- 느린 요청을 자동으로 감지하기
- 민감한 정보(비밀번호 등)는 숨기기
- 로그 파일 자동 관리하기 (너무 커지지 않게)

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

## 💡 꼭 알아야 할 핵심 개념!

### 1. 환경에 맞게 로그 레벨 조정하기

모든 로그를 다 기록하면 너무 많아요! **환경에 따라** 필요한 것만 기록하세요.

```go
type LogLevel int

const (
    DEBUG LogLevel = iota  // 🔍 개발할 때: 모든 정보
    INFO                   // 📝 일반 운영: 중요한 정보
    WARN                   // ⚠️ 경고: 문제가 될 수 있는 것
    ERROR                  // ❌ 에러: 반드시 확인해야 함
    FATAL                  // 💥 치명적: 서버가 멈출 수 있음
)

// 환경별 설정
// 개발 환경: DEBUG (모든 로그 보기)
// 테스트 환경: INFO (중요한 것만)
// 실제 서버: WARN (문제만)
```

**실생활 비유**:
- **개발**: 실험실에서 모든 과정 기록
- **운영**: 공장에서 중요한 사건만 기록
- **문제 발생**: CCTV 녹화 보듯이 로그 확인

### 2. JSON 형식으로 깔끔하게 기록하기

로그를 **구조화**하면 나중에 찾기 쉬워요!

```go
// ✅ 좋은 예: JSON 형식 (프로그램이 읽기 쉬움)
{
    "level": "INFO",
    "timestamp": "2024-01-01T10:00:00Z",
    "request_id": "req-123",
    "method": "GET",
    "path": "/api/users",
    "status_code": 200,
    "latency": "5ms",
    "client_ip": "127.0.0.1"
}

// ❌ 나쁜 예: 그냥 텍스트 (사람은 읽기 쉽지만 검색 어려움)
2024-01-01 10:00:00 INFO GET /api/users 200 5ms 127.0.0.1
```

**왜 JSON이 좋을까요?**
- 나중에 "에러만" 검색하기 쉬움
- 분석 도구로 통계 내기 쉬움
- 다른 시스템으로 전송하기 쉬움

### 3. 민감한 정보는 절대 기록하지 마세요!

비밀번호, 토큰 같은 **민감한 정보**가 로그에 남으면 큰일나요!

```go
// 절대 기록하면 안 되는 것들
sensitiveHeaders := []string{
    "Authorization",   // 인증 토큰
    "Cookie",         // 세션 쿠키
    "X-Api-Key",      // API 키
}

sensitiveFields := []string{
    "password",       // 비밀번호
    "token",         // 토큰
    "secret",        // 비밀 키
    "credit_card",   // 카드 번호
}

// 로그에 기록할 때 마스킹
{
    "username": "hong",
    "password": "***MASKED***"  // 실제 값 대신 이렇게!
}
```

**실생활 비유**: CCTV 녹화본에서 개인정보를 모자이크 처리하는 것!

### 4. 로그 파일 자동 관리하기

로그가 계속 쌓이면 **하드디스크가 꽉 차요**! 자동으로 관리해야 해요.

```go
// 파일이 너무 크면 새 파일로 교체
if fileSize >= 100MB {
    // 1. 기존 파일 백업 (날짜 붙이기)
    os.Rename("app.log", "app.log.20240101")

    // 2. 새 파일 시작
    createNewLogFile()

    // 3. 오래된 백업은 삭제
    deleteOldLogs(30일전)
}
```

**로그 로테이션 전략**:
- **크기 기반**: 100MB 넘으면 새 파일
- **시간 기반**: 매일 자정에 새 파일
- **보관 기간**: 30일 지난 로그는 삭제

**실생활 비유**: 일기장이 다 차면 새 일기장을 사고, 10년 전 일기는 버리는 것!

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