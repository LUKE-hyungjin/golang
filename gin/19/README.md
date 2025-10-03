# Lesson 19: JWT 인증 미들웨어 🔐

> JSON Web Token을 사용한 안전한 인증 시스템 구현 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

JWT(JSON Web Token)는 현대 웹 애플리케이션에서 가장 널리 사용되는 인증 방식입니다. Stateless하고 확장 가능한 인증 시스템을 구현할 수 있습니다. 이번 레슨에서는 Gin에서 JWT를 사용한 완전한 인증 시스템을 구축합니다.

### 핵심 학습 목표
- ✅ JWT 토큰 생성과 검증
- ✅ Access Token과 Refresh Token
- ✅ 역할 기반 접근 제어 (RBAC)
- ✅ 토큰 갱신 메커니즘
- ✅ 보안 모범 사례
- ✅ 미들웨어를 통한 인증

## 🏗 JWT 아키텍처

### 토큰 구조
```
Header.Payload.Signature

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIn0.
TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

### 인증 플로우
```
1. Login Request → Server
2. Server → Validate Credentials
3. Server → Generate Access + Refresh Token
4. Client → Store Tokens
5. Client → API Request with Access Token
6. Server → Validate Token → Process Request
7. Access Token Expired → Use Refresh Token
8. Server → Generate New Access Token
```

## 🛠 구현된 기능

### 1. **토큰 관리**
- Access Token (15분 유효)
- Refresh Token (7일 유효)
- Token ID 추적
- 토큰 취소 (Revoke)

### 2. **인증 엔드포인트**
- 회원가입
- 로그인
- 토큰 갱신
- 로그아웃

### 3. **미들웨어**
- 인증 미들웨어
- 역할 기반 미들웨어
- 선택적 인증 미들웨어

### 4. **보안 기능**
- 비밀번호 해싱 (bcrypt)
- 서명 검증
- 토큰 만료 체크
- Issuer/Audience 검증

## 🎯 주요 API 엔드포인트

### 인증 엔드포인트
```bash
POST /api/v1/register    # 회원가입
POST /api/v1/login       # 로그인
POST /api/v1/refresh     # 토큰 갱신
POST /api/v1/logout      # 로그아웃
```

### 보호된 엔드포인트
```bash
GET  /api/v1/profile     # 프로필 조회 (인증 필요)
GET  /api/v1/protected   # 보호된 리소스
```

### 관리자 엔드포인트
```bash
GET  /api/v1/admin/users     # 사용자 목록 (admin 역할 필요)
GET  /api/v1/admin/dashboard # 관리자 대시보드
```

### 공개 엔드포인트
```bash
GET  /api/v1/public      # 공개 엔드포인트 (선택적 인증)
GET  /health             # 헬스체크
GET  /jwt/info           # JWT 설정 정보
```

## 💻 실습 가이드

### 1. 설치 및 실행
```bash
cd gin/19
go mod init jwt-example
go get -u github.com/gin-gonic/gin
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt

# 실행
go run main.go

# 또는 시크릿 키 설정
JWT_SECRET=your-secret-key go run main.go
```

### 2. 회원가입

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "password": "password123"
  }'

# 응답
{
  "message": "User registered successfully",
  "user": {
    "id": 3,
    "email": "newuser@example.com",
    "username": "newuser",
    "role": "user"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "expires_at": "2024-01-15T11:45:00Z"
  }
}
```

### 3. 로그인

```bash
# 일반 사용자 로그인
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "user123"
  }'

# 관리자 로그인
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# 토큰을 변수에 저장
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.tokens.access_token')

echo $TOKEN
```

### 4. 보호된 엔드포인트 접근

```bash
# 인증 없이 접근 (실패)
curl http://localhost:8080/api/v1/profile

# 응답: {"error":"Authorization header required"}

# 토큰으로 접근 (성공)
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"

# 응답
{
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "username": "admin",
    "role": "admin"
  },
  "token_info": {
    "issued_at": "2024-01-15T11:30:00Z",
    "expires_at": "2024-01-15T11:45:00Z",
    "issuer": "gin-jwt-example"
  }
}
```

### 5. 역할 기반 접근 제어

```bash
# 일반 사용자 토큰으로 관리자 엔드포인트 접근 (실패)
USER_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"user123"}' \
  | jq -r '.tokens.access_token')

curl http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $USER_TOKEN"

# 응답: {"error":"Insufficient permissions"}

# 관리자 토큰으로 접근 (성공)
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.tokens.access_token')

curl http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 응답: {"message":"Admin access granted","data":"Secret admin data"}
```

### 6. 토큰 갱신

```bash
# Refresh Token 저장
REFRESH_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"user123"}' \
  | jq -r '.tokens.refresh_token')

# 새 Access Token 요청
curl -X POST http://localhost:8080/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"

# 응답
{
  "message": "Token refreshed successfully",
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "expires_at": "2024-01-15T12:00:00Z"
  }
}
```

