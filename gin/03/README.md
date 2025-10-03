# 03. Path/Query/Body 파라미터 바인딩

## 📌 개요
Gin 프레임워크의 강력한 파라미터 바인딩 기능을 학습합니다. Path 파라미터, Query 파라미터, JSON/Form Body 데이터를 Go 구조체에 자동으로 바인딩하고 검증하는 방법을 다룹니다.

## 🎯 학습 목표
- Path 파라미터 추출 (`:id`, `*action`)
- Query 파라미터 파싱 및 기본값 설정
- JSON/Form 데이터 구조체 바인딩
- 데이터 검증 (validation tags)
- 파일 업로드 처리
- 복합 파라미터 처리 (Path + Query + Body)

## 📂 파일 구조
```
03/
└── main.go     # 파라미터 바인딩 예제
```

## 💻 코드 설명

### 주요 구성 요소

1. **구조체 태그**: JSON, Form, 검증 규칙 정의
2. **바인딩 메서드**: `ShouldBind`, `ShouldBindJSON`, `ShouldBindQuery`
3. **검증 태그**: `required`, `email`, `min`, `max` 등
4. **옵셔널 필드**: 포인터를 사용한 부분 업데이트

### 구조체 태그 설명
```go
type User struct {
    ID    string `json:"id" form:"id" binding:"required"`
    Name  string `json:"name" form:"name" binding:"required"`
    Email string `json:"email" form:"email" binding:"required,email"`
    Age   int    `json:"age" form:"age" binding:"min=1,max=120"`
}
```
- `json`: JSON 파싱 시 사용할 필드명
- `form`: Query/Form 파라미터 파싱 시 사용할 필드명
- `binding`: 검증 규칙 정의

## 🚀 실행 방법

### 1. 서버 시작
```bash
# gin 폴더에서 실행
cd gin
go run ./03

# 또는 03 폴더로 이동 후 실행
cd gin/03
go run main.go
```

### 2. API 테스트

## 📋 파라미터 타입별 예제

### 1️⃣ Path 파라미터

**단일 Path 파라미터:**
```bash
curl http://localhost:8080/users/123

# 응답:
# {
#   "message": "Path parameter example",
#   "user_id": "123"
# }
```

**다중 Path 파라미터:**
```bash
curl http://localhost:8080/users/456/posts/789

# 응답:
# {
#   "message": "Multiple path parameters",
#   "user_id": "456",
#   "post_id": "789"
# }
```

### 2️⃣ Query 파라미터

**기본 Query 파라미터:**
```bash
# 모든 파라미터 포함
curl "http://localhost:8080/search?q=golang&page=2&limit=20&sort=asc"

# 응답:
# {
#   "message": "Search results",
#   "query": "golang",
#   "page": "2",
#   "structured": {
#     "q": "golang",
#     "page": 2,
#     "limit": 20,
#     "sort": "asc"
#   }
# }
```

**기본값 적용:**
```bash
# 필수 파라미터만 전달
curl "http://localhost:8080/search?q=gin"

# 응답:
# {
#   "message": "Search results",
#   "query": "gin",
#   "page": "1",       # 기본값 적용
#   "structured": {
#     "q": "gin",
#     "page": 1,        # 기본값
#     "limit": 10,      # 기본값
#     "sort": "desc"    # 기본값
#   }
# }
```

### 3️⃣ JSON Body 파라미터

**유효한 데이터:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user001",
    "name": "홍길동",
    "email": "hong@example.com",
    "age": 30
  }'

# 응답:
# {
#   "message": "User created successfully",
#   "user": {
#     "id": "user001",
#     "name": "홍길동",
#     "email": "hong@example.com",
#     "age": 30
#   }
# }
```

**검증 실패 케이스:**
```bash
# 이메일 형식 오류
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user002",
    "name": "김철수",
    "email": "invalid-email",
    "age": 25
  }'

# 응답:
# {
#   "error": "Invalid JSON data",
#   "details": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"
# }
```

```bash
# 나이 범위 초과
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user003",
    "name": "이영희",
    "email": "lee@example.com",
    "age": 150
  }'

# 응답:
# {
#   "error": "Invalid JSON data",
#   "details": "Key: 'User.Age' Error:Field validation for 'Age' failed on the 'max' tag"
# }
```

### 4️⃣ Form 데이터

**HTML Form 데이터:**
```bash
curl -X POST http://localhost:8080/users/form \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=form001&name=박민수&email=park@example.com&age=35"

# 응답:
# {
#   "message": "Form data received",
#   "user": {
#     "id": "form001",
#     "name": "박민수",
#     "email": "park@example.com",
#     "age": 35
#   }
# }
```

### 5️⃣ 복합 파라미터 (Path + Query + Body)

```bash
curl -X PUT "http://localhost:8080/users/user001?notify=true" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user001",
    "name": "홍길동(수정)",
    "email": "hong.new@example.com",
    "age": 31
  }'

# 응답:
# {
#   "message": "User updated",
#   "id": "user001",           # Path 파라미터
#   "notify": "true",          # Query 파라미터
#   "updated_data": {          # Body 파라미터
#     "id": "user001",
#     "name": "홍길동(수정)",
#     "email": "hong.new@example.com",
#     "age": 31
#   }
# }
```

### 6️⃣ 부분 업데이트 (PATCH)

**일부 필드만 업데이트:**
```bash
# 이메일만 수정
curl -X PATCH http://localhost:8080/users/user001 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'

