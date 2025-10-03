# 🏦 Banking System - Advanced Gin Project (Lessons 09-17)

> Production-ready banking system implementing advanced Gin concepts including error handling, logging, configuration management, dependency injection, database operations, migrations, and transaction management

## 📌 프로젝트 소개

이 프로젝트는 Gin 프레임워크의 고급 기능들(Lessons 09-17)을 실전에서 활용하는 방법을 보여주는 완전한 뱅킹 시스템입니다. 실제 프로덕션 환경에서 사용할 수 있는 수준의 코드 구조와 패턴을 제공합니다.

### 🎯 학습 목표
- **Lesson 09**: HTTP 상태 코드와 표준 에러 응답
- **Lesson 10**: 전역 에러 핸들링과 패닉 복구
- **Lesson 11**: 구조화된 로깅과 요청 추적
- **Lesson 12**: Viper를 통한 설정 관리
- **Lesson 13**: 의존성 주입과 클린 아키텍처
- **Lesson 14**: 실행 모드별 최적화 (Debug/Release/Test)
- **Lesson 15**: GORM을 사용한 CRUD 작업
- **Lesson 16**: 마이그레이션과 시드 데이터
- **Lesson 17**: 트랜잭션과 컨텍스트 타임아웃

## 🏗 시스템 아키텍처

```
┌─────────────────────────────────────────┐
│           API Gateway Layer             │
│        (Routes & Middleware)            │
├─────────────────────────────────────────┤
│           Handler Layer                 │
│     (HTTP Request/Response)            │
├─────────────────────────────────────────┤
│           Service Layer                 │
│      (Business Logic & DI)             │
├─────────────────────────────────────────┤
│          Repository Layer               │
│         (Data Access)                   │
├─────────────────────────────────────────┤
│           Database Layer                │
│      (SQLite with GORM)               │
└─────────────────────────────────────────┘
```

## 📁 프로젝트 구조

```
project-advanced/
├── cmd/
│   └── main.go                 # 애플리케이션 진입점
├── internal/
│   ├── config/                 # 설정 관리 (Lesson 12)
│   │   └── config.go
│   ├── handlers/               # HTTP 핸들러
│   │   ├── health.go
│   │   ├── user.go
│   │   ├── account.go
│   │   ├── transaction.go
│   │   └── admin.go
│   ├── middleware/             # 미들웨어 (Lessons 10, 11, 14)
│   │   └── middleware.go
│   ├── models/                 # 데이터 모델 (Lesson 15)
│   │   └── models.go
│   └── services/               # 비즈니스 로직 (Lesson 13)
│       ├── container.go        # DI 컨테이너
│       ├── user_service.go
│       ├── account_service.go
│       ├── transaction_service.go  # (Lesson 17)
│       ├── migration_service.go    # (Lesson 16)
│       └── seeder_service.go
├── pkg/
│   ├── database/               # 데이터베이스 초기화
│   │   └── database.go
│   ├── logger/                 # 로거 구현 (Lesson 11)
│   │   └── logger.go
│   └── validator/              # 입력 검증
│       └── validator.go
├── migrations/                 # 데이터베이스 마이그레이션 (Lesson 16)
├── seeds/                      # 시드 데이터
├── docs/                       # API 문서
├── config.yaml                 # 기본 설정 파일
├── config.development.yaml     # 개발 환경 설정
├── config.production.yaml      # 프로덕션 환경 설정
├── .env.example               # 환경변수 템플릿
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🚀 빠른 시작

### 1. 필수 요구사항
- Go 1.21 이상
- SQLite3
- Make (선택사항)

### 2. 설치 및 실행

```bash
# 프로젝트 클론
cd gin/project-advanced

# 의존성 설치
go mod download

# 환경변수 설정
cp .env.example .env

# 애플리케이션 실행 (개발 모드)
go run cmd/main.go

# 또는 Makefile 사용
make run

# 프로덕션 모드 실행
APP_ENV=production go run cmd/main.go
```

### 3. 초기 데이터
애플리케이션이 시작되면 자동으로:
- 데이터베이스 테이블 생성 (Auto-migration)
- 개발 환경에서 시드 데이터 생성
- 5개 테스트 사용자 계정
- 각 사용자당 2개 은행 계좌

## 💡 핵심 기능 구현

### 1. 에러 처리 시스템 (Lessons 9-10)

#### 표준화된 에러 응답
```go
// 모든 에러는 일관된 형식으로 반환
{
  "error": "Insufficient balance",
  "request_id": "req_1234567890",
  "timestamp": 1642345678,
  "details": {  // Debug 모드에서만
    "account_id": 123,
    "required": 1000,
    "available": 500
  }
}
```

#### 전역 에러 핸들러
```go
// middleware/middleware.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last()

            // 에러 타입별 처리
            switch e := err.Err.(type) {
            case ValidationError:
                c.JSON(400, formatValidationError(e))
            case BusinessError:
                c.JSON(422, formatBusinessError(e))
            default:
                c.JSON(500, formatInternalError(e))
            }
        }
    }
}
```

### 2. 구조화된 로깅 (Lesson 11)

#### JSON 형식 로깅
```go
logger.WithFields(map[string]interface{}{
    "request_id": requestID,
    "user_id": userID,
    "action": "transfer",
    "amount": 1000,
}).Info("Transfer initiated")
```

#### 로그 출력 예시
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "level": "INFO",
  "message": "Transfer initiated",
  "request_id": "req_1234567890",
  "user_id": 42,
  "action": "transfer",
  "amount": 1000
}
```

