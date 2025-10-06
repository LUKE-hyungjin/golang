# 동적 HTML 페이지 만들기 🎨

안녕하세요! 이번에는 **템플릿**을 사용해서 동적으로 변하는 HTML 페이지를 만들어볼 거예요. 사용자 이름을 넣거나, 상품 목록을 보여주는 등 데이터에 따라 바뀌는 페이지를 만들 수 있어요!

## 템플릿이 뭔가요?

템플릿은 **빈칸이 있는 문서**라고 생각하면 돼요. 빈칸에 데이터를 넣어서 완성된 HTML을 만드는 거죠!

### 실생활 비유
- **편지 양식**: "___님께, 저는 ___입니다" → 빈칸에 이름을 채워넣기
- **상장**: "___님은 ___에서 1등을 했습니다" → 이름과 대회명 채우기
- **계약서**: 표준 양식에 이름, 날짜, 금액만 채우기

### 템플릿 vs 정적 HTML
```html
<!-- 정적 HTML: 항상 똑같음 -->
<h1>안녕하세요, 홍길동님!</h1>

<!-- 템플릿: 사용자마다 다름 -->
<h1>안녕하세요, {{.UserName}}님!</h1>
```

## 이번 챕터에서 배울 내용
- HTML 템플릿 파일 만들고 사용하기
- 템플릿에 데이터 전달하기 (이름, 목록 등)
- 조건문으로 다른 내용 보여주기 (if-else)
- 반복문으로 목록 만들기 (for)
- 로그인 폼과 데이터 처리하기
- 에러 페이지 예쁘게 만들기

## 📂 파일 구조
```
08/
├── main.go               # 메인 서버
├── templates/           # HTML 템플릿 파일
│   ├── index.html       # 메인 페이지
│   ├── users.html       # 사용자 목록
│   ├── user-detail.html # 사용자 상세
│   ├── products.html    # 제품 카탈로그
│   ├── contact.html     # 문의 폼
│   ├── contact-success.html # 문의 성공
│   ├── login.html       # 로그인
│   ├── dashboard.html   # 대시보드
│   ├── 404.html        # 404 에러
│   └── error.html      # 일반 에러
└── static/             # 정적 파일 (CSS, JS, 이미지)
```

## 핵심 개념 이해하기

### 1. 템플릿 파일 불러오기

서버가 시작될 때 템플릿 파일들을 읽어와야 해요!

```go
// 방법 1: templates 폴더의 모든 파일 불러오기
r.LoadHTMLGlob("templates/*")

// 방법 2: 특정 파일만 선택해서 불러오기
r.LoadHTMLFiles("templates/index.html", "templates/users.html")
```

**실생활 비유**: 학교에서 수업 시작 전에 교재를 준비해두는 것!

### 2. 나만의 함수 추가하기

템플릿에서 사용할 수 있는 특별한 함수를 만들 수 있어요!

```go
r.SetFuncMap(template.FuncMap{
    "formatDate": formatDate,           // 날짜 예쁘게 표시
    "formatCurrency": formatCurrency,   // 금액에 쉼표 넣기
})
```

**예시**: `{{formatCurrency 10000}}` → "10,000원"

**실생활 비유**: 계산기에 "10% 할인" 버튼을 추가하는 것!

### 3. 템플릿에 데이터 전달하고 HTML 만들기

데이터를 템플릿에 넣어서 최종 HTML을 만들어요!

```go
c.HTML(http.StatusOK, "index.html", gin.H{
    "title": "우리 쇼핑몰",
    "userName": "홍길동",
    "products": productList,
})
```

**어떻게 동작할까요?**
1. `index.html` 템플릿 파일을 찾아요
2. `{{.title}}`, `{{.userName}}` 같은 빈칸에 데이터를 채워요
3. 완성된 HTML을 브라우저로 보내요

**실생활 비유**: 주문서 양식에 고객 정보를 채워서 완성하는 것!

## 🚀 실행 방법

```bash
cd gin
go run ./08

# 브라우저에서 접속
http://localhost:8080
```

## 📋 페이지별 기능 테스트

### 1️⃣ 메인 페이지
```bash
# 브라우저 접속
http://localhost:8080/

# curl 테스트
curl http://localhost:8080/
```

### 2️⃣ 사용자 목록
```bash
# 사용자 목록 페이지
http://localhost:8080/users

# 개별 사용자 상세
http://localhost:8080/users/1
```

### 3️⃣ 제품 카탈로그
```bash
# 제품 목록
http://localhost:8080/products
```

### 4️⃣ 문의하기
```bash
# 문의 폼
http://localhost:8080/contact

# POST 테스트
curl -X POST http://localhost:8080/contact \
  -d "name=홍길동&email=hong@example.com&message=문의사항입니다"
```

### 5️⃣ 로그인/대시보드
```bash
# 로그인 페이지
http://localhost:8080/login

# 로그인 (admin/1234)
curl -X POST http://localhost:8080/login \
  -d "username=admin&password=1234" \
  -c cookies.txt \
  -L

# 대시보드 접근
curl http://localhost:8080/dashboard \
  -b cookies.txt

# 로그아웃
curl http://localhost:8080/logout \
  -b cookies.txt \
  -L
```

### 6️⃣ 에러 페이지
```bash
# 404 에러
http://localhost:8080/nonexistent

# 500 에러
http://localhost:8080/error
```

