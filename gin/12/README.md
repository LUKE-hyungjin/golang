# Lesson 12: 환경변수와 설정 파일 로딩 (Viper) 🔧

> Go 애플리케이션의 설정 관리를 체계적으로 다루는 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

애플리케이션이 커질수록 설정 관리는 매우 중요해집니다. 하드코딩된 값들을 환경변수나 설정 파일로 관리하면 환경별 배포가 쉬워지고, 보안도 강화됩니다. 이번 레슨에서는 Go의 대표적인 설정 관리 라이브러리인 **Viper**를 사용하여 전문적인 설정 시스템을 구축합니다.

### 핵심 학습 목표
- ✅ Viper를 사용한 설정 파일 로딩
- ✅ 환경변수와 설정 파일의 우선순위 관리
- ✅ 환경별 설정 파일 분리 (development, production)
- ✅ 설정값 검증과 기본값 처리
- ✅ 실시간 설정 변경 감지 (Hot Reload)
- ✅ 구조화된 설정 타입 시스템

## 🛠 구현된 기능

### 1. **설정 구조 설계**
```go
type Config struct {
    Server   ServerConfig   // 서버 설정
    Database DatabaseConfig // 데이터베이스 설정
    Redis    RedisConfig    // Redis 캐시 설정
    JWT      JWTConfig      // 인증 토큰 설정
    Email    EmailConfig    // 이메일 서비스 설정
    Storage  StorageConfig  // 파일 스토리지 설정
    Logging  LoggingConfig  // 로깅 설정
    Security SecurityConfig // 보안 설정
    Features FeatureFlags   // 기능 플래그
    External ExternalAPIs   // 외부 API 설정
}
```

### 2. **설정 로더 인터페이스**
```go
type ConfigLoader interface {
    Load() (*Config, error)           // 설정 로딩
    Watch(callback func(*Config))      // 변경 감지
    Get(key string) interface{}       // 키로 값 조회
    Set(key string, value interface{}) // 런타임 설정 변경
}
```

### 3. **우선순위 시스템**
1. 환경변수 (최우선)
2. 환경별 설정 파일 (config.production.yaml)
3. 기본 설정 파일 (config.yaml)
4. 하드코딩된 기본값

### 4. **환경별 설정 파일**
- `config.yaml` - 기본 설정
- `config.development.yaml` - 개발 환경
- `config.production.yaml` - 프로덕션 환경
- `.env.example` - 환경변수 템플릿

## 🎯 주요 API 엔드포인트

### 설정 관리 API
```bash
GET    /config              # 현재 설정 조회 (민감정보 마스킹)
GET    /config/:key         # 특정 설정값 조회
POST   /config/:key         # 설정값 업데이트
POST   /config/reload       # 설정 파일 재로드
GET    /config/environment  # 현재 환경 정보
GET    /features           # 기능 플래그 조회
POST   /features/:flag     # 기능 플래그 토글
```

### 헬스체크 API
```bash
GET    /health             # 애플리케이션 상태
GET    /health/db          # 데이터베이스 연결 상태
GET    /health/redis       # Redis 연결 상태
```

## 💻 실습 가이드

### 1. 기본 실행
```bash
# 기본 설정으로 실행
go run main.go

# 개발 환경으로 실행
APP_ENV=development go run main.go

# 프로덕션 환경으로 실행
APP_ENV=production go run main.go
```

### 2. 환경변수 사용
```bash
# .env 파일 복사
cp .env.example .env

# 환경변수로 포트 변경
APP_PORT=3000 go run main.go

# 여러 환경변수 설정
DB_HOST=localhost DB_PORT=5432 DB_USERNAME=myuser go run main.go
```

### 3. API 테스트

#### 설정 조회
```bash
# 전체 설정 조회 (민감정보 마스킹됨)
curl http://localhost:8080/config | jq

# 특정 설정값 조회
curl http://localhost:8080/config/server.port
curl http://localhost:8080/config/features.beta_features

# 환경 정보 조회
curl http://localhost:8080/environment
```

#### 설정 업데이트
```bash
# 로그 레벨 변경
curl -X POST http://localhost:8080/config/logging.level \
  -H "Content-Type: application/json" \
  -d '{"value": "debug"}'

# Rate Limit 설정 변경
curl -X POST http://localhost:8080/config/security.rate_limit.requests_per_minute \
  -H "Content-Type: application/json" \
  -d '{"value": 120}'
```

