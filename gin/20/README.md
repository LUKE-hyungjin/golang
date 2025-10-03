# Lesson 20: 입력 검증 (Input Validation) 🛡️

> Gin의 강력한 바인딩과 검증 시스템으로 안전한 API 구축하기

## 📌 이번 레슨에서 배우는 내용

입력 검증은 웹 애플리케이션 보안의 첫 번째 방어선입니다. Gin은 go-playground/validator를 기반으로 한 강력한 검증 시스템을 제공합니다. 이번 레슨에서는 기본 검증부터 커스텀 검증, 다국어 에러 메시지까지 모든 것을 다룹니다.

### 핵심 학습 목표
- ✅ 구조체 태그 기반 검증
- ✅ 커스텀 검증 규칙 구현
- ✅ 다국어 에러 메시지
- ✅ 중첩 구조체 검증
- ✅ 조건부 검증
- ✅ 크로스 필드 검증

## 🏗 검증 아키텍처

### 검증 플로우
```
1. HTTP Request → Gin Router
2. ShouldBind* → Parse & Bind
3. Validator → Apply Rules
4. Error Translation → Format Messages
5. Response → User-Friendly Errors
```

### 검증 레이어
```
┌──────────────────┐
│   HTTP Request   │
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  Binding Layer   │ → JSON/XML/Form 파싱
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│ Validation Layer │ → 규칙 적용
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│ Translation Layer│ → 에러 메시지 변환
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  Error Response  │
└──────────────────┘
```

## 🛠 구현된 기능

### 1. **기본 검증 태그**
- required: 필수 필드
- email: 이메일 형식
- min/max: 최소/최대 길이
- gte/lte: 숫자 범위
- oneof: 열거형 값
- eqfield: 필드 일치

### 2. **커스텀 검증자**
- strong_password: 강력한 비밀번호
- korean_phone: 한국 전화번호
- postal_code: 우편번호
- category: 카테고리 검증
- before_today: 날짜 검증
- credit_card: 신용카드 번호

### 3. **다국어 지원**
- 한국어 에러 메시지
- 영어 에러 메시지
- 커스텀 메시지 등록

### 4. **고급 검증**
- 중첩 구조체 검증
- 배열/슬라이스 검증
- 맵 검증
- 동적 검증 규칙

## 🎯 주요 API 엔드포인트

### 사용자 등록
```bash
POST /api/v1/register    # 회원가입 검증
POST /api/v1/login       # 로그인 검증
PUT  /api/v1/profile     # 프로필 수정
```

### 상품 관리
```bash
POST /api/v1/products    # 상품 생성
PUT  /api/v1/products/:id # 상품 수정
```

### 주문 처리
```bash
POST /api/v1/orders      # 주문 생성 (중첩 검증)
```

### 검색 및 필터
```bash
GET  /api/v1/search      # 검색 파라미터 검증
```

### 결제 검증
```bash
POST /api/v1/payment     # 신용카드 검증
```

### 동적 검증
```bash
POST /api/v1/dynamic     # 런타임 검증 규칙
```

## 💻 실습 가이드

### 1. 설치 및 실행
```bash
cd gin/20
go mod init validation-example
go get -u github.com/gin-gonic/gin
go get -u github.com/go-playground/validator/v10
go get -u github.com/go-playground/universal-translator
go get -u github.com/go-playground/locales/ko

# 실행
go run main.go
```

### 2. 회원가입 테스트

#### 성공 케이스
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser123",
    "password": "Test1234!@",
    "confirm_password": "Test1234!@",
    "phone": "+821012345678",
    "age": 25,
    "agree_terms": true
  }'

