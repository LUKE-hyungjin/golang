# 04. Context 사용 (c.Param, c.Query, c.Bind, c.JSON)

## 📌 개요
Gin의 핵심인 Context 객체를 완벽하게 이해하고 활용하는 방법을 학습합니다. Context는 요청 처리 중 모든 정보를 담고 있으며, 요청 데이터 추출부터 응답 생성까지 모든 작업의 중심입니다.

## 🎯 학습 목표
- gin.Context의 주요 메서드 완벽 이해
- 요청 데이터 추출 방법 (Param, Query, Header, Body)
- 다양한 응답 포맷 생성 (JSON, XML, YAML, String)
- Context를 통한 값 저장 및 전달
- 쿠키, 파일 업로드/다운로드 처리
- 스트림 응답 및 비동기 처리

## 📂 파일 구조
```
04/
└── main.go     # Context 활용 예제
```

## 💻 Context 주요 메서드 카테고리

### 1. 요청 데이터 추출
- `c.Param()`: Path 파라미터
- `c.Query()`, `c.DefaultQuery()`: Query 파라미터
- `c.GetHeader()`: 헤더 값
- `c.ShouldBind()`, `c.ShouldBindJSON()`: Body 바인딩
- `c.Cookie()`: 쿠키 값

### 2. 응답 생성
- `c.JSON()`: JSON 응답
- `c.XML()`: XML 응답
- `c.YAML()`: YAML 응답
- `c.String()`: 문자열 응답
- `c.Data()`: 바이너리 응답
- `c.HTML()`: HTML 템플릿 응답

### 3. Context 데이터 관리
- `c.Set()`: 값 저장
- `c.Get()`, `c.MustGet()`: 값 가져오기
- `c.Copy()`: Context 복사 (고루틴용)

### 4. 플로우 제어
- `c.Next()`: 다음 핸들러 실행
- `c.Abort()`: 체인 중단
- `c.AbortWithStatusJSON()`: 중단 + 에러 응답
- `c.Redirect()`: 리다이렉트

## 🚀 실행 방법

### 서버 시작
```bash
cd gin
go run ./04
```

## 📋 API 테스트 예제

### 1️⃣ 기본 Context 메서드 테스트
```bash
# Path, Query 파라미터와 헤더 정보
curl "http://localhost:8080/context/basic/john?age=25&city=Seoul&page=2" \
  -H "User-Agent: MyApp/1.0"

# 응답:
# {
#   "name": "john",
#   "age": "25",
#   "city": "Seoul",
#   "page": "2",
#   "user_agent": "MyApp/1.0",
#   "content_type": ""
# }
```

### 2️⃣ Request Body 바인딩
```bash
# 유효한 데이터
curl -X POST http://localhost:8080/context/bind \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user001",
    "username": "johndoe",
    "email": "john@example.com",
    "age": 30,
    "role": "admin"
  }'

# 검증 실패 (잘못된 role)
curl -X POST http://localhost:8080/context/bind \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user002",
    "username": "jane",
    "email": "jane@example.com",
    "age": 25,
    "role": "superadmin"
  }'
```

### 3️⃣ 다양한 응답 포맷
```bash
# JSON 응답
curl http://localhost:8080/context/response/json

# XML 응답
curl http://localhost:8080/context/response/xml

# YAML 응답
curl http://localhost:8080/context/response/yaml

# 문자열 응답
curl http://localhost:8080/context/response/string

# 바이너리 데이터 응답
curl http://localhost:8080/context/response/data
```

### 4️⃣ 로그인과 Context 값 저장
```bash
# 로그인 (Context에 값 저장)
curl -X POST http://localhost:8080/context/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'

# 응답:
# {
#   "message": "Login successful",
#   "username": "admin",
#   "role": "admin",
#   "authenticated": true
# }

# 실패 케이스
curl -X POST http://localhost:8080/context/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "wrongpass"
  }'
```

### 5️⃣ 파일 업로드/다운로드
```bash
# 파일 생성
echo "Test file content" > test.txt

# 파일 업로드
curl -X POST http://localhost:8080/context/upload \
  -F "file=@test.txt"

# 파일 다운로드
curl -O -J http://localhost:8080/context/download
```

### 6️⃣ 쿠키 처리
```bash
# 쿠키 설정
curl -c cookies.txt http://localhost:8080/context/cookie/set

# 쿠키 읽기
curl -b cookies.txt http://localhost:8080/context/cookie/get

# 응답:
# {"session_id":"abc123xyz"}
```

### 7️⃣ Request 정보 조회
```bash
curl http://localhost:8080/context/request-info \
  -H "Referer: https://google.com" \
  -H "Custom-Header: CustomValue"
```