## 📝 템플릿 문법

### 1. 변수 출력
```html
{{.title}}                <!-- 변수 출력 -->
{{.user.Name}}           <!-- 중첩 필드 -->
```

### 2. 조건문
```html
{{if .IsActive}}
    <span>활성</span>
{{else}}
    <span>비활성</span>
{{end}}

{{if eq .Role "admin"}}
    <span>관리자</span>
{{else if eq .Role "user"}}
    <span>사용자</span>
{{end}}
```

### 3. 반복문
```html
{{range .users}}
    <li>{{.Name}} - {{.Email}}</li>
{{else}}
    <li>사용자가 없습니다</li>
{{end}}

{{range $index, $user := .users}}
    <tr class="{{if isEven $index}}even{{end}}">
        <td>{{$index}}</td>
        <td>{{$user.Name}}</td>
    </tr>
{{end}}
```

### 4. 커스텀 함수 사용
```html
{{formatDate .JoinedAt}}
{{formatCurrency .Price}}
{{add 1 2}}
```

### 5. 비교 연산자
```html
{{if gt .Age 18}}          <!-- 크다 -->
{{if lt .Count 10}}        <!-- 작다 -->
{{if ge .Score 60}}        <!-- 크거나 같다 -->
{{if le .Price 100}}       <!-- 작거나 같다 -->
{{if eq .Status "active"}} <!-- 같다 -->
{{if ne .Role "guest"}}    <!-- 같지 않다 -->
```

### 6. 논리 연산자
```html
{{if and .IsActive .IsVerified}}
{{if or .IsAdmin .IsModerator}}
{{if not .IsBlocked}}
```

## 🎨 구현된 페이지 설명

### 메인 페이지 (index.html)
- 네비게이션 메뉴
- 환영 메시지
- 기능 소개 카드

### 사용자 목록 (users.html)
- 통계 대시보드
- 테이블 형식의 사용자 목록
- 역할별 배지
- 상태 표시

### 사용자 상세 (user-detail.html)
- 프로필 정보
- 활동 통계
- 상세 정보 표시

### 제품 카탈로그 (products.html)
- 카테고리별 분류
- 그리드 레이아웃
- 가격 포맷팅
- 재고 상태

### 문의하기 (contact.html)
- 폼 입력
- 유효성 검사
- 제출 처리

### 대시보드 (dashboard.html)
- 로그인 상태 확인
- 통계 카드
- 최근 활동
- 로그아웃

## 🔍 트러블슈팅

### 템플릿 파일을 찾을 수 없음
```go
// 상대 경로 확인
r.LoadHTMLGlob("08/templates/*")

// 또는 절대 경로 사용
dir, _ := os.Getwd()
r.LoadHTMLGlob(filepath.Join(dir, "templates/*"))
```

### 커스텀 함수가 동작하지 않음
```go
// SetFuncMap을 LoadHTMLGlob보다 먼저 호출
r.SetFuncMap(template.FuncMap{...})
r.LoadHTMLGlob("templates/*")
```

### XSS 방지
```go
// html/template은 자동으로 이스케이프
{{.UserInput}}  // 안전

// HTML 그대로 출력 (위험)
{{.HTMLContent | safe}}  // template.HTML 타입 사용
```

### 템플릿 캐싱
```go
// 프로덕션 모드에서는 템플릿 캐싱
gin.SetMode(gin.ReleaseMode)
// 개발 모드에서는 매번 리로드
gin.SetMode(gin.DebugMode)
```

## 🏗️ 실전 활용 팁

### 1. 레이아웃 템플릿
```html
<!-- layout.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    {{template "content" .}}
</body>
</html>

<!-- page.html -->
{{define "content"}}
    <h1>Page Content</h1>
{{end}}
```

### 2. 파셜 템플릿
```html
<!-- header.html -->
{{define "header"}}
<header>...</header>
{{end}}

<!-- main.html -->
{{template "header" .}}
```

### 3. 데이터 전달 패턴
```go
type PageData struct {
    Title   string
    User    *User
    IsAuth  bool
    Data    interface{}
}

c.HTML(200, "page.html", PageData{
    Title: "페이지 제목",
    User:  currentUser,
    IsAuth: true,
    Data:  specificData,
})
```

### 4. 에러 처리 패턴
```go
func RenderError(c *gin.Context, code int, message string) {
    c.HTML(code, "error.html", gin.H{
        "code":    code,
        "message": message,
    })
}
```

### 5. CSRF 토큰
```go
func CSRFToken() string {
    // 토큰 생성 로직
    return token
}

r.SetFuncMap(template.FuncMap{
    "csrf": CSRFToken,
})

// 템플릿에서
// <input type="hidden" name="csrf" value="{{csrf}}">
```

## 📚 다음 학습 단계
- 에러 처리 & 로깅
- 구성 & 설정
- 데이터베이스 연동
- 보안 (CORS, 인증/인가)
- 테스트 작성

## 🔗 참고 자료
- [Gin HTML 렌더링 문서](https://gin-gonic.com/docs/examples/html-rendering/)
- [Go html/template 패키지](https://pkg.go.dev/html/template)
- [템플릿 보안 가이드](https://github.com/golang/go/wiki/WebAssemblySecurityModel)