# 응답
{
  "message": "User registered successfully",
  "user": {
    "email": "test@example.com",
    "username": "testuser123"
  }
}
```

#### 실패 케이스 - 약한 비밀번호
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ko" \
  -d '{
    "email": "test@example.com",
    "username": "user",
    "password": "weak",
    "confirm_password": "weak",
    "phone": "01012345678"
  }'

# 응답 (한국어)
{
  "errors": {
    "password": "비밀번호는 대문자, 소문자, 숫자, 특수문자를 포함해야 합니다",
    "phone": "올바른 한국 전화번호 형식이 아닙니다"
  }
}
```

### 3. 상품 생성 검증

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "노트북",
    "price": 1500000,
    "category": "electronics",
    "stock": 10,
    "tags": ["laptop", "computer", "tech"],
    "images": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ]
  }'

# 응답
{
  "message": "Product created successfully",
  "product": {
    "name": "노트북",
    "price": 1500000,
    "category": "electronics"
  }
}
```

### 4. 주문 생성 (중첩 구조체)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_date": "2024-01-15",
    "total_amount": 150000,
    "customer": {
      "name": "김철수",
      "email": "kim@example.com",
      "phone": "+821012345678"
    },
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "price": 50000
      },
      {
        "product_id": 2,
        "quantity": 1,
        "price": 50000
      }
    ],
    "shipping_address": {
      "street": "강남대로 123",
      "city": "서울",
      "postal_code": "06234",
      "country": "KR"
    }
  }'
```

### 5. 검색 파라미터 검증

```bash
# 올바른 검색
curl "http://localhost:8080/api/v1/search?q=laptop&category=electronics&min_price=100000&max_price=2000000&sort=price&order=asc&page=1&limit=10"

# 잘못된 카테고리
curl "http://localhost:8080/api/v1/search?q=laptop&category=invalid"

# 응답
{
  "error": "Invalid category"
}
```

### 6. 신용카드 검증

```bash
# 유효한 카드번호 (테스트용)
curl -X POST http://localhost:8080/api/v1/payment \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": "4532015112830366",
    "holder_name": "John Doe",
    "expiry_date": "12/25",
    "cvv": "123"
  }'

# 무효한 카드번호
curl -X POST http://localhost:8080/api/v1/payment \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": "1234567890123456",
    "holder_name": "John Doe",
    "expiry_date": "12/25",
    "cvv": "123"
  }'

# 응답
{
  "error": "Invalid credit card number"
}
```

### 7. 동적 검증

```bash
curl -X POST http://localhost:8080/api/v1/dynamic \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "age": 25,
    "preferences": {
      "notifications": true,
      "theme": "dark"
    }
  }'
```

## 🔍 코드 하이라이트

### 커스텀 비밀번호 검증자
```go
func strongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool

    if len(password) >= 8 {
        hasMinLen = true
    }

    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
```

### 한국 전화번호 검증
```go
func koreanPhoneNumber(fl validator.FieldLevel) bool {
    phone := fl.Field().String()

    // 010-1234-5678 or 01012345678
    matched, _ := regexp.MatchString(
        `^01[0-9]-?[0-9]{3,4}-?[0-9]{4}$`,
        phone,
    )
    return matched
}
```

### 다국어 에러 메시지
```go
func translateError(err error, trans ut.Translator) map[string]string {
    errs := err.(validator.ValidationErrors)
    result := make(map[string]string)

    for _, e := range errs {
        result[e.Field()] = e.Translate(trans)
    }

    return result
}

// 커스텀 메시지 등록
validate.RegisterTranslation("strong_password", trans,
    func(ut ut.Translator) error {
        return ut.Add("strong_password",
            "비밀번호는 대문자, 소문자, 숫자, 특수문자를 포함해야 합니다",
            true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("strong_password", fe.Field())
        return t
    },
)
```

### 조건부 검증
```go
type ConditionalForm struct {
    Type     string `json:"type" binding:"required,oneof=personal business"`

    // Personal fields
    FirstName string `json:"first_name" binding:"required_if=Type personal"`
    LastName  string `json:"last_name" binding:"required_if=Type personal"`

    // Business fields
    CompanyName string `json:"company_name" binding:"required_if=Type business"`
    TaxID       string `json:"tax_id" binding:"required_if=Type business"`
}
```