### 8️⃣ 스트림 응답 (Server-Sent Events)
```bash
# SSE 스트림 받기
curl -N http://localhost:8080/context/stream

# 출력:
# event:message
# data:0
#
# event:message
# data:1
#
# ... (1초 간격으로 계속)
```

### 9️⃣ 비동기 처리
```bash
curl http://localhost:8080/context/async

# 즉시 응답:
# {"message":"Request is being processed asynchronously"}
#
# 서버 로그에 2초 후:
# Async request from: 127.0.0.1
```

### 🔟 요청 중단 (Abort)
```bash
# Authorization 헤더 없이 (중단됨)
curl http://localhost:8080/context/abort

# 응답:
# {"error":"Authorization header required"}

# Authorization 헤더와 함께
curl http://localhost:8080/context/abort \
  -H "Authorization: Bearer token123"

# 응답:
# {"message":"Authorized","token":"Bearer token123"}
```

### 1️⃣1️⃣ 리다이렉트
```bash
# 외부 리다이렉트
curl -L http://localhost:8080/context/redirect

# 내부 리다이렉트
curl http://localhost:8080/old-endpoint

# 응답:
# {"message":"This is the new endpoint"}
```

### 1️⃣2️⃣ Content Negotiation
```bash
# JSON 요청
curl http://localhost:8080/context/negotiate \
  -H "Accept: application/json"

# XML 요청
curl http://localhost:8080/context/negotiate \
  -H "Accept: application/xml"

# YAML 요청
curl http://localhost:8080/context/negotiate \
  -H "Accept: application/x-yaml"
```

## 📝 핵심 포인트

### 1. Context 생명주기
```go
// Context는 요청별로 생성되고 소멸
// 하나의 요청 처리 동안 모든 미들웨어와 핸들러에서 공유
Request → Middleware1 → Middleware2 → Handler → Response
        ↓            ↓            ↓
    c.Set("key")  c.Get("key")  c.Get("key")
```

### 2. 안전한 타입 변환
```go
// Get with type assertion
value, exists := c.Get("key")
if exists {
    strValue, ok := value.(string)
    if ok {
        // 사용
    }
}

// MustGet (없으면 panic)
strValue := c.MustGet("key").(string)
```

### 3. 고루틴에서 Context 사용
```go
// 절대 원본 Context를 고루틴에 전달하지 마세요!
go func(c *gin.Context) {
    // 위험: 요청 완료 후 Context가 재사용될 수 있음
}(c)

// 올바른 방법: Copy() 사용
cCp := c.Copy()
go func() {
    // 안전: 복사된 Context 사용
    log.Println(cCp.Request.URL.Path)
}()
```

### 4. 에러 처리 패턴
```go
// 패턴 1: Early return
if err := c.ShouldBindJSON(&data); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}

// 패턴 2: Abort
if !authorized {
    c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
    return // return 필수!
}
```

## 🔍 트러블슈팅

### ShouldBind vs Bind
```go
// Bind: 자동으로 400 응답 (커스텀 에러 메시지 불가)
c.Bind(&data)

// ShouldBind: 에러만 반환 (커스텀 처리 가능)
if err := c.ShouldBind(&data); err != nil {
    // 커스텀 에러 응답
}
```

### Context 값이 없을 때
```go
// 안전한 처리
if value, exists := c.Get("key"); exists {
    // 사용
} else {
    // 기본값 처리
}
```

### 파일 업로드 크기 제한
```go
// main()에서 설정
r.MaxMultipartMemory = 8 << 20  // 8 MiB
```

## 🏗️ 실전 활용 팁

### 1. Request ID 추가
```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("RequestID", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

### 2. 사용자 정보 저장
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 토큰 검증 후
        c.Set("UserID", userID)
        c.Set("UserRole", role)
        c.Next()
    }
}
```

### 3. 응답 래퍼
```go
func SuccessResponse(c *gin.Context, data interface{}) {
    c.JSON(200, gin.H{
        "success": true,
        "data":    data,
        "timestamp": time.Now().Unix(),
    })
}

func ErrorResponse(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "success": false,
        "error":   message,
        "timestamp": time.Now().Unix(),
    })
}
```

## 📚 다음 단계
- [05. 미들웨어](../05/README.md): 전역/그룹/개별 미들웨어
- [06. 라우트 그룹](../06/README.md): API 버저닝과 그룹화
- [07. 정적 파일](../07/README.md): Static 파일 서빙

## 🔗 참고 자료
- [Gin Context 공식 문서](https://pkg.go.dev/github.com/gin-gonic/gin#Context)
- [HTTP 헤더 MDN](https://developer.mozilla.org/ko/docs/Web/HTTP/Headers)
- [Server-Sent Events](https://developer.mozilla.org/ko/docs/Web/API/Server-sent_events)