#### 기능 플래그
```bash
# 기능 플래그 조회
curl http://localhost:8080/features

# 베타 기능 활성화
curl -X POST http://localhost:8080/features/beta_features \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'

# 유지보수 모드 전환
curl -X POST http://localhost:8080/features/maintenance_mode \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

#### 설정 재로드
```bash
# 설정 파일 수정 후 재로드
curl -X POST http://localhost:8080/config/reload
```

## 🔍 코드 하이라이트

### Viper 초기화
```go
func NewViperConfigLoader(configPath, environment string) *ViperConfigLoader {
    v := viper.New()

    // 설정 파일 경로
    v.AddConfigPath(configPath)
    v.SetConfigName("config")
    v.SetConfigType("yaml")

    // 환경변수 바인딩
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 기본값 설정
    setDefaults(v)

    return &ViperConfigLoader{
        viper:       v,
        configPath:  configPath,
        environment: environment,
    }
}
```

### 계층적 설정 로딩
```go
func (v *ViperConfigLoader) Load() (*Config, error) {
    // 1. 기본 설정 로드
    if err := v.viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    // 2. 환경별 설정 오버라이드
    envConfig := fmt.Sprintf("config.%s", v.environment)
    v.viper.SetConfigName(envConfig)
    if err := v.viper.MergeInConfig(); err == nil {
        log.Printf("Loaded environment config: %s", envConfig)
    }

    // 3. 환경변수 오버라이드 (자동)

    // 4. 구조체로 언마샬
    var config Config
    if err := v.viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    // 5. 검증
    if err := v.validate(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}
```

### 실시간 설정 변경 감지
```go
func (v *ViperConfigLoader) Watch(callback func(*Config)) {
    v.viper.WatchConfig()
    v.viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)

        var newConfig Config
        if err := v.viper.Unmarshal(&newConfig); err == nil {
            if err := v.validate(&newConfig); err == nil {
                v.mu.Lock()
                v.config = &newConfig
                v.mu.Unlock()

                callback(&newConfig)
            }
        }
    })
}
```

### 민감정보 마스킹
```go
func maskSensitiveConfig(config Config) Config {
    masked := config // 복사

    // 데이터베이스 비밀번호 마스킹
    if len(masked.Database.Password) > 0 {
        masked.Database.Password = "***"
    }

    // JWT 시크릿 마스킹
    if len(masked.JWT.Secret) > 0 {
        masked.JWT.Secret = "***"
    }

    // API 키 마스킹
    if len(masked.External.PaymentGateway.APIKey) > 0 {
        masked.External.PaymentGateway.APIKey = "sk_***"
    }

    return masked
}
```

## 🎨 설정 파일 구조

### config.yaml (기본 설정)
```yaml
server:
  host: 0.0.0.0
  port: 8080
  mode: debug

database:
  driver: postgres
  host: localhost
  port: 5432

features:
  new_dashboard: false
  beta_features: false
  maintenance_mode: false
```

### config.production.yaml (프로덕션 오버라이드)
```yaml
server:
  mode: release

database:
  host: ${DB_HOST}        # 환경변수에서 읽기
  password: ${DB_PASSWORD}
  ssl_mode: require

logging:
  level: info
  format: json

security:
  ssl_redirect: true
  rate_limit:
    enabled: true
```

## 📝 베스트 프랙티스

### 1. **환경변수 네이밍 규칙**
```bash
# 계층 구조를 언더스코어로 표현
DATABASE_HOST=localhost      # database.host
DATABASE_PORT=5432           # database.port
SECURITY_CORS_ENABLED=true   # security.cors.enabled
```

### 2. **설정 검증**
```go
func validate(config *Config) error {
    // 필수 값 체크
    if config.Server.Port <= 0 {
        return errors.New("server port must be positive")
    }

    // 범위 체크
    if config.Security.RateLimit.RequestsPerMinute < 1 {
        return errors.New("rate limit must be at least 1")
    }

    // 형식 체크
    if !isValidEmail(config.Email.From) {
        return errors.New("invalid email address")
    }

    return nil
}
```

### 3. **설정 분리 원칙**
- 🔒 **비밀 정보**: 환경변수로만 관리 (절대 커밋하지 않음)
- 🌍 **환경별 설정**: 환경별 파일로 분리
- 📦 **기본값**: 코드나 기본 설정 파일에 포함
- 🚀 **기능 플래그**: 독립적으로 관리

### 4. **12-Factor App 원칙**
```go
// 설정을 코드에서 분리
// 환경변수를 통한 설정 주입
// 환경 간 이식성 보장
```

## 🚀 프로덕션 체크리스트

- [ ] 모든 비밀 정보가 환경변수로 관리되는가?
- [ ] .env 파일이 .gitignore에 포함되어 있는가?
- [ ] 프로덕션 설정 파일이 준비되어 있는가?
- [ ] 설정 검증 로직이 충분한가?
- [ ] 기본값이 안전한 값으로 설정되어 있는가?
- [ ] 설정 변경 시 재시작이 필요한 항목이 문서화되어 있는가?
- [ ] 환경별 설정 차이가 명확히 구분되는가?
- [ ] 설정 관련 로그가 민감정보를 노출하지 않는가?

## 📚 추가 학습 자료

- [Viper 공식 문서](https://github.com/spf13/viper)
- [12-Factor App - Config](https://12factor.net/config)
- [환경변수 vs 설정 파일](https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d)
- [Go 설정 관리 패턴](https://medium.com/@felipedutratine/manage-config-in-golang-to-get-environment-variables-c2f6d4b8b93e)

## 🎯 다음 레슨 예고

**Lesson 13: 의존성 주입 (Dependency Injection)**
- 인터페이스 기반 설계
- 팩토리 패턴
- Wire를 사용한 DI
- 테스트 가능한 구조 설계

설정 관리는 확장 가능한 애플리케이션의 핵심입니다! 🔧