## 🎨 검증 태그 치트시트

### 기본 검증
```go
// 필수/선택
`binding:"required"`           // 필수
`binding:"omitempty"`          // 비어있으면 무시

// 문자열
`binding:"min=3,max=10"`       // 길이 제한
`binding:"alpha"`              // 알파벳만
`binding:"alphanum"`           // 알파벳+숫자
`binding:"email"`              // 이메일
`binding:"url"`                // URL

// 숫자
`binding:"gte=0,lte=100"`      // 범위
`binding:"gt=0,lt=100"`        // 범위 (경계 제외)

// 날짜
`binding:"datetime=2006-01-02"` // 특정 형식

// 배열
`binding:"dive"`               // 각 요소 검증
`binding:"unique"`             // 중복 제거
`binding:"min=1,max=10"`       // 크기 제한
```

### 크로스 필드
```go
// 필드 비교
`binding:"eqfield=Password"`   // 같은 값
`binding:"nefield=OldPassword"` // 다른 값
`binding:"gtfield=MinPrice"`   // 크거나 같음
`binding:"ltefield=MaxPrice"`  // 작거나 같음

// 조건부
`binding:"required_if=Type admin"`
`binding:"required_unless=Guest true"`
`binding:"required_with=Email"`
`binding:"required_without=Phone"`
```

## 📝 베스트 프랙티스

### 1. **명확한 에러 메시지**
```go
// ❌ Bad: 기술적 메시지
"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"

// ✅ Good: 사용자 친화적
"이메일 주소를 입력해주세요"
```

### 2. **적절한 검증 위치**
```go
// ❌ Bad: 핸들러에서 직접 검증
if len(user.Email) == 0 {
    c.JSON(400, gin.H{"error": "Email required"})
    return
}

// ✅ Good: 구조체 태그 사용
type User struct {
    Email string `json:"email" binding:"required,email"`
}
```

### 3. **재사용 가능한 검증자**
```go
// 공통 검증자 등록
func SetupValidators(v *validator.Validate) {
    v.RegisterValidation("phone", phoneValidator)
    v.RegisterValidation("postal", postalValidator)
    v.RegisterValidation("ssn", ssnValidator)
}
```

### 4. **보안 고려사항**
```go
// SQL Injection 방지
`binding:"excludesall='\""`

// XSS 방지
`binding:"html_encoded"`

// 파일 업로드 검증
`binding:"required,file_ext=jpg|png|pdf"`
```

## 🚀 프로덕션 체크리스트

- [ ] 모든 입력에 검증 적용
- [ ] 다국어 에러 메시지 준비
- [ ] 커스텀 검증자 테스트
- [ ] 에러 로깅 구현
- [ ] Rate limiting 적용
- [ ] 입력 크기 제한
- [ ] 파일 업로드 검증
- [ ] SQL Injection 방지

## 🔒 보안 고려사항

### 입력 검증 원칙
- **모든 입력을 의심하라**
- **화이트리스트 방식 사용**
- **클라이언트 검증 신뢰 금지**
- **에러 메시지에 민감정보 노출 금지**

### 추가 보안 레이어
- Content-Type 검증
- 입력 정규화 (normalization)
- 인코딩 검증
- 바이러스 스캔 (파일 업로드)

## 📚 추가 학습 자료

- [go-playground/validator](https://github.com/go-playground/validator)
- [Validator Documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Universal Translator](https://github.com/go-playground/universal-translator)
- [OWASP Input Validation](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html)

## 🎯 다음 레슨 예고

**Lesson 21: 파일 업로드 처리**
- 단일/다중 파일 업로드
- 파일 타입 검증
- 크기 제한
- 바이러스 스캔
- S3/클라우드 저장소 연동

강력한 입력 검증으로 안전한 API를 구축하세요! 🛡️