### 3. 설정 관리 (Lesson 12)

#### 계층적 설정
```yaml
# config.yaml
app:
  environment: ${APP_ENV:development}
  mode: ${APP_MODE:debug}

database:
  dsn: ${DATABASE_URL:banking.db}

security:
  jwt_secret: ${JWT_SECRET:secret}
```

#### 환경변수 우선순위
1. 환경변수 (최우선)
2. 환경별 설정 파일
3. 기본 설정 파일
4. 하드코딩된 기본값

### 4. 의존성 주입 (Lesson 13)

#### Clean Architecture 패턴
```go
// DI Container
type Container struct {
    DB     *gorm.DB
    Config *config.Config
    Logger *logger.Logger
}

// Service with injected dependencies
type TransactionService struct {
    container *Container
}

func NewTransactionService(c *Container) *TransactionService {
    return &TransactionService{container: c}
}
```

### 5. 실행 모드 (Lesson 14)

#### Debug Mode
- 상세 로깅
- 프로파일링 엔드포인트 활성화
- 에러 스택 트레이스
- 메모리 통계

#### Release Mode
- 보안 헤더
- Rate Limiting
- 최소 로깅
- 에러 정보 마스킹

#### Test Mode
- Mock 데이터
- 빠른 실행
- 테스트 헬퍼

### 6. GORM CRUD (Lesson 15)

#### Repository 패턴
```go
// 페이지네이션과 필터링
func (r *AccountRepository) List(userID uint, offset, limit int) ([]*Account, int64, error) {
    var accounts []*Account
    var total int64

    query := r.db.Model(&Account{}).Where("user_id = ?", userID)
    query.Count(&total)

    err := query.Offset(offset).Limit(limit).Find(&accounts).Error
    return accounts, total, err
}
```

### 7. 마이그레이션 시스템 (Lesson 16)

#### 버전 관리
```go
migrations := []Migration{
    {
        Version: "001_create_users",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&User{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&User{})
        },
    },
}
```

### 8. 트랜잭션 처리 (Lesson 17)

#### ACID 보장 송금
```go
func (s *TransactionService) Transfer(ctx context.Context, from, to uint, amount float64) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. 비관적 잠금으로 계좌 조회
        var fromAccount Account
        tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromAccount, from)

        // 2. 잔액 확인
        if fromAccount.Balance < amount {
            return ErrInsufficientBalance
        }

        // 3. 잔액 업데이트
        fromAccount.Balance -= amount
        toAccount.Balance += amount

        // 4. 저장
        tx.Save(&fromAccount)
        tx.Save(&toAccount)

        return nil
    }, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
}
```

## 🔧 API 엔드포인트

### 인증 및 사용자 관리
```bash
POST   /api/v1/register           # 회원가입
POST   /api/v1/login              # 로그인
GET    /api/v1/users              # 사용자 목록
GET    /api/v1/users/:id          # 사용자 상세
PUT    /api/v1/users/:id          # 사용자 수정
DELETE /api/v1/users/:id          # 사용자 삭제
```

### 계좌 관리
```bash
POST   /api/v1/accounts           # 계좌 생성
GET    /api/v1/accounts           # 계좌 목록
GET    /api/v1/accounts/:id       # 계좌 상세
POST   /api/v1/accounts/:id/deposit  # 입금
POST   /api/v1/accounts/:id/withdraw # 출금
GET    /api/v1/accounts/:id/balance  # 잔액 조회
```

### 거래 관리
```bash
POST   /api/v1/transactions/transfer  # 송금
GET    /api/v1/transactions          # 거래 내역
GET    /api/v1/transactions/:id      # 거래 상세
GET    /api/v1/transactions/report   # 거래 리포트
```

### 관리자 기능
```bash
GET    /admin/migrations          # 마이그레이션 상태
POST   /admin/migrations/run      # 마이그레이션 실행
POST   /admin/migrations/rollback # 롤백
POST   /admin/seed                # 시드 데이터 생성
GET    /admin/metrics             # 시스템 메트릭
GET    /admin/config              # 설정 조회
```

### 디버그 엔드포인트 (Debug Mode)
```bash
GET    /debug/pprof/*             # 프로파일링
GET    /debug/vars                # 런타임 변수
GET    /debug/routes              # 라우트 목록
```

## 📝 API 사용 예제

### 1. 회원가입
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 2. 계좌 생성
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "account_type": "savings",
    "currency": "USD"
  }'
```

### 3. 송금 (트랜잭션)
```bash
curl -X POST http://localhost:8080/api/v1/transactions/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100.50,
    "description": "Monthly rent"
  }'

