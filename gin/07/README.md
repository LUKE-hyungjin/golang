# 이미지, CSS, JS 파일 서빙하기 📦

안녕하세요! 웹사이트를 만들 때는 HTML뿐만 아니라 이미지, CSS, JavaScript 같은 **정적 파일**도 필요해요. 이번 챕터에서는 Gin으로 이런 파일들을 어떻게 제공하는지 배워봅시다!

## 정적 파일이 뭔가요?

정적 파일은 **변하지 않는 파일**을 말해요. 사용자가 요청할 때마다 똑같은 내용을 보여주는 파일들이에요.

### 정적 파일의 예시
- **이미지**: logo.png, profile.jpg
- **스타일시트**: style.css (웹사이트 꾸미기)
- **JavaScript**: app.js (인터랙션 추가)
- **문서**: manual.pdf, terms.txt

### 실생활 비유
- **도서관의 책**: 누가 빌려가든 같은 내용
- **박물관의 전시물**: 누가 봐도 똑같은 작품
- **자판기의 음료**: 누가 뽑아도 같은 음료

## 이번 챕터에서 배울 내용
- 정적 파일(이미지, CSS, JS) 제공하기
- 파일 업로드 받고 저장하기
- 파일 다운로드 기능 만들기
- SPA(React, Vue 같은 프론트) 지원하기
- 캐시로 속도 빠르게 만들기

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

## 핵심 개념 이해하기

### 1. Static() - 폴더 전체 공개하기

폴더 안의 모든 파일을 한 번에 제공할 수 있어요!

```go
r.Static("/static", "./static")
// URL: /static/css/style.css → 실제 파일: ./static/css/style.css
```

**어떻게 동작할까요?**
- 브라우저에서 `/static/logo.png`를 요청하면
- 서버의 `./static/logo.png` 파일을 찾아서 보내줍니다

**실생활 비유**: 도서관 전체를 개방하는 것과 같아요. 누구나 들어와서 책을 찾아볼 수 있죠!

### 2. StaticFile() - 특정 파일 하나만 제공하기

딱 한 개의 파일만 제공하고 싶을 때 사용해요!

```go
r.StaticFile("/favicon.ico", "./static/favicon.ico")
```

**언제 사용할까요?**
- 웹사이트 아이콘(favicon)
- robots.txt (검색엔진 안내)
- 홈페이지 대표 이미지

**실생활 비유**: 도서관의 베스트셀러 코너에 특정 책 한 권만 진열하는 것!

### 3. StaticFS() - 고급 옵션으로 파일 제공

더 세밀하게 제어하고 싶을 때 사용해요!

```go
r.StaticFS("/assets", http.Dir("./static"))
```

**Static()과 뭐가 다른가요?**
- 파일 시스템을 직접 지정할 수 있어요
- 압축된 파일이나 메모리의 파일도 제공 가능해요
- 고급 사용자용!

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