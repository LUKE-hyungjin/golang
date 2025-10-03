# Lesson 18: CORS (Cross-Origin Resource Sharing) 설정 🌐

> 웹 애플리케이션의 크로스 오리진 요청을 안전하게 처리하는 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

CORS는 웹 브라우저가 다른 도메인의 리소스에 접근할 수 있도록 허용하는 메커니즘입니다. SPA(Single Page Application)와 API 서버가 분리된 현대 웹 아키텍처에서 필수적인 설정입니다. 이번 레슨에서는 Gin에서 CORS를 구성하는 다양한 방법을 학습합니다.

### 핵심 학습 목표
- ✅ CORS 동작 원리 이해
- ✅ 환경별 CORS 설정
- ✅ Preflight 요청 처리
- ✅ 동적 CORS 구성
- ✅ 보안 헤더와 함께 사용
- ✅ 라우트별 CORS 설정

## 🏗 CORS 아키텍처

### CORS 요청 흐름
```
1. Preflight Request (OPTIONS)
   Browser → Server
   ├── Origin: http://localhost:3000
   ├── Access-Control-Request-Method: POST
   └── Access-Control-Request-Headers: Content-Type

2. Preflight Response
   Server → Browser
   ├── Access-Control-Allow-Origin: http://localhost:3000
   ├── Access-Control-Allow-Methods: POST
   ├── Access-Control-Allow-Headers: Content-Type
   └── Access-Control-Max-Age: 86400

3. Actual Request
   Browser → Server (with actual data)

4. Actual Response
   Server → Browser (with CORS headers)
```

## 🛠 구현된 기능

### 1. **환경별 CORS 설정**
- Development: 모든 localhost 허용
- Production: 특정 도메인만 허용
- Default: 기본 보안 설정

### 2. **커스텀 CORS 미들웨어**
- 세밀한 제어 가능
- 와일드카드 서브도메인 지원
- 동적 오리진 검증

### 3. **gin-contrib/cors 패키지**
- 표준화된 CORS 구현
- 간편한 설정
- 검증된 보안

### 4. **라우트별 CORS**
- 엔드포인트별 다른 설정
- Public/Private API 구분
- Admin 전용 CORS

## 🎯 주요 API 엔드포인트

### 공개 API
```bash
GET  /public/info        # 모든 오리진 허용
GET  /public/health      # 헬스체크
```

### 일반 API
```bash
GET  /api/data          # 설정된 오리진만 허용
POST /api/echo          # Echo 서비스
```

### 비공개 API
```bash
GET  /private/user      # 특정 오리진만 허용
POST /private/upload    # 파일 업로드
```

### 관리자 API
```bash
GET  /admin/users       # 관리자 도메인만 허용
GET  /admin/settings    # 보안 설정
```

### 테스트 엔드포인트
```bash
GET    /test            # GET 메서드 테스트
POST   /test            # POST 메서드 테스트
PUT    /test            # PUT 메서드 테스트
DELETE /test            # DELETE 메서드 테스트
PATCH  /test            # PATCH 메서드 테스트
```

## 💻 실습 가이드

### 1. 실행
```bash
cd gin/18
go mod init cors-example
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/cors

# 개발 모드
APP_ENV=development go run main.go

# 프로덕션 모드
APP_ENV=production go run main.go
```

### 2. CORS 테스트

#### 브라우저 콘솔에서 테스트
```javascript
// 개발자 도구 콘솔에서 실행
fetch('http://localhost:8080/api/data', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```

#### cURL로 테스트
```bash
# 일반 요청
curl -H "Origin: http://localhost:3000" \
     http://localhost:8080/api/data

# Preflight 요청
curl -X OPTIONS \
     -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -v http://localhost:8080/api/data

# 응답 헤더 확인
curl -I -H "Origin: http://localhost:3000" \
     http://localhost:8080/api/data
```

### 3. 다양한 오리진 테스트

#### 허용된 오리진
```bash
# localhost:3000 (허용됨)
curl -H "Origin: http://localhost:3000" \
     http://localhost:8080/api/data

# 응답
HTTP/1.1 200 OK
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Credentials: true
```

#### 허용되지 않은 오리진
```bash
# example.org (차단됨)
curl -H "Origin: http://example.org" \
     http://localhost:8080/api/data

# 프로덕션 모드에서는 403 Forbidden
```

### 4. 크로스 오리진 POST 요청

```bash
# Preflight 요청 (자동)
curl -X OPTIONS \
     -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     http://localhost:8080/api/echo

# 실제 POST 요청
curl -X POST \
     -H "Origin: http://localhost:3000" \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello CORS"}' \
     http://localhost:8080/api/echo
```

### 5. 파일 업로드 with CORS

```bash
# HTML 파일 생성 (test.html)
cat > test.html << 'EOF'
<!DOCTYPE html>
<html>
<body>
  <input type="file" id="fileInput">
  <button onclick="uploadFile()">Upload</button>
  <script>
    function uploadFile() {
      const file = document.getElementById('fileInput').files[0];
      const formData = new FormData();
      formData.append('file', file);

      fetch('http://localhost:8080/private/upload', {
        method: 'POST',
        body: formData
      })
      .then(response => response.json())
      .then(data => console.log(data))
      .catch(error => console.error('Error:', error));
    }
  </script>
</body>
</html>
EOF

# 브라우저에서 test.html 열고 테스트
```

### 6. 설정 확인

