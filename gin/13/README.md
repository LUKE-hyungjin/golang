# Lesson 13: 의존성 주입 (Dependency Injection) 💉

> 테스트 가능하고 유지보수가 쉬운 Go 애플리케이션을 위한 DI 패턴 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

의존성 주입(DI)은 객체 간의 결합도를 낮추고 테스트 가능성을 높이는 핵심 설계 패턴입니다. Go에서는 인터페이스와 구조체를 활용하여 우아하게 DI를 구현할 수 있습니다. 이번 레슨에서는 실전에서 사용할 수 있는 다양한 DI 패턴을 학습합니다.

### 핵심 학습 목표
- ✅ 인터페이스 기반 설계 원칙
- ✅ Constructor Injection 패턴
- ✅ Factory 패턴과 Container
- ✅ 의존성 역전 원칙 (DIP)
- ✅ Mock을 활용한 테스트 전략
- ✅ 환경별 구현체 교체

## 🏗 아키텍처 설계

### Clean Architecture 레이어
```
┌─────────────────────────────────────────┐
│           Presentation Layer            │
│         (HTTP Handlers, Routes)         │
├─────────────────────────────────────────┤
│           Business Logic Layer          │
│      (Services, Use Cases, Rules)       │
├─────────────────────────────────────────┤
│            Domain Layer                 │
│       (Entities, Value Objects)         │
├─────────────────────────────────────────┤
│         Infrastructure Layer            │
│    (Repositories, External Services)    │
└─────────────────────────────────────────┘
```

### 의존성 방향
```
Handler → Service → Repository
   ↓         ↓          ↓
Interface Interface Interface
```

## 🛠 구현된 DI 패턴

### 1. **인터페이스 정의 (Port)**
```go
// Repository 인터페이스
type UserRepository interface {
    FindByID(ctx context.Context, id int) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}

// Service 인터페이스
type UserService interface {
    GetUser(ctx context.Context, id int) (*User, error)
    CreateUser(ctx context.Context, email, name string) (*User, error)
    UpdateUser(ctx context.Context, id int, name string) (*User, error)
}

// 외부 서비스 인터페이스
type EmailService interface {
    SendEmail(to, subject, body string) error
    SendOrderConfirmation(order *Order, user *User) error
}
```

### 2. **Constructor Injection**
```go
type UserServiceImpl struct {
    userRepo UserRepository  // 인터페이스 의존
    cache    CacheService    // 인터페이스 의존
    email    EmailService    // 인터페이스 의존
}

// 생성자를 통한 의존성 주입
func NewUserService(
    userRepo UserRepository,
    cache CacheService,
    email EmailService,
) UserService {
    return &UserServiceImpl{
        userRepo: userRepo,
        cache:    cache,
        email:    email,
    }
}
```

### 3. **Factory Pattern & Container**
```go
type Container struct {
    config         *Config
    db             *sql.DB
    userRepository UserRepository
    userService    UserService
    emailService   EmailService
    cacheService   CacheService
}

// Lazy initialization with singleton
func (c *Container) GetUserService() UserService {
    if c.userService == nil {
        c.userService = NewUserService(
            c.GetUserRepository(),
            c.GetCacheService(),
            c.GetEmailService(),
        )
    }
    return c.userService
}
```

### 4. **환경별 구현체 교체**
```go
func (c *Container) GetEmailService() EmailService {
    if c.emailService == nil {
        if c.config.Environment == "test" {
            // 테스트 환경: Mock 서비스
            c.emailService = NewMockEmailService()
        } else {
            // 프로덕션 환경: 실제 SMTP 서비스
            c.emailService = NewSMTPEmailService(
                c.config.SMTPHost,
                c.config.SMTPPort,
            )
        }
    }
    return c.emailService
}
```

## 🎯 주요 API 엔드포인트

