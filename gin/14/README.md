# Lesson 14: 실행 모드 (Release/Debug/Test) 🎯

> Gin 애플리케이션의 환경별 최적화와 실행 모드 관리 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

애플리케이션은 개발, 테스트, 프로덕션 환경에서 다르게 동작해야 합니다. Gin은 세 가지 실행 모드를 제공하며, 각 모드별로 최적화된 설정과 기능을 활용할 수 있습니다. 이번 레슨에서는 실행 모드를 효과적으로 관리하는 방법을 학습합니다.

### 핵심 학습 목표
- ✅ Debug, Release, Test 모드의 차이점
- ✅ 모드별 미들웨어와 설정 최적화
- ✅ 프로파일링과 디버깅 도구 활용
- ✅ 메트릭스 수집과 모니터링
- ✅ 리소스 제한과 성능 튜닝
- ✅ 환경별 에러 처리 전략

## 🎨 실행 모드 비교

| 특성 | Debug Mode | Release Mode | Test Mode |
|------|------------|--------------|-----------|
| **로깅 레벨** | Debug | Info | Error |
| **에러 상세** | 스택 트레이스 포함 | 일반 메시지만 | 테스트 정보 포함 |
| **프로파일링** | ✅ 활성화 | ❌ 비활성화 | ❌ 비활성화 |
| **Swagger UI** | ✅ 활성화 | ❌ 비활성화 | ❌ 비활성화 |
| **요청 로깅** | ✅ 상세 | ✅ 간단 | ❌ 비활성화 |
| **컬러 출력** | ✅ 활성화 | ❌ 비활성화 | ❌ 비활성화 |
| **패닉 복구** | ✅ 활성화 | ✅ 활성화 | ❌ 비활성화 |
| **Rate Limiting** | ❌ 없음 | ✅ 100 req/min | ❌ 없음 |
| **보안 헤더** | ❌ 없음 | ✅ 모두 포함 | ❌ 없음 |
| **메모리 제한** | 무제한 | 1GB | 256MB |
| **타임아웃** | 30초 | 15초 | 5초 |

## 🛠 구현된 기능

### 1. **모드별 설정 구조**
```go
type ModeConfig struct {
    Mode            RunMode
    LogLevel        string
    EnableProfiling bool
    EnableMetrics   bool
    EnableSwagger   bool
    ErrorDetails    bool
    PanicRecovery   bool
    RequestLogging  bool
    ColoredOutput   bool
    MaxMemory       int64
    MaxCPU          int
    RateLimit       int
    Timeout         time.Duration
}
```

### 2. **모드별 미들웨어**
- **Debug**: 상세 로깅, 메모리 추적, 요청 시간 측정
- **Release**: 보안 헤더, Rate Limiting, 민감정보 제거
- **Test**: 테스트 ID 생성, 간소화된 응답

### 3. **디버깅 도구**
- `/debug/pprof/*` - CPU/메모리 프로파일링
- `/debug/vars` - 런타임 변수
- `/debug/gc` - 가비지 컬렉션 트리거
- `/debug/mem` - 메모리 통계

### 4. **모니터링 엔드포인트**
- `/metrics` - Prometheus 형식 메트릭스
- `/health` - 헬스체크
- `/mode` - 현재 모드 정보

## 🎯 주요 API 엔드포인트

### 공통 엔드포인트
```bash
GET  /health          # 헬스체크
GET  /mode            # 모드 정보
GET  /api/users       # 샘플 API
GET  /api/error       # 에러 테스트
GET  /api/panic       # 패닉 테스트
GET  /api/slow        # 느린 응답 테스트
GET  /api/memory      # 메모리 사용 테스트
```

### Debug 모드 전용
```bash
GET  /debug/pprof/           # 프로파일링 인덱스
GET  /debug/pprof/profile    # CPU 프로파일
GET  /debug/pprof/heap       # 힙 프로파일
GET  /debug/pprof/trace      # 실행 트레이스
GET  /debug/vars             # 런타임 변수
GET  /debug/gc               # GC 트리거
GET  /debug/mem              # 메모리 통계
GET  /debug/config           # 설정 정보
GET  /debug/routes           # 라우트 목록
GET  /debug/env              # 환경변수
POST /mode/:mode             # 모드 전환 (재시작 필요)
GET  /swagger/*              # Swagger UI
```

