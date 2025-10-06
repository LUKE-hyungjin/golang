# Gin 웹 프레임워크로 첫 서버 만들기 🚀

안녕하세요! 오늘은 Go 언어의 가장 인기 있는 웹 프레임워크인 **Gin**을 사용해서 여러분의 첫 번째 웹 서버를 만들어보겠습니다.

## 왜 Gin을 배워야 할까요?

Gin은 Go로 웹 서버를 만들 때 가장 많이 사용되는 프레임워크입니다. 마치 Python의 Flask나 Node.js의 Express처럼, 빠르고 간단하게 API 서버를 만들 수 있죠.

### Gin의 장점
- ⚡ **엄청 빠른 성능**: Gin은 속도에 최적화되어 있어요
- 🎯 **간단한 문법**: 몇 줄만으로도 서버를 만들 수 있어요
- 🛠 **풍부한 기능**: 로깅, 에러 처리 등이 기본으로 제공돼요
- 📦 **JSON 처리**: API 만들기에 최적화되어 있어요

## 이번 챕터에서 배울 내용
- Gin 프레임워크 설치하고 시작하기
- 기본 라우터 만들기
- 간단한 GET 요청 처리하기
- JSON 형식으로 응답 보내기
- 서버 실행하고 테스트하기

## 📂 파일 구조
```
01/
└── main.go     # 메인 서버 파일
```

## 💻 코드 이해하기

### 핵심 개념 설명

웹 서버를 만들 때 꼭 알아야 할 핵심 요소들을 설명해드릴게요!

#### 1. `gin.Default()` - 서버의 시작점
```go
r := gin.Default()
```
이 한 줄이 여러분의 웹 서버를 만들어줍니다! `Default()`는 기본적으로 필요한 기능들을 자동으로 포함해요:
- 📝 **Logger**: 모든 요청을 자동으로 기록해줘요
- 🛡 **Recovery**: 에러가 나도 서버가 죽지 않게 보호해줘요

#### 2. `r.GET()` - 요청 받기
```go
r.GET("/ping", func(c *gin.Context) {
    // 여기서 요청을 처리해요
})
```
GET 요청을 받는 경로를 만들어요. `/ping` 경로로 들어온 요청을 처리하겠다는 뜻이죠!

#### 3. `c.JSON()` - JSON으로 응답하기
```go
c.JSON(200, gin.H{"message": "pong"})
```
클라이언트에게 JSON 형식으로 데이터를 보내요. 200은 "성공"을 의미하는 상태 코드예요.

#### 4. `gin.H{}` - 간편한 데이터 구조
```go
gin.H{"key": "value"}
```
`gin.H`는 JSON 객체를 만들 때 쓰는 편리한 방법이에요. Go의 `map[string]interface{}`와 같지만 훨씬 간단하죠!

#### 5. `r.Run()` - 서버 실행!
```go
r.Run(":8080")
```
드디어 서버를 시작합니다! 기본 포트는 8080번이에요.

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

## 💡 실습하면서 알아두면 좋은 팁!

### 1️⃣ `gin.Default()` vs `gin.New()` - 뭘 써야 할까요?

**초보자라면 무조건 `gin.Default()`를 사용하세요!**

```go
// 🟢 초보자 추천: 필요한 기능이 다 들어있어요
r := gin.Default()

// 🔴 고급 사용자용: 모든 걸 직접 설정해야 해요
r := gin.New()
```

`Default()`는 로깅과 에러 복구 기능이 자동으로 켜져있어서 편해요!

### 2️⃣ 다양한 응답 형식 사용하기

Gin은 여러 가지 형식으로 응답을 보낼 수 있어요:

```go
// 📋 JSON 응답 (API 만들 때 주로 사용)
c.JSON(200, gin.H{"message": "성공!"})

// 📝 문자열 응답 (간단한 텍스트)
c.String(200, "안녕하세요!")

// 🎨 HTML 응답 (웹페이지 보여줄 때)
c.HTML(200, "index.html", gin.H{"title": "내 웹사이트"})
```

대부분의 경우 JSON을 사용하니까, 처음엔 `c.JSON()`만 익혀도 충분해요!

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