### 7. 선택적 인증

```bash
# 인증 없이 접근
curl http://localhost:8080/api/v1/public

# 응답
{
  "message": "This is a public endpoint",
  "authenticated": false
}

# 인증과 함께 접근
curl http://localhost:8080/api/v1/public \
  -H "Authorization: Bearer $TOKEN"

# 응답
{
  "message": "This is a public endpoint",
  "authenticated": true,
  "user": "admin"
}
```

### 8. 로그아웃

```bash
curl -X POST http://localhost:8080/api/v1/logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"

# 응답: {"message":"Logged out successfully"}
```

## 🔍 코드 하이라이트

### JWT Claims 구조
```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

### 토큰 생성
```go
func generateAccessToken(user *User) (string, time.Time, error) {
    expiresAt := time.Now().Add(15 * time.Minute)

    claims := Claims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "gin-jwt-example",
            Subject:   fmt.Sprintf("%d", user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}
```

### 인증 미들웨어
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authorization 헤더 확인
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Bearer 토큰 파싱
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        // 토큰 검증
        claims, err := ValidateToken(bearerToken[1])
        if err != nil {
            c.JSON(401, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        // Context에 저장
        c.Set("claims", claims)
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### 역할 기반 미들웨어
```go
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, _ := c.Get("claims")
        userClaims := claims.(*Claims)

        for _, role := range requiredRoles {
            if userClaims.Role == role {
                c.Next()
                return
            }
        }

        c.JSON(403, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    }
}
```

## 🎨 JWT 보안 모범 사례

### 1. **적절한 토큰 유효 기간**
```go
// Access Token: 짧게 (5-15분)
AccessTokenExpiry: 15 * time.Minute

// Refresh Token: 길게 (7-30일)
RefreshTokenExpiry: 7 * 24 * time.Hour
```

### 2. **강력한 시크릿 키**
```go
// ❌ Bad: 약한 키
SecretKey: "secret"

// ✅ Good: 강력한 키
SecretKey: os.Getenv("JWT_SECRET") // 최소 32자 이상
```

### 3. **토큰 저장 위치**
```javascript
// ❌ Bad: localStorage (XSS 취약)
localStorage.setItem('token', token)

// ✅ Good: HttpOnly Cookie
document.cookie = `token=${token}; HttpOnly; Secure; SameSite=Strict`

// ✅ Good: Memory + Refresh in HttpOnly Cookie
```

### 4. **토큰 검증**
```go
// 서명 검증
// 만료 시간 체크
// Issuer 확인
// Audience 확인
// 토큰 ID 추적
```

## 📝 베스트 프랙티스

### 1. **토큰 블랙리스트**
```go
// Redis를 사용한 토큰 무효화
type TokenBlacklist interface {
    Add(tokenID string, expiry time.Duration)
    Exists(tokenID string) bool
}
```

### 2. **Rate Limiting**
```go
// 로그인 시도 제한
loginAttempts := make(map[string]int)
if loginAttempts[email] > 5 {
    return errors.New("too many login attempts")
}
```

### 3. **토큰 로테이션**
```go
// Refresh 시 새로운 Refresh Token 발급
oldRefreshToken.Revoke()
newTokenPair := GenerateTokenPair(user)
```

### 4. **다중 디바이스 관리**
```go
// 사용자별 활성 세션 추적
type Session struct {
    UserID    uint
    TokenID   string
    Device    string
    CreatedAt time.Time
}
```

## 🚀 프로덕션 체크리스트

- [ ] 강력한 시크릿 키 사용
- [ ] HTTPS 전용 전송
- [ ] 토큰 유효 기간 최적화
- [ ] Refresh Token 안전한 저장
- [ ] 토큰 무효화 메커니즘
- [ ] Rate Limiting 적용
- [ ] 로그인 이상 징후 모니터링
- [ ] 정기적인 키 로테이션

## 🔒 보안 고려사항

### 공격 방어
- **XSS**: HttpOnly Cookie 사용
- **CSRF**: CSRF 토큰 추가
- **Replay Attack**: 토큰 ID와 nonce 사용
- **Brute Force**: Rate Limiting

### 추가 보안 레이어
- 2FA (Two-Factor Authentication)
- IP 화이트리스트
- 디바이스 핑거프린팅
- 이상 행동 탐지

## 📚 추가 학습 자료

- [JWT.io](https://jwt.io/)
- [RFC 7519 - JSON Web Token](https://datatracker.ietf.org/doc/html/rfc7519)
- [OWASP JWT Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)
- [golang-jwt Documentation](https://github.com/golang-jwt/jwt)

## 🎯 다음 레슨 예고

**Lesson 20: 입력 검증 (Binding + Validator)**
- 구조체 태그 검증
- 커스텀 검증 규칙
- 에러 메시지 커스터마이징
- 다국어 검증 메시지

JWT로 안전하고 확장 가능한 인증 시스템을 구축하세요! 🔐