# 07. 정적 파일 서빙 (Static, StaticFS)

## 📌 개요
Gin에서 정적 파일(CSS, JavaScript, 이미지, 문서 등)을 서빙하는 다양한 방법을 학습합니다. SPA(Single Page Application) 지원, 파일 업로드/다운로드, 캐시 제어 등 실전에서 필요한 기능들을 구현합니다.

## 🎯 학습 목표
- Static(), StaticFile(), StaticFS() 메서드 활용
- 파일 업로드와 다운로드 구현
- SPA 라우팅 처리
- 캐시 제어 헤더 설정
- 파일 관리 API 구축
- 정적 파일과 API 엔드포인트 조합

## 📂 파일 구조
```
07/
├── main.go              # 메인 서버
├── static/             # 정적 파일 디렉토리
│   ├── index.html      # 메인 HTML
│   ├── css/
│   │   └── style.css   # 스타일시트
│   ├── js/
│   │   └── app.js      # JavaScript
│   ├── images/         # 이미지 파일
│   └── robots.txt      # robots.txt
└── uploads/            # 업로드 파일 저장소
```

## 💻 정적 파일 서빙 방법

### 1. Static() - 디렉토리 전체 서빙
```go
r.Static("/static", "./static")
// URL: /static/css/style.css → 파일: ./static/css/style.css
```

### 2. StaticFile() - 단일 파일 서빙
```go
r.StaticFile("/favicon.ico", "./static/favicon.ico")
```

### 3. StaticFS() - FileSystem 인터페이스 사용
```go
r.StaticFS("/assets", http.Dir("./static"))
```

## 🚀 실행 방법

```bash
cd gin
go run ./07

# 브라우저에서 접속
# http://localhost:8080
```

## 📋 기능 테스트

### 1️⃣ 정적 파일 접근

**CSS 파일:**
```bash
curl http://localhost:8080/static/css/style.css
```

**JavaScript 파일:**
```bash
curl http://localhost:8080/static/js/app.js
```

**robots.txt:**
```bash
curl http://localhost:8080/static/robots.txt
```

### 2️⃣ 파일 업로드

```bash
# 텍스트 파일 생성
echo "Hello, Gin!" > test.txt

# 파일 업로드
curl -X POST http://localhost:8080/upload \
  -F "file=@test.txt"

# 응답:
# {
#   "message": "File uploaded successfully",
#   "filename": "test.txt",
#   "size": 12,
#   "url": "/uploads/test.txt"
# }

# 이미지 파일 업로드
curl -X POST http://localhost:8080/upload \
  -F "file=@image.jpg"
```

### 3️⃣ 업로드된 파일 접근

```bash
# 업로드된 파일 직접 접근
curl http://localhost:8080/uploads/test.txt

# 다운로드 (Content-Disposition 헤더 포함)
curl -O -J http://localhost:8080/download/test.txt
```

### 4️⃣ 파일 관리 API

**파일 목록 조회:**
```bash
curl http://localhost:8080/api/files

# 응답:
# {
#   "files": [
#     {
#       "name": "test.txt",
#       "size": 12,
#       "url": "/uploads/test.txt"
#     }
#   ],
#   "total": 1
# }
```

**파일 정보 조회:**
```bash
curl http://localhost:8080/api/files/test.txt/info

# 응답:
# {
#   "name": "test.txt",
#   "size": 12,
#   "modified": "2024-01-01T10:00:00Z",
#   "is_directory": false
# }
```

**파일 삭제:**
```bash
curl -X DELETE http://localhost:8080/api/files/test.txt

# 응답:
# {
#   "message": "File deleted successfully",
#   "filename": "test.txt"
# }
```

### 5️⃣ 캐시 제어

```bash
# 캐시 헤더가 포함된 응답
curl -I http://localhost:8080/cached/css/style.css

# 헤더:
# Cache-Control: public, max-age=3600
# ETag: W/"123456"
```

### 6️⃣ SPA 지원

