# RESTful API 만들기: HTTP 라우팅 완벽 가이드 🛣️

이전 챕터에서는 간단한 GET 요청만 다뤘죠? 이번에는 **진짜 API**를 만들어볼 거예요! 사용자를 생성하고, 조회하고, 수정하고, 삭제하는 완전한 API를 만들어봅시다.

## 무엇을 배우게 될까요?

실제 서비스에서 사용하는 것처럼, 데이터를 다루는 4가지 기본 동작을 배워요:
- 📖 **조회하기** (GET): 데이터 읽기
- ✍️ **만들기** (POST): 새로운 데이터 추가
- ✏️ **수정하기** (PUT): 기존 데이터 변경
- 🗑️ **삭제하기** (DELETE): 데이터 제거

이걸 개발자들은 **CRUD**(Create, Read, Update, Delete)라고 부릅니다!

## 이번 챕터에서 배울 내용
- 다양한 HTTP 메서드 사용하기 (GET, POST, PUT, DELETE)
- RESTful API란 무엇인지 이해하기
- 실전 사용자 관리 API 만들기
- 올바른 HTTP 상태 코드 사용하기
- 메모리에 데이터 저장하기 (간단한 데이터베이스처럼!)

## 📂 파일 구조
```
02/
└── main.go     # 라우팅 예제 서버
```

## 💻 코드 이해하기

### 핵심 개념 살펴보기

#### 1️⃣ User 구조체 - 사용자 데이터 모델
```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```
실제 서비스의 "사용자"를 코드로 표현한 거예요. 마치 엑셀의 한 행처럼, 사용자 한 명의 정보를 담고 있죠!

#### 2️⃣ 메모리 저장소 - 간단한 데이터베이스
```go
var users []User  // 사용자들을 담는 배열
```
진짜 데이터베이스를 사용하기 전에, 메모리에 임시로 데이터를 저장해요. 서버를 끄면 데이터가 사라지지만, 연습하기엔 충분해요!

#### 3️⃣ RESTful 라우팅 - URL 체계적으로 만들기
```
/users          → 모든 사용자 목록
/users/:id      → 특정 사용자 한 명
/users/:id      → 사용자 정보 수정
/users/:id      → 사용자 삭제
```
URL만 봐도 무슨 기능인지 알 수 있죠? 이게 바로 RESTful API의 철학이에요!

### API 엔드포인트

| 메서드 | 경로 | 설명 | 상태 코드 |
|--------|------|------|-----------|
| GET | /users | 모든 사용자 조회 | 200 OK |
| GET | /users/:id | 특정 사용자 조회 | 200 OK / 404 Not Found |
| POST | /users | 새 사용자 생성 | 201 Created / 400 Bad Request |
| PUT | /users/:id | 사용자 정보 수정 | 200 OK / 404 Not Found |
| DELETE | /users/:id | 사용자 삭제 | 204 No Content / 404 Not Found |

## 🚀 실행 방법

### 1. 서버 시작
```bash
# gin 폴더에서 실행
cd gin
go run ./02

# 또는 02 폴더로 이동 후 실행
cd gin/02
go run main.go
```

### 2. API 테스트

**모든 사용자 조회:**
```bash
curl http://localhost:3001/users

# 응답:
# []  (초기 상태)
```

**새 사용자 생성:**
```bash
curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"홍길동","email":"hong@example.com"}'

# 응답:
# {
#   "id": "1",
#   "user": {
#     "id": "1",
#     "name": "홍길동",
#     "email": "hong@example.com"
#   }
# }
```

**특정 사용자 조회:**
```bash
curl http://localhost:3001/users/1

# 응답:
# {
#   "id": "1",
#   "name": "홍길동",
#   "email": "hong@example.com"
# }
```

**사용자 정보 수정 (PUT - 전체 수정):**
```bash
curl -X PUT http://localhost:3001/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"김철수","email":"kim@example.com"}'

# 응답:
# {
#   "user": {
#     "id": "1",
#     "name": "김철수",
#     "email": "kim@example.com"
#   }
# }
```

**사용자 삭제:**
```bash
curl -X DELETE http://localhost:3001/users/1

# 응답: 204 No Content (본문 없음)
```

## 💡 꼭 이해하고 넘어가야 할 개념!

