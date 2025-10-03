# 06. 라우트 그룹, 버저닝 (v1, v2)

## 📌 개요
Gin의 라우트 그룹 기능을 활용하여 API를 체계적으로 구성하고, 버전별로 관리하는 방법을 학습합니다. 대규모 애플리케이션에서 필수적인 API 버저닝과 라우트 조직화 패턴을 다룹니다.

## 🎯 학습 목표
- 라우트 그룹을 사용한 API 구조화
- URL 경로 기반 API 버저닝 (v1, v2)
- 헤더 기반 API 버저닝
- 중첩 라우트 그룹 구성
- 그룹별 미들웨어 적용
- 관리자, 내부 API, 공개 API 분리

## 📂 파일 구조
```
06/
└── main.go     # 라우트 그룹과 버저닝 예제
```

## 💻 라우트 그룹 구조

### 구현된 그룹 구조
```
/
├── /api
│   ├── /v1                    # API 버전 1
│   │   ├── /users
│   │   │   ├── /:id
│   │   │   ├── /:id/profile
│   │   │   └── /:id/settings
│   │   └── /products
│   │       └── /:id
│   └── /v2                    # API 버전 2 (개선된 버전)
│       ├── /users
│       │   ├── /:id
│       │   ├── /:id/activities
│       │   └── /:id/follow
│       └── /products
│           ├── /search
│           ├── /:id
│           └── /:id/reviews
├── /admin                     # 관리자 패널
│   ├── /dashboard
│   ├── /users
│   └── /system
├── /public                    # 공개 API
│   ├── /status
│   └── /docs
├── /internal                  # 내부 서비스용
└── /webhooks                  # Webhook 엔드포인트
```

## 🚀 실행 방법

```bash
cd gin
go run ./06
```

## 📋 API 테스트 예제

### 1️⃣ API 버전 1 (v1)

**사용자 목록 조회:**
```bash
curl http://localhost:8080/api/v1/users

# 응답:
# {
#   "version": "v1",
#   "users": [
#     {
#       "id": "1",
#       "username": "user1",
#       "email": "user1@example.com",
#       "created_at": "2024-01-01T10:00:00Z",
#       "api_version": "v1"
#     }
#   ]
# }
```

**특정 사용자 조회:**
```bash
curl http://localhost:8080/api/v1/users/123

# 사용자 프로필 조회
curl http://localhost:8080/api/v1/users/123/profile

# 사용자 설정 조회
curl http://localhost:8080/api/v1/users/123/settings
```

**제품 조회:**
```bash
# 제품 목록
curl http://localhost:8080/api/v1/products

# 특정 제품
curl http://localhost:8080/api/v1/products/456
```

### 2️⃣ API 버전 2 (v2) - 개선된 버전

**v2의 개선사항:**
- 페이지네이션 지원
- 더 상세한 응답 형식
- 필터링 기능
- 추가 엔드포인트

```bash
# 페이지네이션이 추가된 사용자 목록
curl "http://localhost:8080/api/v2/users?page=1&limit=20"

# 응답:
# {
#   "version": "v2",
#   "data": [...],
#   "pagination": {
#     "page": "1",
#     "limit": "20",
#     "total": 100,
#     "total_pages": 5
#   }
# }
```

**v2 전용 기능:**
```bash
# 사용자 활동 내역 (v2에서 추가)
curl http://localhost:8080/api/v2/users/123/activities

# 사용자 팔로우 (v2에서 추가)
curl -X POST http://localhost:8080/api/v2/users/123/follow

# 제품 검색 (v2에서 추가)
curl "http://localhost:8080/api/v2/products/search?q=laptop"

# 제품 리뷰 (v2에서 추가)
curl http://localhost:8080/api/v2/products/456/reviews

# 필터링된 제품 목록
curl "http://localhost:8080/api/v2/products?category=Electronics&min_price=100&max_price=1000"
```

### 3️⃣ 관리자 패널 (인증 필요)

**인증 없이 접근 (실패):**
```bash
curl http://localhost:8080/admin/dashboard

# 응답:
# {"error":"Admin authentication required"}
```

**관리자 토큰으로 접근:**
```bash
# 대시보드
curl http://localhost:8080/admin/dashboard \
  -H "X-Admin-Token: admin-secret-token"

# 모든 사용자 조회
curl http://localhost:8080/admin/users \
  -H "X-Admin-Token: admin-secret-token"

# 사용자 차단
curl -X PUT http://localhost:8080/admin/users/123/ban \
  -H "X-Admin-Token: admin-secret-token"

# 시스템 로그
curl http://localhost:8080/admin/system/logs \
  -H "X-Admin-Token: admin-secret-token"

# 시스템 메트릭
curl http://localhost:8080/admin/system/metrics \
  -H "X-Admin-Token: admin-secret-token"
```

### 4️⃣ 공개 API (인증 불필요)

```bash
# 서비스 상태
curl http://localhost:8080/public/status

# API 문서
curl http://localhost:8080/public/docs
```

### 5️⃣ 내부 API (내부 서비스용)