```bash
# 존재하지 않는 경로도 index.html 반환 (SPA 라우팅)
curl http://localhost:8080/about
curl http://localhost:8080/users/123
# 모두 index.html 내용 반환

# API 경로는 404 반환
curl http://localhost:8080/api/nonexistent

# 응답:
# {"error":"API endpoint not found"}
```

## 🌐 웹 UI 사용

브라우저에서 http://localhost:8080 접속 시:

1. **파일 업로드**: 파일 선택 후 업로드 버튼 클릭
2. **파일 목록**: 업로드된 파일 목록 자동 표시
3. **다운로드**: 각 파일의 다운로드 버튼 클릭
4. **삭제**: 각 파일의 삭제 버튼 클릭
5. **드래그 앤 드롭**: 파일을 브라우저로 드래그하여 업로드

## 📝 핵심 포인트

### 1. 정적 파일 서빙 옵션

```go
// 방법 1: Static (가장 일반적)
r.Static("/static", "./static")

// 방법 2: StaticFile (특정 파일)
r.StaticFile("/favicon.ico", "./favicon.ico")

// 방법 3: StaticFS (커스텀 FileSystem)
r.StaticFS("/fs", http.Dir("./public"))

// 방법 4: gin.H와 c.File() 조합
r.GET("/file/:name", func(c *gin.Context) {
    c.File("./files/" + c.Param("name"))
})
```

### 2. 파일 업로드 처리

```go
// 단일 파일
file, _ := c.FormFile("file")
c.SaveUploadedFile(file, dst)

// 다중 파일
form, _ := c.MultipartForm()
files := form.File["files"]
for _, file := range files {
    c.SaveUploadedFile(file, dst)
}
```

### 3. 다운로드 헤더 설정

```go
c.Header("Content-Disposition", "attachment; filename=file.txt")
c.Header("Content-Type", "application/octet-stream")
c.File(filepath)
```

### 4. SPA 라우팅 처리

```go
r.NoRoute(func(c *gin.Context) {
    // API 경로는 404
    if strings.HasPrefix(c.Request.URL.Path, "/api/") {
        c.JSON(404, gin.H{"error": "Not found"})
        return
    }
    // 나머지는 index.html
    c.File("./static/index.html")
})
```

## 🔍 트러블슈팅

### 파일 업로드 크기 제한

```go
// 기본값: 32 MB
r.MaxMultipartMemory = 8 << 20  // 8 MB로 설정
```

### MIME 타입 설정

```go
c.Header("Content-Type", "image/png")  // 수동 설정
// 또는
c.File()  // 자동 감지
```

### 경로 순회 공격 방지

```go
// 위험: 경로 순회 가능
filepath := "./uploads/" + c.Param("filename")

// 안전: 경로 정규화
filename := filepath.Base(c.Param("filename"))
filepath := filepath.Join("./uploads", filename)
```

### CORS 설정 (정적 파일)

```go
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Next()
})
```

## 🏗️ 실전 활용 팁

### 1. 압축 미들웨어

```go
import "github.com/gin-contrib/gzip"

r.Use(gzip.Gzip(gzip.DefaultCompression))
```

### 2. 조건부 캐싱

```go
func CacheMiddleware(maxAge int) gin.HandlerFunc {
    return func(c *gin.Context) {
        if strings.HasPrefix(c.Request.URL.Path, "/static/") {
            c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
        }
        c.Next()
    }
}
```

### 3. 파일 타입 검증

```go
func ValidateFileType(file *multipart.FileHeader) bool {
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }

    // MIME 타입 체크
    buffer := make([]byte, 512)
    f, _ := file.Open()
    f.Read(buffer)
    contentType := http.DetectContentType(buffer)

    return allowedTypes[contentType]
}
```

### 4. 보안 헤더

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Next()
    }
}
```

## 📚 다음 단계
- [08. 템플릿 렌더링](../08/README.md)

## 🔗 참고 자료
- [Gin Static Files 문서](https://gin-gonic.com/docs/examples/serving-static-files/)
- [MDN Web Docs - HTTP 캐싱](https://developer.mozilla.org/ko/docs/Web/HTTP/Caching)
- [SPA 라우팅 가이드](https://blog.pshrmn.com/entry/single-page-applications-and-the-server/)