### Release 모드 전용
```bash
GET  /metrics         # Prometheus 메트릭스
```

### Test 모드 전용
```bash
POST /test/reset      # 테스트 데이터 리셋
POST /test/seed       # 테스트 데이터 시딩
```

## 💻 실습 가이드

### 1. 모드별 실행

#### Debug 모드
```bash
# 환경변수로 설정
GIN_MODE=debug go run main.go

# 또는 기본값 (환경변수 없을 때)
go run main.go
```

#### Release 모드
```bash
GIN_MODE=release go run main.go

# 프로덕션 배포 예시
GIN_MODE=release ./app
```

#### Test 모드
```bash
GIN_MODE=test go run main.go

# 테스트 실행 시
GIN_MODE=test go test ./...
```

### 2. 모드 정보 확인
```bash
# 현재 모드 확인
curl http://localhost:8080/mode | jq

# 응답 예시 (Debug 모드)
{
  "mode": "debug",
  "debug": true,
  "release": false,
  "test": false,
  "profiling": true,
  "metrics": true,
  "swagger": true,
  "error_details": true,
  "request_logging": true,
  "colored_output": true
}
```

### 3. Debug 모드 기능 테스트

#### 프로파일링
```bash
# CPU 프로파일 (30초)
go tool pprof http://localhost:8080/debug/pprof/profile

# 메모리 프로파일
go tool pprof http://localhost:8080/debug/pprof/heap

# 고루틴 확인
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 실행 트레이스
wget http://localhost:8080/debug/pprof/trace
go tool trace trace
```

#### 메모리 통계
```bash
# 메모리 상태
curl http://localhost:8080/debug/mem | jq

# GC 트리거
curl -X GET http://localhost:8080/debug/gc

# 런타임 변수
curl http://localhost:8080/debug/vars | jq
```

#### 라우트 정보
```bash
# 등록된 라우트 목록
curl http://localhost:8080/debug/routes | jq

# 설정 정보
curl http://localhost:8080/debug/config | jq

# 환경변수
curl http://localhost:8080/debug/env | jq
```

### 4. Release 모드 기능 테스트

#### 메트릭스 수집
```bash
# Prometheus 형식 메트릭스
curl http://localhost:8080/metrics

# 출력 예시
# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 142

# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines 8
```

#### 보안 헤더 확인
```bash
# Release 모드에서 보안 헤더 확인
curl -I http://localhost:8080/health

# 응답 헤더
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

### 5. 에러 처리 테스트

#### Debug 모드 에러
```bash
curl http://localhost:8080/api/error | jq

# 상세한 에러 정보
{
  "error": "test error",
  "type": "*errors.errorString",
  "stack_trace": "goroutine 1 [running]:\n...",
  "request_id": "abc123",
  "path": "/api/error",
  "method": "GET"
}
```

#### Release 모드 에러
```bash
curl http://localhost:8080/api/error | jq

# 간소화된 에러 정보
{
  "error": "Internal server error",
  "request_id": "abc123"
}
```

### 6. Test 모드 기능
```bash
# 테스트 데이터 리셋
curl -X POST http://localhost:8080/test/reset

# 테스트 데이터 시딩
curl -X POST http://localhost:8080/test/seed

# 테스트 ID 확인
curl http://localhost:8080/api/users -H "X-Test-Mode: true"
```

## 🔍 코드 하이라이트

### 모드별 라우터 초기화
```go
func NewApplication(mode RunMode) *Application {
    config := GetModeConfig(mode)

    // 리소스 제한 설정
    if config.MaxMemory > 0 {
        debug.SetMemoryLimit(config.MaxMemory)
    }

    if config.MaxCPU > 0 {
        runtime.GOMAXPROCS(config.MaxCPU)
    }

    // 모드별 라우터 생성
    var router *gin.Engine
    switch mode {
    case DebugMode:
        router = SetupDebugRouter(config)
    case ReleaseMode:
        router = SetupReleaseRouter(config)
    case TestMode:
        router = SetupTestRouter(config)
    }

    return &Application{
        Router: router,
        Config: config,
        Mode:   mode,
    }
}
```

### Debug 미들웨어
```go
func DebugMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 상세 로깅
        log.Printf("[DEBUG] %s %s from %s",
            c.Request.Method, c.Request.URL.Path, c.ClientIP())

        // 메모리 추적
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        log.Printf("[DEBUG] Memory - Alloc: %v MB, NumGC: %v",
            m.Alloc/1024/1024, m.NumGC)

        // 요청 시간 측정
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        log.Printf("[DEBUG] Request completed in %v", latency)
    }
}
```

### Release 보안 설정
```go
func ReleaseMiddleware(config *ModeConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 보안 헤더
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security",
            "max-age=31536000; includeSubDomains")

        c.Next()

        // 서버 정보 숨기기
        c.Header("Server", "")
    }
}
```

## 📝 베스트 프랙티스

### 1. **환경별 설정 분리**
```go
// ❌ Bad: 하드코딩된 조건
if os.Getenv("ENV") == "prod" {
    // production code
}