# 응답:
# {
#   "message": "User partially updated",
#   "updates": {
#     "id": "user001",
#     "email": "newemail@example.com"
#   }
# }
```

```bash
# 이름과 나이 수정
curl -X PATCH http://localhost:8080/users/user001 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "새이름",
    "age": 40
  }'

# 응답:
# {
#   "message": "User partially updated",
#   "updates": {
#     "id": "user001",
#     "name": "새이름",
#     "age": 40
#   }
# }
```

### 7️⃣ 파일 업로드

**단일 파일:**
```bash
# test.txt 파일 생성
echo "Hello, Gin!" > test.txt

# 파일 업로드
curl -X POST http://localhost:8080/upload \
  -F "file=@test.txt"

# 응답:
# {
#   "message": "File uploaded successfully",
#   "filename": "test.txt",
#   "size": 12
# }
```

**다중 파일:**
```bash
# 여러 파일 생성
echo "File 1" > file1.txt
echo "File 2" > file2.txt
echo "File 3" > file3.txt

# 다중 파일 업로드
curl -X POST http://localhost:8080/upload/multiple \
  -F "files=@file1.txt" \
  -F "files=@file2.txt" \
  -F "files=@file3.txt"

# 응답:
# {
#   "message": "Multiple files uploaded",
#   "count": 3,
#   "files": [
#     {"filename": "file1.txt", "size": 7},
#     {"filename": "file2.txt", "size": 7},
#     {"filename": "file3.txt", "size": 7}
#   ]
# }
```

## 📝 주요 학습 포인트

### 1. 바인딩 메서드 비교

| 메서드 | 용도 | 에러 시 동작 |
|--------|------|------------|
| `c.Bind()` | 자동 타입 감지 | 400 응답 자동 반환 |
| `c.ShouldBind()` | 자동 타입 감지 | 에러만 반환 |
| `c.ShouldBindJSON()` | JSON 전용 | 에러만 반환 |
| `c.ShouldBindQuery()` | Query 파라미터 전용 | 에러만 반환 |
| `c.ShouldBindUri()` | Path 파라미터 전용 | 에러만 반환 |

### 2. 검증 태그 종류
```go
// 필수 필드
`binding:"required"`

// 문자열 길이
`binding:"min=3,max=20"`

// 숫자 범위
`binding:"min=1,max=100"`

// 이메일 형식
`binding:"email"`

// URL 형식
`binding:"url"`

// 정규식
`binding:"alphanum"`  // 영숫자만

// 조건부 검증
`binding:"required_if=Role Admin"`

// 커스텀 검증
`binding:"customValidator"`
```

### 3. 옵셔널 필드 처리
```go
// 포인터 사용
type UpdateRequest struct {
    Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty"`
}

// nil 체크로 수정 여부 판단
if req.Name != nil {
    user.Name = *req.Name
}
```

### 4. Content-Type별 바인딩
```go
// application/json
c.ShouldBindJSON(&data)

// application/x-www-form-urlencoded
c.ShouldBind(&data)

// multipart/form-data
c.ShouldBind(&data)

// application/xml
c.ShouldBindXML(&data)

// application/yaml
c.ShouldBindYAML(&data)
```

## 🔍 트러블슈팅

### JSON 바인딩 실패
```bash
# Content-Type 헤더 확인
-H "Content-Type: application/json"

# JSON 구조 확인
# 올바른 구조: {"key": "value"}
# 잘못된 구조: {'key': 'value'}  # 작은따옴표 사용
```

### 검증 에러 메시지 파싱
```go
// 상세한 검증 에러 처리
if err := c.ShouldBindJSON(&user); err != nil {
    // 검증 에러를 파싱하여 필드별 에러 메시지 생성
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        for _, fe := range ve {
            // 필드명과 태그 정보 추출
            field := fe.Field()
            tag := fe.Tag()
        }
    }
}
```

### 파일 업로드 크기 제한
```go
// main() 함수에서 설정
r.MaxMultipartMemory = 8 << 20  // 8 MiB (기본값)
```

## 🏗️ 실전 활용 팁

### 1. 커스텀 검증 함수
```go
// 전화번호 형식 검증
func ValidatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^010-\d{4}-\d{4}$`, phone)
    return matched
}

// 등록
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("phone", ValidatePhone)
}
```

### 2. 에러 응답 표준화
```go
type ErrorResponse struct {
    Error   string            `json:"error"`
    Details map[string]string `json:"details,omitempty"`
}
```

### 3. 페이지네이션 구조체
```go
type Pagination struct {
    Page  int `form:"page,default=1" binding:"min=1"`
    Limit int `form:"limit,default=20" binding:"min=1,max=100"`
    Sort  string `form:"sort,default=created_at"`
    Order string `form:"order,default=desc" binding:"oneof=asc desc"`
}
```

## 📚 다음 단계
- 컨텍스트 활용: Request/Response 처리
- 미들웨어 구현: 인증, 로깅, CORS
- 에러 핸들링: 전역 에러 처리기

## 🔗 참고 자료
- [Gin 모델 바인딩과 검증](https://gin-gonic.com/docs/examples/binding-and-validation/)
- [Go Validator 문서](https://github.com/go-playground/validator)
- [HTTP Form 처리 가이드](https://developer.mozilla.org/ko/docs/Learn/Forms)