# 응답
{
  "transaction_id": "TXN1234567890",
  "status": "completed",
  "amount": 100.50,
  "processing_time_ms": 125,
  "timestamp": "2024-01-15T10:30:45Z"
}
```

### 4. 타임아웃 설정 (Context)
```bash
# 5초 타임아웃 설정
curl -X POST http://localhost:8080/api/v1/transactions/transfer \
  -H "Content-Type: application/json" \
  -H "X-Request-Timeout: 5000" \
  -d '{...}'
```

### 5. 동시성 테스트
```bash
# 10개 동시 요청 테스트
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/transactions/transfer \
    -H "Content-Type: application/json" \
    -d '{"from_account_id": 1, "to_account_id": 2, "amount": 10}' &
done
```

## 🧪 테스트

### 유닛 테스트
```bash
go test ./internal/services/... -v
```

### 통합 테스트
```bash
go test ./internal/handlers/... -v
```

### 부하 테스트
```bash
# Apache Bench 사용
ab -n 1000 -c 10 http://localhost:8080/health

# Hey 사용
hey -n 1000 -c 10 http://localhost:8080/api/v1/accounts
```

## 📊 모니터링

### 메트릭 수집
```bash
# Prometheus 형식 메트릭
curl http://localhost:8080/admin/metrics

# 응답 예시
http_requests_total{method="GET",path="/api/v1/accounts",status="200"} 142
http_request_duration_seconds{quantile="0.99"} 0.05
go_goroutines 12
go_memstats_alloc_bytes 2.5e+06
```

### 헬스체크
```bash
# 애플리케이션 헬스
curl http://localhost:8080/health

# 데이터베이스 헬스
curl http://localhost:8080/health/db

# Ready 상태
curl http://localhost:8080/ready
```

## 🔒 보안 고려사항

### 구현된 보안 기능
- ✅ SQL 인젝션 방지 (GORM 파라미터 바인딩)
- ✅ XSS 방지 (보안 헤더)
- ✅ CSRF 보호 (선택적)
- ✅ Rate Limiting
- ✅ 비밀번호 해싱 (bcrypt)
- ✅ JWT 토큰 인증
- ✅ 민감정보 로깅 마스킹
- ✅ HTTPS 강제 (프로덕션)

### 프로덕션 체크리스트
- [ ] 환경변수로 모든 시크릿 관리
- [ ] TLS/SSL 인증서 설정
- [ ] 데이터베이스 백업 전략
- [ ] 로그 수집 및 분석
- [ ] 알림 시스템 구성
- [ ] 장애 복구 계획
- [ ] 성능 모니터링
- [ ] 보안 감사

## 🚀 배포

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o banking cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/banking .
COPY --from=builder /app/config*.yaml .
EXPOSE 8080
CMD ["./banking"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DATABASE_URL=postgres://...
    volumes:
      - ./logs:/app/logs
```

## 📈 성능 최적화

### 구현된 최적화
- **연결 풀링**: 데이터베이스 연결 재사용
- **컨텍스트 타임아웃**: 장시간 실행 방지
- **비동기 처리**: 고루틴 활용
- **캐싱**: 자주 조회되는 데이터 캐싱
- **인덱싱**: 데이터베이스 인덱스 최적화
- **페이지네이션**: 대량 데이터 처리
- **배치 처리**: 트랜잭션 배치 실행

## 📚 학습 포인트

### Lesson 9: HTTP 상태 코드
- 적절한 상태 코드 선택
- 표준화된 에러 응답 구조
- 비즈니스 에러 vs 시스템 에러

### Lesson 10: 에러 핸들링
- 패닉 복구 미들웨어
- 전역 에러 핸들러
- 에러 타입별 처리

### Lesson 11: 로깅
- 구조화된 로그 포맷
- 로그 레벨 관리
- 요청 추적 (Request ID)
- 민감정보 마스킹

### Lesson 12: 설정 관리
- Viper 활용
- 환경별 설정 분리
- 환경변수 우선순위
- 설정 검증

### Lesson 13: 의존성 주입
- Constructor Injection
- Interface 기반 설계
- 테스트 가능한 구조
- Clean Architecture

### Lesson 14: 실행 모드
- Debug/Release/Test 모드
- 모드별 최적화
- 프로파일링
- 보안 설정

### Lesson 15: GORM CRUD
- Repository 패턴
- 관계 설정
- Preloading
- 페이지네이션

### Lesson 16: 마이그레이션
- 버전 관리
- Up/Down 마이그레이션
- 시드 데이터
- 롤백 전략

### Lesson 17: 트랜잭션
- ACID 보장
- 비관적/낙관적 잠금
- 컨텍스트 타임아웃
- 동시성 제어

## 🤝 기여하기

이 프로젝트는 학습 목적으로 만들어졌습니다. 개선 사항이나 버그를 발견하면 이슈를 생성하거나 PR을 보내주세요!

## 📄 라이선스

MIT License

## 🎯 다음 단계

이 프로젝트를 완료했다면:
1. JWT 인증 구현
2. Redis 캐싱 추가
3. WebSocket 실시간 알림
4. gRPC 서비스 추가
5. Kubernetes 배포
6. CI/CD 파이프라인 구성

---

**Happy Banking! 🏦**