```bash
# 상세 헬스체크
curl http://localhost:8080/internal/health/detailed \
  -H "X-Internal-API-Key: internal-api-key-123"

# 캐시 클리어
curl -X POST http://localhost:8080/internal/cache/clear \
  -H "X-Internal-API-Key: internal-api-key-123"

# 작업 트리거
curl -X POST "http://localhost:8080/internal/jobs/trigger?type=backup" \
  -H "X-Internal-API-Key: internal-api-key-123"
```

### 6️⃣ Webhook 엔드포인트

```bash
# GitHub webhook
curl -X POST http://localhost:8080/webhooks/github \
  -H "X-GitHub-Event: push" \
  -d '{"ref":"refs/heads/main"}'

# Stripe webhook
curl -X POST http://localhost:8080/webhooks/stripe \
  -d '{"type":"payment.succeeded"}'

# Slack webhook
curl -X POST http://localhost:8080/webhooks/slack \
  -d '{"text":"Hello from Slack"}'
```

### 7️⃣ 헤더 기반 버저닝

```bash
# 헤더로 v1 지정
curl http://localhost:8080/api/users \
  -H "API-Version: 1.0"

# 헤더로 v2 지정
curl http://localhost:8080/api/users \
  -H "API-Version: 2.0"

# 헤더 없이 (기본값: 최신 버전)
curl http://localhost:8080/api/users
```

## 📝 핵심 포인트

### 1. 라우트 그룹 생성

```go
// 기본 그룹 생성
v1 := r.Group("/api/v1")

// 미들웨어와 함께
v2 := r.Group("/api/v2", middleware())

// 또는
v2 := r.Group("/api/v2")
v2.Use(middleware())
```

### 2. 중첩 그룹

```go
api := r.Group("/api")
{
    v1 := api.Group("/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", getUsers)
            users.POST("", createUser)

            profile := users.Group("/:id/profile")
            {
                profile.GET("", getProfile)
                profile.PUT("", updateProfile)
            }
        }
    }
}
```

### 3. 버저닝 전략

**URL 경로 방식 (추천):**
```
/api/v1/users
/api/v2/users
```

**헤더 방식:**
```
API-Version: 1.0
Accept: application/vnd.api+json;version=1
```

**Query 파라미터 방식:**
```
/api/users?version=1
```

### 4. 그룹별 미들웨어 적용

```go
// 공개 API
public := r.Group("/public")

// 인증이 필요한 API
protected := r.Group("/api")
protected.Use(AuthMiddleware())

// 관리자 API
admin := r.Group("/admin")
admin.Use(AuthMiddleware(), RequireRole("admin"))
```

## 🔍 트러블슈팅

### 라우트 충돌

```go
// 충돌 발생 가능
r.GET("/users/:id", handler1)
r.GET("/users/me", handler2)  // :id가 "me"를 잡아버림

// 해결: 순서 변경
r.GET("/users/me", handler2)  // 구체적인 경로를 먼저
r.GET("/users/:id", handler1)
```

### 버전 마이그레이션

```go
// v1을 v2로 리다이렉트
v1.GET("/old-endpoint", func(c *gin.Context) {
    c.Redirect(301, "/api/v2/new-endpoint")
})
```

### 버전별 응답 형식

```go
// v1 응답
type ResponseV1 struct {
    Data interface{} `json:"data"`
}

// v2 응답 (개선된 형식)
type ResponseV2 struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Meta    interface{} `json:"meta,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}
```

## 🏗️ 실전 활용 팁

### 1. API 버전 관리 모범 사례

```go
// 버전별 핸들러 분리
package v1
func GetUsers(c *gin.Context) { }

package v2
func GetUsers(c *gin.Context) { }

// 라우터 설정
v1Group.GET("/users", v1.GetUsers)
v2Group.GET("/users", v2.GetUsers)
```

### 2. 버전 지원 종료 알림

```go
func DeprecationMiddleware(version string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-API-Deprecated", "true")
        c.Header("X-API-Sunset-Date", "2025-01-01")
        c.Next()
    }
}
```

### 3. 기능 플래그

```go
func FeatureFlag(feature string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !IsFeatureEnabled(feature) {
            c.AbortWithStatusJSON(404, gin.H{
                "error": "Feature not available",
            })
            return
        }
        c.Next()
    }
}
```

### 4. 라우트 문서화

```go
// Swagger 주석 추가
// @Summary 사용자 목록 조회
// @Tags Users
// @Version 2.0
// @Router /api/v2/users [get]
func GetUsersV2(c *gin.Context) { }
```

## 📚 다음 단계
- [07. 정적 파일 서빙](../07/README.md)
- [08. 템플릿 렌더링](../08/README.md)

## 🔗 참고 자료
- [Gin 라우트 그룹 문서](https://gin-gonic.com/docs/examples/grouping-routes/)
- [REST API 버저닝 가이드](https://www.baeldung.com/rest-versioning)
- [API 설계 모범 사례](https://swagger.io/resources/articles/best-practices-in-api-design/)