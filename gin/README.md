# Gin 학습 체크리스트

> Gin(Golang) 웹 프레임워크를 체계적으로 학습하기 위한 가이드입니다. 완료 시 체크하세요.

## 시작하기
- [x] Go 설치 및 환경 변수(GOROOT, GOPATH) 점검
- [x] 모듈 초기화(go mod init)와 기본 빌드/실행(go run, go build)
- [x] IDE/에디터 설정(goimports, gopls, golangci-lint)

## Gin 핵심
- [x] Gin 설치 및 기본 서버 구동
- [x] 라우팅(GET, POST, PUT, DELETE)
- [x] Path/Query/Body 파라미터 바인딩
- [x] 컨텍스트 사용(c.Param, c.Query, c.Bind, c.JSON)
- [x] 미들웨어(전역/그룹/개별)와 next() 흐름
- [x] 라우트 그룹, 버저닝(v1, v2)
- [x] 정적 파일 서빙(Static, StaticFS)
- [x] 템플릿 렌더링(HTML)

## 에러 처리 & 로깅
- [ ] HTTP 상태코드와 에러 응답 규약
- [ ] 에러 핸들링 미들웨어
- [ ] 로깅 미들웨어(요청/응답 로깅)

## 구성 & 설정
- [ ] 환경변수/설정파일(viper 등) 로딩
- [ ] 의존성 주입(인터페이스/팩토리 분리)
- [ ] 실행 모드(release/debug/test)

## 데이터 계층
- [ ] GORM(or sqlx) 연결 및 CRUD
- [ ] 마이그레이션과 시드 데이터
- [ ] 트랜잭션, 컨텍스트 타임아웃

## 보안
- [ ] CORS 설정
- [ ] 인증/인가(JWT, 세션) 미들웨어
- [ ] 입력 검증(binding + validator)

## 테스트 & 품질
- [ ] 핸들러 유닛 테스트(httptest)
- [ ] 통합 테스트(테스트 라우터 구성)
- [ ] 린팅(golangci-lint)과 포맷팅(go fmt)

---

## 현재 폴더 구조
```
.
01/                  # 기본 서버 구동
├── main.go
└── README.md
02/                  # 라우팅 기본
├── main.go
└── README.md
03/                  # 파라미터 바인딩
├── main.go
└── README.md
04/                  # Context 사용
├── main.go
└── README.md
05/                  # 미들웨어
├── main.go
└── README.md
06/                  # 라우트 그룹과 버저닝
├── main.go
└── README.md
07/                  # 정적 파일 서빙
├── main.go
├── README.md
├── static/
│   ├── index.html
│   ├── css/
│   ├── js/
│   └── ...
└── uploads/
08/                  # 템플릿 렌더링
├── main.go
├── README.md
└── templates/
    ├── index.html
    ├── users.html
    ├── products.html
    └── ...
project/             # 📌 통합 프로젝트 (01~08 모든 내용 포함)
├── cmd/
│   └── main.go     # 메인 애플리케이션
├── internal/       # 내부 패키지
│   ├── handlers/   # HTTP 핸들러
│   ├── middleware/ # 미들웨어
│   └── models/     # 데이터 모델
├── web/           # 웹 리소스
│   ├── templates/ # HTML 템플릿
│   └── static/    # 정적 파일
├── Documentation/ # 프로젝트 문서
│   ├── README.md
│   ├── API.md
│   └── QUICKSTART.md
├── Makefile
└── go.mod
README.md
go.mod
go.sum
```

## 실행 방법
```bash
# gin 폴더로 이동
cd gin

# 각 예제 실행 (포트 8080 사용)
go run ./01  # 01: 기본 서버
go run ./02  # 02: 라우팅 기본
go run ./03  # 03: 파라미터 바인딩
go run ./04  # 04: Context 사용
go run ./05  # 05: 미들웨어
go run ./06  # 06: 라우트 그룹과 버저닝
go run ./07  # 07: 정적 파일 서빙
go run ./08  # 08: 템플릿 렌더링

# 브라우저에서 접속
http://localhost:8080

# API 테스트 예시
curl http://localhost:8080/ping
curl http://localhost:8080/users/123
curl "http://localhost:8080/search?q=hello"
curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{"name":"Jin","email":"jin@example.com"}'
```

## 학습 순서
1. **01 - 기본 서버**: Gin 설치와 Hello World
2. **02 - 라우팅**: HTTP 메서드별 라우팅과 RESTful API
3. **03 - 파라미터 바인딩**: Path, Query, Body 파라미터 처리
4. **04 - Context**: gin.Context의 모든 기능
5. **05 - 미들웨어**: 인증, 로깅, CORS 등 미들웨어 구현
6. **06 - 라우트 그룹**: API 버저닝과 그룹화
7. **07 - 정적 파일**: 파일 업로드/다운로드, SPA 지원
8. **08 - 템플릿**: HTML 템플릿 렌더링과 웹 애플리케이션

### 🚀 통합 프로젝트
**project 폴더**: 01~08의 모든 내용을 통합한 실전 블로그/커뮤니티 플랫폼
- 사용자 인증 시스템
- 블로그 포스트 CRUD
- 댓글 시스템
- 파일 업로드
- 관리자 대시보드
- API v1/v2 버저닝
- 완전한 웹 인터페이스

```bash
# 통합 프로젝트 실행
cd gin/project
make install
make run
# http://localhost:8080 접속
```

각 폴더의 README.md에서 상세한 설명과 테스트 방법을 확인할 수 있습니다.