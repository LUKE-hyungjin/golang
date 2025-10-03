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
- [ ] 컨텍스트 사용(c.Param, c.Query, c.Bind, c.JSON)
- [ ] 미들웨어(전역/그룹/개별)와 next() 흐름
- [ ] 라우트 그룹, 버저닝(v1, v2)
- [ ] 정적 파일 서빙(Static, StaticFS)
- [ ] 템플릿 렌더링(HTML)

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
01
01/main.go
01/README.md
02
02/main.go
02/README.md
03
03/main.go
03/README.md
README.md
cmd
go.mod
go.sum
```