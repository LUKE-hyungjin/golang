# 02. HTTP 라우팅 기본

## 📌 개요
Gin 프레임워크의 라우팅 시스템을 학습합니다. GET, POST, PUT, DELETE 등 다양한 HTTP 메서드별 라우팅 구현과 기본적인 RESTful API 패턴을 익힙니다.

## 🎯 학습 목표
- HTTP 메서드별 라우팅 구현 (GET, POST, PUT, DELETE, PATCH)
- RESTful API 설계 원칙 이해
- 간단한 CRUD 작업 구현
- HTTP 상태 코드 적절한 사용
- 메모리 기반 데이터 저장소 구현

## 📂 파일 구조
```
02/
└── main.go     # 라우팅 예제 서버
```

## 💻 코드 설명

### 주요 구성 요소

1. **User 구조체**: 사용자 데이터 모델
2. **메모리 저장소**: `map`을 사용한 간단한 데이터 저장
3. **뮤텍스**: 동시성 제어를 위한 sync.RWMutex
4. **RESTful 라우팅**: 리소스 기반 URL 설계

### API 엔드포인트

| 메서드 | 경로 | 설명 | 상태 코드 |
|--------|------|------|-----------|
| GET | /users | 모든 사용자 조회 | 200 OK |
| GET | /users/:id | 특정 사용자 조회 | 200 OK / 404 Not Found |
| POST | /users | 새 사용자 생성 | 201 Created / 400 Bad Request |
| PUT | /users/:id | 사용자 정보 수정 | 200 OK / 404 Not Found |
| DELETE | /users/:id | 사용자 삭제 | 204 No Content / 404 Not Found |
| PATCH | /users/:id | 사용자 부분 수정 | 200 OK / 404 Not Found |

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
curl http://localhost:8080/users

# 응답:
# []  (초기 상태)
```

**새 사용자 생성:**
```bash
curl -X POST http://localhost:8080/users \
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
curl http://localhost:8080/users/1

# 응답:
# {
#   "id": "1",
#   "name": "홍길동",
#   "email": "hong@example.com"
# }
```

**사용자 정보 수정 (PUT - 전체 수정):**
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"김철수","email":"kim@example.com"}'

# 응답:
# {
#   "message": "User updated",
#   "user": {
#     "id": "1",
#     "name": "김철수",
#     "email": "kim@example.com"
#   }
# }
```

**사용자 정보 부분 수정 (PATCH):**
```bash
curl -X PATCH http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'

# 응답:
# {
#   "message": "User partially updated",
#   "user": {
#     "id": "1",
#     "name": "김철수",
#     "email": "new@example.com"
#   }
# }
```

**사용자 삭제:**
```bash
curl -X DELETE http://localhost:8080/users/1

# 응답: 204 No Content (본문 없음)
```

## 📝 주요 학습 포인트

### 1. RESTful API 설계 원칙
- **리소스 중심**: URL은 리소스를 나타냄 (/users, /users/:id)
- **HTTP 메서드 활용**: 동작은 메서드로 표현 (GET, POST, PUT, DELETE)
- **상태 코드**: 적절한 HTTP 상태 코드 반환
- **일관성**: 예측 가능한 API 구조

### 2. HTTP 메서드별 용도
```go
// GET: 리소스 조회 (Read)
r.GET("/users", getAllUsers)

// POST: 새 리소스 생성 (Create)
r.POST("/users", createUser)

// PUT: 리소스 전체 수정 (Update - Full)
r.PUT("/users/:id", updateUser)

// PATCH: 리소스 부분 수정 (Update - Partial)
r.PATCH("/users/:id", patchUser)

// DELETE: 리소스 삭제 (Delete)
r.DELETE("/users/:id", deleteUser)
```

### 3. 상태 코드 사용
```go
// 200 OK: 요청 성공
c.JSON(http.StatusOK, data)

// 201 Created: 리소스 생성 성공
c.JSON(http.StatusCreated, newResource)

// 204 No Content: 성공했지만 반환할 내용 없음
c.Status(http.StatusNoContent)

// 400 Bad Request: 잘못된 요청
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})

// 404 Not Found: 리소스를 찾을 수 없음
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
```

### 4. 동시성 제어
```go
var mu sync.RWMutex  // 읽기/쓰기 뮤텍스

// 읽기 작업 시
mu.RLock()
defer mu.RUnlock()

// 쓰기 작업 시
mu.Lock()
defer mu.Unlock()
```

## 🧪 테스트 시나리오

### 전체 CRUD 플로우 테스트
```bash
# 1. 초기 상태 확인
curl http://localhost:8080/users

# 2. 사용자 3명 생성
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자1","email":"user1@test.com"}'

curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자2","email":"user2@test.com"}'

curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"사용자3","email":"user3@test.com"}'

# 3. 모든 사용자 확인
curl http://localhost:8080/users

# 4. 특정 사용자 수정
curl -X PUT http://localhost:8080/users/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"수정된사용자2","email":"modified@test.com"}'

# 5. 사용자 삭제
curl -X DELETE http://localhost:8080/users/1

# 6. 최종 상태 확인
curl http://localhost:8080/users
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
curl http://localhost:8080/users  # 모든 사용자 목록 확인

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