# 01. Gin 기본 서버 구동

## 📌 개요
Gin 웹 프레임워크를 사용한 가장 기본적인 HTTP 서버 구현입니다. "Hello, World" 수준의 첫 번째 예제로, Gin의 기본 구조와 서버 시작 방법을 학습합니다.

## 🎯 학습 목표
- Gin 프레임워크 설치 및 임포트
- 기본 라우터 생성
- 간단한 GET 엔드포인트 구현
- JSON 응답 반환
- 서버 실행 및 테스트

## 📂 파일 구조
```
01/
└── main.go     # 메인 서버 파일
```

## 💻 코드 설명

### 주요 구성 요소

1. **gin.Default()**: 기본 미들웨어(Logger, Recovery)가 포함된 라우터 생성
2. **r.GET()**: HTTP GET 메서드 핸들러 등록
3. **c.JSON()**: JSON 형식으로 응답 반환
4. **gin.H{}**: `map[string]interface{}`의 단축 표현
5. **r.Run()**: 서버 시작 (기본 포트: 8080)

### 엔드포인트

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | /ping | 서버 상태 확인 (pong 응답) |
| GET | / | 환영 메시지 반환 |

## 🚀 실행 방법

### 1. 서버 시작
```bash
# gin 폴더에서 실행
cd gin
go run ./01

# 또는 01 폴더로 이동 후 실행
cd gin/01
go run main.go
```

### 2. 서버 테스트

**브라우저에서 접속:**
- http://localhost:8080
- http://localhost:8080/ping

**curl 명령어 사용:**
```bash
# 루트 경로 테스트
curl http://localhost:8080

# 응답:
# {"message":"Welcome to Gin server!"}

# /ping 경로 테스트
curl http://localhost:8080/ping

# 응답:
# {"message":"pong"}
```

**HTTPie 사용 (설치된 경우):**
```bash
http GET localhost:8080
http GET localhost:8080/ping
```

## 📝 주요 학습 포인트

### 1. Gin의 장점
- **빠른 성능**: 고성능 HTTP 라우터 사용
- **간단한 API**: 직관적인 메서드 체이닝
- **미들웨어 지원**: Logger, Recovery 등 기본 제공
- **JSON 처리**: 간편한 JSON 직렬화/역직렬화

### 2. gin.Default() vs gin.New()
```go
// 기본 미들웨어 포함 (추천)
r := gin.Default()  // Logger + Recovery 미들웨어 자동 적용

// 미들웨어 없는 순수 라우터
r := gin.New()      // 필요한 미들웨어를 수동으로 추가해야 함
```

### 3. 응답 형식
```go
// JSON 응답
c.JSON(200, gin.H{"key": "value"})

// 문자열 응답
c.String(200, "Hello")

// HTML 응답 (템플릿 사용 시)
c.HTML(200, "index.html", gin.H{})
```

## 🔍 트러블슈팅

### 포트 충돌 시
```bash
# 8080 포트가 이미 사용 중인 경우
# main.go에서 포트 변경
r.Run(":3000")  // 3000번 포트로 변경

# 또는 환경변수로 지정
PORT=3000 go run main.go
```

### 서버가 시작되지 않는 경우
1. Go 모듈 초기화 확인: `go mod init`
2. Gin 설치 확인: `go get -u github.com/gin-gonic/gin`
3. 포트 사용 여부 확인: `lsof -i :8080`

## 📚 다음 단계
- [02. 라우팅 기본](../02/README.md): HTTP 메서드별 라우팅 구현
- [03. 파라미터 바인딩](../03/README.md): Path, Query, Body 파라미터 처리

## 🔗 참고 자료
- [Gin 공식 문서](https://gin-gonic.com/docs/)
- [Gin GitHub 저장소](https://github.com/gin-gonic/gin)
- [Go 웹 프로그래밍 가이드](https://golang.org/doc/articles/wiki/)