### RESTful API란 무엇일까요?

REST는 웹 API를 만드는 하나의 "스타일"이에요. 마치 건축에 여러 양식이 있듯이, API를 만드는 좋은 방법론이죠!

#### RESTful API의 4가지 원칙

**1. 리소스 중심으로 URL 만들기**
```
❌ 나쁜 예: /getUser, /createUser, /deleteUser
✅ 좋은 예: /users  (동작은 HTTP 메서드로!)
```

**2. HTTP 메서드로 동작 구분하기**
```go
GET    /users     → 조회 (읽기)
POST   /users     → 생성 (쓰기)
PUT    /users/:id → 수정
DELETE /users/:id → 삭제
```

**3. 적절한 상태 코드 사용하기**
```
200 OK          → 성공했어요!
201 Created     → 새로 만들었어요!
404 Not Found   → 찾을 수 없어요!
```

**4. 일관성 있게 만들기**
모든 API가 같은 패턴을 따르면, 사용하기 쉬워져요!

### HTTP 메서드 쉽게 이해하기 🎯

각 메서드를 실생활에 비유해볼게요!

```go
// 📖 GET: 책을 "읽는" 것처럼, 데이터를 조회만 해요
r.GET("/users", getAllUsers)
// "사용자 목록 좀 보여주세요!"

// ✍️ POST: 새 글을 "쓰는" 것처럼, 새로운 데이터를 만들어요
r.POST("/users", createUser)
// "새 사용자를 등록할게요!"

// ✏️ PUT: 글을 "수정하는" 것처럼, 기존 데이터를 바꿔요
r.PUT("/users/:id", updateUser)
// "이 사용자 정보를 바꿀게요!"

// 🗑️ DELETE: 글을 "삭제하는" 것처럼, 데이터를 지워요
r.DELETE("/users/:id", deleteUser)
// "이 사용자를 삭제할게요!"
```

💡 **기억하기**: 메서드가 "동사"의 역할을 하고, URL이 "명사"의 역할을 해요!

## 🧪 테스트 시나리오

### 전체 CRUD 플로우 테스트
```bash
# 1. 초기 상태 확인
curl http://localhost:3001/users

# 2. 사용자 3명 생성
curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자1","email":"user1@test.com"}'

curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자2","email":"user2@test.com"}'

curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자3","email":"user3@test.com"}'

# 3. 모든 사용자 확인
curl http://localhost:3001/users

# 4. 특정 사용자 수정
curl -X PUT http://localhost:3001/users/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"수정된사용자2","email":"modified@test.com"}'

# 5. 사용자 삭제
curl -X DELETE http://localhost:3001/users/1

# 6. 최종 상태 확인
curl http://localhost:3001/users
```

## 🔍 트러블슈팅

### JSON 파싱 에러
```bash
# Content-Type 헤더 확인
-H "Content-Type: application/json"

# JSON 형식 검증
# 올바른 형식: {"name":"test","email":"test@test.com"}
# 잘못된 형식: {name:"test",email:"test@test.com"}  # 따옴표 누락
```

### 404 Not Found 에러
```bash
# ID 존재 여부 확인
curl http://localhost:3001/users  # 모든 사용자 목록 확인

# URL 경로 확인
# 올바른 경로: /users/1
# 잘못된 경로: /user/1, /users/1/
```

## 🏗️ 확장 아이디어

1. **검증 추가**: 이메일 형식 검증, 필수 필드 체크
2. **페이지네이션**: GET /users?page=1&limit=10
3. **정렬**: GET /users?sort=name&order=desc
4. **필터링**: GET /users?name=김&email=gmail
5. **데이터베이스 연동**: 메모리 대신 실제 DB 사용

## 📚 다음 단계
- [03. 파라미터 바인딩](../03/README.md): Path, Query, Body 파라미터 처리
- 미들웨어 구현: 인증, 로깅, CORS 처리
- 데이터베이스 연동: GORM을 사용한 영구 저장소

## 🔗 참고 자료
- [REST API 설계 가이드](https://restfulapi.net/)
- [HTTP 상태 코드 참조](https://developer.mozilla.org/ko/docs/Web/HTTP/Status)
- [Gin 라우팅 문서](https://gin-gonic.com/docs/examples/routes-grouping/)