// ✅ Good: 모드 기반 설정
config := GetModeConfig(mode)
app := NewApplication(config)
```

### 2. **민감정보 보호**
```go
// Debug 모드에서만 상세 정보
if mode == DebugMode {
    return detailedError
} else {
    return genericError
}
```

### 3. **리소스 제한**
```go
// 환경별 리소스 제한
switch mode {
case TestMode:
    debug.SetMemoryLimit(256 << 20) // 256MB
    runtime.GOMAXPROCS(2)
case ReleaseMode:
    debug.SetMemoryLimit(1 << 30) // 1GB
    runtime.GOMAXPROCS(runtime.NumCPU())
}
```

### 4. **조건부 미들웨어**
```go
if config.EnableProfiling {
    router.Use(ProfilingMiddleware())
}

if config.RequestLogging {
    router.Use(gin.Logger())
}
```

## 🚀 프로덕션 체크리스트

### Debug 모드
- [ ] 모든 디버깅 도구 활성화
- [ ] 상세한 로깅 설정
- [ ] 프로파일링 엔드포인트 접근 가능
- [ ] 에러 스택 트레이스 포함

### Release 모드
- [ ] 디버깅 엔드포인트 비활성화
- [ ] 보안 헤더 설정
- [ ] Rate limiting 활성화
- [ ] 민감정보 마스킹
- [ ] 메트릭스 수집 활성화
- [ ] 리소스 제한 설정
- [ ] 로그 레벨 적절히 설정

### Test 모드
- [ ] 불필요한 로깅 비활성화
- [ ] 패닉 복구 비활성화 (테스트 실패 감지)
- [ ] 테스트 헬퍼 엔드포인트 활성화
- [ ] 짧은 타임아웃 설정

## 🎨 성능 튜닝 팁

### CPU 프로파일링
```bash
# 1. 프로파일 수집
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# 2. 분석
go tool pprof cpu.prof

# 3. 웹 UI로 확인
go tool pprof -http=:6060 cpu.prof
```

### 메모리 프로파일링
```bash
# 1. 힙 프로파일
curl http://localhost:8080/debug/pprof/heap > heap.prof

# 2. 분석
go tool pprof heap.prof

# 3. 메모리 누수 확인
go tool pprof -alloc_space heap.prof
```

### 실행 트레이스
```bash
# 1. 트레이스 수집
curl http://localhost:8080/debug/pprof/trace?seconds=5 > trace.out

# 2. 분석
go tool trace trace.out
```

## 📚 추가 학습 자료

- [Gin Mode Configuration](https://gin-gonic.com/docs/examples/run-multiple-service/)
- [Go Profiling Guide](https://go.dev/blog/pprof)
- [Production-Ready Go](https://www.oreilly.com/library/view/production-go/9781788993746/)
- [Monitoring Go Applications](https://prometheus.io/docs/guides/go-application/)

## 🎯 정리

실행 모드 관리는 애플리케이션의 안정성과 성능을 좌우하는 중요한 요소입니다. 각 환경에 맞는 최적화된 설정으로 개발 생산성을 높이고, 프로덕션 안정성을 확보할 수 있습니다.

**핵심 포인트:**
- Debug 모드에서 충분히 테스트하고 프로파일링
- Release 모드에서 보안과 성능 최적화
- Test 모드에서 빠르고 정확한 테스트 실행
- 환경별 설정을 명확히 분리

이제 환경에 따라 최적화된 Gin 애플리케이션을 운영할 수 있습니다! 🚀