### 사용자 관리 API
```bash
GET    /users/:id      # 사용자 조회
POST   /users          # 사용자 생성
PUT    /users/:id      # 사용자 수정
DELETE /users/:id      # 사용자 삭제
GET    /users          # 사용자 목록
```

### DI 정보 API
```bash
GET    /di/info             # DI 패턴 정보
GET    /patterns/constructor # Constructor Injection 예제
GET    /patterns/factory    # Factory 패턴 예제
GET    /patterns/interface  # Interface Segregation 예제
```

## 💻 실습 가이드

### 1. 기본 실행
```bash
# 개발 환경 (Mock 서비스 사용)
APP_ENV=development go run main.go

# 테스트 환경 (In-Memory 구현체)
APP_ENV=test go run main.go

# 프로덕션 환경 (실제 서비스)
APP_ENV=production DATABASE_URL=postgres://... go run main.go
```

### 2. API 테스트

#### 사용자 생성
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "role": "admin"
  }'
```

#### 사용자 조회
```bash
# 단일 사용자 조회
curl http://localhost:8080/users/1

# 사용자 목록 조회
curl "http://localhost:8080/users?page=1&page_size=10"
```

#### 사용자 수정
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'
```

#### 사용자 삭제
```bash
curl -X DELETE http://localhost:8080/users/1
```

#### DI 패턴 정보 조회
```bash
# DI 정보
curl http://localhost:8080/di/info

# Constructor Injection 패턴
curl http://localhost:8080/patterns/constructor

# Factory 패턴
curl http://localhost:8080/patterns/factory

# Interface Segregation
curl http://localhost:8080/patterns/interface
```

## 🔍 코드 하이라이트

### 서비스 계층의 비즈니스 로직
```go
func (s *UserServiceImpl) CreateUser(ctx context.Context, email, name string) (*User, error) {
    // 1. 도메인 객체 생성
    user := &User{
        Email: email,
        Name:  name,
    }

    // 2. Repository를 통한 저장 (인터페이스 의존)
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // 3. 외부 서비스 호출 (인터페이스 의존)
    s.email.SendEmail(email, "Welcome!", fmt.Sprintf("Welcome %s!", name))

    return user, nil
}
```

### 테스트 가능한 Mock 구현체
```go
type MockUserRepository struct {
    users map[int]*User
}

func (r *MockUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    if user, exists := r.users[id]; exists {
        return user, nil
    }
    return nil, fmt.Errorf("user not found")
}

// 테스트에서 사용
func TestUserService_CreateUser(t *testing.T) {
    // Given
    mockRepo := NewMockUserRepository()
    mockEmail := NewMockEmailService()
    mockCache := NewInMemoryCacheService()

    service := NewUserService(mockRepo, mockCache, mockEmail)

    // When
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")

    // Then
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### 의존성 그래프 자동 생성
```go
func (c *Container) Initialize() error {
    // 1단계: 기본 서비스
    c.cacheService = NewInMemoryCacheService()

    // 2단계: 외부 서비스
    c.emailService = c.createEmailService()

    // 3단계: Repository
    c.userRepository = c.createUserRepository()

    // 4단계: Business Service
    c.userService = NewUserService(
        c.userRepository,
        c.cacheService,
        c.emailService,
    )

    return nil
}
```

## 🎨 DI 패턴 비교

### Constructor Injection
```go
// ✅ 장점: 명확한 의존성, 불변성 보장
// ❌ 단점: 많은 의존성 시 복잡

func NewOrderService(
    orderRepo OrderRepository,
    productRepo ProductRepository,
    userRepo UserRepository,
    payment PaymentService,
    email EmailService,
    notification NotificationService,
) OrderService {
    return &OrderServiceImpl{...}
}
```

### Setter Injection
```go
// ✅ 장점: 유연성, 선택적 의존성
// ❌ 단점: 불완전한 초기화 가능성

type Service struct {
    repo Repository
}

