# API를 깔끔하게 정리하는 라우트 그룹 📁

안녕하세요! 프로젝트가 커지면 API가 수십, 수백 개가 됩니다. 이때 **라우트 그룹**을 사용하면 API를 폴더처럼 깔끔하게 정리할 수 있어요. 또한 API 버전(v1, v2)을 관리하는 방법도 배워봅시다!

## 라우트 그룹이 뭔가요?

라우트 그룹은 **비슷한 API들을 묶어서 관리하는 기능**입니다. 마치 컴퓨터에서 파일을 폴더별로 정리하는 것처럼, API도 용도별로 그룹을 만들어 관리할 수 있어요!

### 실생활 비유
- **백화점 층별 안내**: 1층은 화장품, 2층은 의류, 3층은 식품처럼 구역을 나누는 것
- **학교 학년별 반**: 1학년 1반, 1학년 2반, 2학년 1반... 이렇게 체계적으로 분류
- **도서관 서가**: 소설, 역사, 과학 등 카테고리별로 책을 배치

## 이번 챕터에서 배울 내용
- 라우트 그룹으로 API 깔끔하게 정리하기
- API 버전 관리하기 (v1, v2로 구분)
- 관리자 API, 공개 API 분리하기
- 그룹마다 다른 미들웨어 적용하기
- 중첩 그룹으로 복잡한 구조 만들기

## 📂 파일 구조
```
06/
└── main.go     # 라우트 그룹과 버저닝 예제
```

## 핵심 개념 이해하기

### 우리 프로젝트의 API 구조
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

## 💡 꼭 알아야 할 핵심 개념!

### 1. 라우트 그룹 만들기

그룹을 만드는 방법은 정말 간단해요!

```go
// 방법 1: 기본 그룹 만들기
v1 := r.Group("/api/v1")

// 방법 2: 그룹 만들면서 미들웨어 함께 적용
v2 := r.Group("/api/v2", 로깅미들웨어())

// 방법 3: 나중에 미들웨어 추가
v2 := r.Group("/api/v2")
v2.Use(인증미들웨어())
```

**실생활 비유**: 아파트 단지를 만들고(Group), 각 동마다 CCTV를 설치하는(미들웨어) 것과 같아요!

### 2. 중첩 그룹 - 그룹 안에 그룹 만들기

복잡한 구조도 깔끔하게 만들 수 있어요!

```go
api := r.Group("/api")           // 1층: API 전체
{
    v1 := api.Group("/v1")       // 2층: v1 버전
    {
        users := v1.Group("/users")  // 3층: 사용자 관련
        {
            users.GET("", 사용자목록)
            users.POST("", 사용자생성)

            profile := users.Group("/:id/profile")  // 4층: 프로필
            {
                profile.GET("", 프로필조회)
                profile.PUT("", 프로필수정)
            }
        }
    }
}
```

**실생활 비유**: 백화점에서 1층 → 화장품 코너 → 스킨케어 섹션 → 특정 브랜드처럼 단계별로 들어가는 것!

### 3. API 버전 관리 - 어떤 방법이 좋을까요?

**방법 1: URL에 버전 넣기 (가장 많이 사용!) ⭐**
```
/api/v1/users  ← 구버전
/api/v2/users  ← 신버전
```
- **장점**: 주소만 봐도 버전을 알 수 있어요
- **단점**: URL이 조금 길어져요

**방법 2: 헤더로 버전 지정**
```
API-Version: 1.0
API-Version: 2.0
```
- **장점**: URL이 깔끔해요
- **단점**: 헤더를 봐야 버전을 알 수 있어요

**방법 3: Query 파라미터 사용**
```
/api/users?version=1
```
- **장점**: 유연하게 버전을 바꿀 수 있어요
- **단점**: 캐싱이 복잡해져요

💡 **추천**: 대부분의 회사에서는 **방법 1(URL 경로)**을 사용합니다!

### 4. 그룹마다 다른 보안 적용하기

```go
// 공개 API - 누구나 접근 가능
public := r.Group("/public")
public.GET("/status", 서버상태)

// 일반 API - 로그인 필요
api := r.Group("/api")
api.Use(로그인체크())

// 관리자 API - 관리자만 접근
admin := r.Group("/admin")
admin.Use(로그인체크(), 관리자체크())
```

**실생활 비유**:
- **공개 API**: 공원처럼 누구나 들어갈 수 있는 곳
- **일반 API**: 아파트처럼 주민만 들어갈 수 있는 곳
- **관리자 API**: 관리사무소처럼 직원만 들어갈 수 있는 곳

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