```bash
# 현재 CORS 설정 조회
curl http://localhost:8080/cors/config | jq

# 응답 예시
{
  "environment": "development",
  "allowed_origins": [
    "http://localhost:3000",
    "http://localhost:3001",
    "http://localhost:8080",
    "http://127.0.0.1:3000",
    "http://127.0.0.1:8080"
  ],
  "allowed_methods": ["*"],
  "allowed_headers": ["*"],
  "expose_headers": [
    "Content-Length",
    "Content-Type",
    "X-Request-ID",
    "X-RateLimit-Limit"
  ],
  "allow_credentials": true,
  "max_age": "24h0m0s"
}
```

## 🔍 코드 하이라이트

### 커스텀 CORS 미들웨어
```go
func CustomCORS(config CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // 오리진 검증
        if isOriginAllowed(origin, config.AllowOrigins) {
            c.Header("Access-Control-Allow-Origin", origin)

            // Preflight 요청 처리
            if c.Request.Method == "OPTIONS" {
                c.Header("Access-Control-Allow-Methods",
                    strings.Join(config.AllowMethods, ", "))
                c.Header("Access-Control-Allow-Headers",
                    strings.Join(config.AllowHeaders, ", "))
                c.Header("Access-Control-Max-Age",
                    fmt.Sprintf("%d", int(config.MaxAge.Seconds())))
                c.AbortWithStatus(http.StatusNoContent)
                return
            }
        }

        c.Next()
    }
}
```

### 와일드카드 서브도메인 지원
```go
func isOriginAllowed(origin string, allowedOrigins []string) bool {
    for _, allowed := range allowedOrigins {
        // 와일드카드 처리
        if strings.HasPrefix(allowed, "*.") {
            domain := strings.TrimPrefix(allowed, "*")
            if strings.HasSuffix(origin, domain) {
                return true
            }
        }
        if allowed == "*" || allowed == origin {
            return true
        }
    }
    return false
}
```

### 동적 CORS 설정
```go
func DynamicCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // 데이터베이스나 캐시에서 확인
        if isDynamicallyAllowed(origin) {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Credentials", "true")
        }

        c.Next()
    }
}
```

### gin-contrib/cors 사용
```go
import "github.com/gin-contrib/cors"

router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
        return origin == "https://github.com"
    },
    MaxAge: 12 * time.Hour,
}))
```

## 🎨 CORS 설정 패턴

### 1. **개발 환경**
```go
// 모든 localhost 허용
AllowOrigins: []string{
    "http://localhost:*",
    "http://127.0.0.1:*",
}
```

### 2. **프로덕션 환경**
```go
// 명시적 도메인만 허용
AllowOrigins: []string{
    "https://app.example.com",
    "https://www.example.com",
}
```

### 3. **마이크로서비스**
```go
// 내부 서비스 간 통신
AllowOrigins: []string{
    "https://*.internal.example.com",
}
```

### 4. **파트너 API**
```go
// 동적 파트너 도메인 관리
func isPartnerDomain(origin string) bool {
    // 데이터베이스에서 파트너 목록 조회
    partners := getPartnersFromDB()
    return partners[origin]
}
```

## 📝 베스트 프랙티스

### 1. **최소 권한 원칙**
```go
// ❌ Bad: 모든 것 허용
AllowOrigins: []string{"*"}
AllowHeaders: []string{"*"}
AllowMethods: []string{"*"}

// ✅ Good: 필요한 것만 허용
AllowOrigins: []string{"https://app.example.com"}
AllowHeaders: []string{"Content-Type", "Authorization"}
AllowMethods: []string{"GET", "POST"}
```

### 2. **Credentials 주의**
```go
// Credentials를 허용할 때는 와일드카드 사용 불가
if config.AllowCredentials {
    // ❌ Bad
    c.Header("Access-Control-Allow-Origin", "*")

    // ✅ Good
    c.Header("Access-Control-Allow-Origin", specificOrigin)
}
```

### 3. **Preflight 캐싱**
```go
// 브라우저가 Preflight 결과를 캐시하도록 설정
c.Header("Access-Control-Max-Age", "86400") // 24시간
```

### 4. **보안 헤더 함께 사용**
```go
// CORS와 함께 보안 헤더 설정
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "SAMEORIGIN")
c.Header("X-XSS-Protection", "1; mode=block")
```

## 🚀 프로덕션 체크리스트

- [ ] 허용 오리진 목록 최소화
- [ ] 와일드카드 사용 제한
- [ ] Credentials 필요성 검토
- [ ] Preflight 캐시 시간 설정
- [ ] 보안 헤더 추가
- [ ] 로깅 및 모니터링
- [ ] Rate limiting 적용
- [ ] HTTPS 사용

## 🔒 보안 고려사항

### CORS 우회 공격 방지
- 오리진 검증 철저히
- Referer 헤더 추가 검증
- 세션/토큰 기반 인증 병행

### 민감한 데이터 보호
- Expose Headers 최소화
- Credentials 신중히 사용
- 암호화된 연결 (HTTPS)

## 📚 추가 학습 자료

- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [CORS in Action](https://livebook.manning.com/book/cors-in-action)
- [gin-contrib/cors](https://github.com/gin-contrib/cors)
- [OWASP CORS Security](https://owasp.org/www-community/attacks/CORS_OriginHeaderScrutiny)

## 🎯 다음 레슨 예고

**Lesson 19: JWT 인증 미들웨어**
- JWT 토큰 생성과 검증
- Refresh Token 구현
- 역할 기반 접근 제어
- 보안 모범 사례

CORS를 제대로 설정하여 안전한 크로스 오리진 통신을 구현하세요! 🌐