func (s *Service) SetRepository(repo Repository) {
    s.repo = repo
}
```

### Interface Injection
```go
// ✅ 장점: 다형성, 테스트 용이
// ❌ 단점: 인터페이스 관리 복잡도

type RepositoryAware interface {
    SetRepository(repo Repository)
}
```

## 📝 베스트 프랙티스

### 1. **인터페이스는 사용하는 쪽에서 정의**
```go
// ❌ Bad: Repository 패키지에서 정의
package repository
type UserRepository interface {...}

// ✅ Good: Service 패키지에서 정의
package service
type UserRepository interface {
    // 필요한 메서드만 정의
    FindByID(id int) (*User, error)
}
```

### 2. **작은 인터페이스 선호**
```go
// ❌ Bad: 거대한 인터페이스
type UserRepository interface {
    FindByID(id int) (*User, error)
    FindByEmail(email string) (*User, error)
    FindByUsername(username string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
    List() ([]*User, error)
    Count() (int, error)
    // ... 20개 더
}

// ✅ Good: 분리된 인터페이스
type UserReader interface {
    FindByID(id int) (*User, error)
}

type UserWriter interface {
    Create(user *User) error
    Update(user *User) error
}
```

### 3. **의존성 순환 방지**
```go
// ❌ Bad: 순환 의존성
// Package A → Package B → Package A

// ✅ Good: 인터페이스를 통한 의존성 역전
// Package A → Interface ← Package B
```

### 4. **테스트를 위한 설계**
```go
// 모든 외부 의존성을 인터페이스로
type Service struct {
    db       Database    // 인터페이스
    cache    Cache       // 인터페이스
    logger   Logger      // 인터페이스
    metrics  Metrics     // 인터페이스
}

// 테스트에서 쉽게 Mock 주입
func TestService(t *testing.T) {
    service := NewService(
        NewMockDB(),
        NewMockCache(),
        NewMockLogger(),
        NewMockMetrics(),
    )
    // ...
}
```

## 🚀 고급 DI 도구

### Wire (Google)
```go
// +build wireinject

func InitializeApp() (*App, error) {
    wire.Build(
        NewDatabase,
        NewUserRepository,
        NewEmailService,
        NewUserService,
        NewApp,
    )
    return nil, nil
}
```

### Dig (Uber)
```go
container := dig.New()
container.Provide(NewDatabase)
container.Provide(NewUserRepository)
container.Provide(NewUserService)

err := container.Invoke(func(service UserService) {
    // service 사용
})
```

### Fx (Uber)
```go
app := fx.New(
    fx.Provide(
        NewDatabase,
        NewUserRepository,
        NewUserService,
    ),
    fx.Invoke(StartServer),
)
app.Run()
```

## 🎯 실전 체크리스트

- [ ] 모든 비즈니스 로직이 인터페이스에 의존하는가?
- [ ] 순환 의존성이 없는가?
- [ ] 각 컴포넌트가 단일 책임을 가지는가?
- [ ] 테스트에서 Mock을 쉽게 주입할 수 있는가?
- [ ] 환경별로 다른 구현체를 사용할 수 있는가?
- [ ] 의존성 그래프가 명확한가?
- [ ] 생성자가 너무 복잡하지 않은가?
- [ ] 인터페이스가 적절한 크기로 분리되어 있는가?

## 📚 추가 학습 자료

- [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)
- [Wire: Automated Dependency Injection](https://github.com/google/wire)
- [SOLID Principles in Go](https://dave.cheney.net/2016/08/20/solid-go-design)
- [Clean Architecture in Go](https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1)

## 🎯 다음 레슨 예고

**Lesson 14: 실행 모드 (Release/Debug/Test)**
- Gin의 실행 모드 설정
- 환경별 최적화
- 디버깅과 프로파일링
- 테스트 모드 활용

의존성 주입은 확장 가능하고 테스트 가능한 애플리케이션의 기초입니다! 💉