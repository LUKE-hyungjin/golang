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
- [x] HTTP 상태코드와 에러 응답 규약
- [x] 에러 핸들링 미들웨어
- [x] 로깅 미들웨어(요청/응답 로깅)

## 구성 & 설정
- [x] 환경변수/설정파일(viper 등) 로딩
- [x] 의존성 주입(인터페이스/팩토리 분리)
- [x] 실행 모드(release/debug/test)

## 데이터 계층
- [x] GORM(or sqlx) 연결 및 CRUD
- [x] 마이그레이션과 시드 데이터
- [x] 트랜잭션, 컨텍스트 타임아웃

## 보안
- [x] CORS 설정
- [x] 인증/인가(JWT, 세션) 미들웨어
- [x] 입력 검증(binding + validator)

## 테스트 & 품질
- [x] 핸들러 유닛 테스트(httptest)
- [x] 통합 테스트(테스트 라우터 구성)
- [x] 린팅(golangci-lint)과 포맷팅(go fmt)

## 🚀 Advanced - 심화 학습
> 기본기를 모두 익힌 후, 더 큰 규모의 시스템을 구축하기 위한 고급 주제들

### API 문서화
- [ ] **Swagger/OpenAPI 3.0 (swaggo)**: API 문서 자동 생성
  - 코드 주석에서 API 스펙 자동 추출
  - 대화형 API 테스트 UI 제공
  - 클라이언트 SDK 자동 생성 지원
- [ ] **API 버전 관리**: 하위 호환성 유지 전략

### 아키텍처 & 설계
- [ ] **마이크로서비스 아키텍처 (MSA)**: 서비스 분리와 독립적 배포
  - 서비스 간 통신 패턴
  - 분산 트랜잭션 처리
  - 서비스 디스커버리
- [ ] **메시지 큐 (Kafka, RabbitMQ, NATS)**: 비동기 메시지 처리
  - 이벤트 드리븐 아키텍처
  - Pub/Sub 패턴 구현
  - 메시지 신뢰성 보장
- [ ] **gRPC**: 고성능 RPC 프레임워크
  - Protocol Buffers 정의
  - 양방향 스트리밍
  - 서비스 간 효율적 통신
- [ ] **GraphQL**: 유연한 데이터 쿼리
  - 스키마 정의
  - Resolver 구현
  - DataLoader 패턴
- [ ] **WebSocket**: 실시간 양방향 통신
  - 채팅 시스템 구현
  - 실시간 알림
  - 협업 도구 개발

### 인프라 & DevOps
- [ ] **컨테이너 (Docker)**: 애플리케이션 패키징
  - 멀티스테이지 빌드
  - 이미지 최적화
  - Docker Compose 활용
- [ ] **컨테이너 오케스트레이션 (Kubernetes)**: 대규모 컨테이너 관리
  - Deployment, Service, Ingress
  - Auto-scaling (HPA, VPA)
  - ConfigMap, Secret 관리
  - Helm Charts 작성
- [ ] **CI/CD 파이프라인**: 자동화된 배포
  - GitHub Actions / GitLab CI
  - ArgoCD (GitOps)
  - Blue-Green / Canary 배포
  - 자동 롤백 전략
- [ ] **클라우드 플랫폼 (AWS, GCP, Azure)**: 클라우드 네이티브 개발
  - Serverless (Lambda, Cloud Functions)
  - 관리형 서비스 활용 (RDS, S3, CDN)
  - Infrastructure as Code (Terraform)
  - 비용 최적화

### 데이터 관리 심화
- [ ] **캐싱 전략 (Redis, Memcached)**: 성능 최적화
  - 캐시 무효화 전략
  - 분산 캐싱
  - Session Store
  - Pub/Sub 메시징
- [ ] **검색 엔진 (Elasticsearch)**: 고급 검색 기능
  - 풀텍스트 검색
  - 집계(Aggregation)
  - 로그 분석 (ELK Stack)
  - 실시간 인덱싱
- [ ] **NoSQL 데이터베이스**: 유연한 데이터 모델
  - MongoDB (Document DB)
  - DynamoDB (Key-Value)
  - Cassandra (Wide Column)
  - Neo4j (Graph DB)
- [ ] **이벤트 소싱 & CQRS**: 복잡한 도메인 처리
  - 이벤트 스토어 구현
  - 읽기/쓰기 모델 분리
  - 이벤트 재생

### 모니터링 & 관찰성
- [ ] **분산 추적 (Jaeger, Zipkin)**: 마이크로서비스 디버깅
- [ ] **메트릭 수집 (Prometheus + Grafana)**: 시스템 모니터링
- [ ] **로그 집계 (ELK, Fluentd)**: 중앙화된 로깅
- [ ] **APM (Application Performance Monitoring)**: 성능 분석

### 보안 심화
- [ ] **OAuth 2.0 / OIDC**: 소셜 로그인 구현
- [ ] **API Gateway**: 중앙화된 인증/인가
- [ ] **서비스 메시 (Istio, Linkerd)**: 서비스 간 보안 통신
- [ ] **Vault**: 시크릿 관리

### 성능 최적화
- [ ] **로드 밸런싱**: 트래픽 분산
- [ ] **데이터베이스 샤딩**: 수평적 확장
- [ ] **CDN 활용**: 정적 자원 배포
- [ ] **비동기 처리**: 백그라운드 작업

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
09/                  # HTTP 상태코드와 에러 응답
├── main.go
└── README.md
10/                  # 에러 핸들링 미들웨어
├── main.go
└── README.md
11/                  # 로깅 미들웨어
├── main.go
├── README.md
└── logs/            # 로그 파일 (자동 생성)
12/                  # 환경변수/설정파일 로딩 (Viper)
├── main.go
├── README.md
├── config/
│   ├── config.yaml
│   ├── config.development.yaml
│   └── config.production.yaml
└── .env.example
13/                  # 의존성 주입 (DI)
├── main.go
└── README.md
14/                  # 실행 모드 (Release/Debug/Test)
├── main.go
└── README.md
15/                  # GORM과 SQLite CRUD
├── main.go
├── README.md
└── blog.db         # SQLite 데이터베이스 (자동 생성)
16/                  # 마이그레이션과 시드 데이터
├── main.go
├── README.md
└── blog.db         # SQLite 데이터베이스 (자동 생성)
17/                  # 트랜잭션과 컨텍스트 타임아웃
├── main.go
├── README.md
└── transaction.db  # SQLite 데이터베이스 (자동 생성)
18/                  # CORS 설정
├── main.go
└── README.md
19/                  # JWT 인증 미들웨어
├── main.go
└── README.md
20/                  # 입력 검증 (Binding + Validator)
├── main.go
└── README.md
21/                  # 핸들러 유닛 테스트 (httptest)
├── main.go
└── README.md
22/                  # 통합 테스트
├── main.go
└── README.md
23/                  # 린팅과 포맷팅 (golangci-lint)
├── main.go
├── README.md
├── .golangci.yml
└── Makefile
project(01~08)/             # 📌 통합 프로젝트 (01~08 모든 내용 포함)
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
go run ./09  # 09: HTTP 상태코드와 에러 응답
go run ./10  # 10: 에러 핸들링 미들웨어
go run ./11  # 11: 로깅 미들웨어
go run ./12  # 12: 환경변수/설정파일 로딩 (Viper)
go run ./13  # 13: 의존성 주입 (DI)
go run ./14  # 14: 실행 모드 (Release/Debug/Test)
go run ./15  # 15: GORM과 SQLite CRUD
go run ./16  # 16: 마이그레이션과 시드 데이터
go run ./17  # 17: 트랜잭션과 컨텍스트 타임아웃
go run ./18  # 18: CORS 설정
go run ./19  # 19: JWT 인증 미들웨어
go run ./20  # 20: 입력 검증
go run ./21  # 21: 핸들러 유닛 테스트
go run ./22  # 22: 통합 테스트
go run ./23  # 23: 린팅과 포맷팅

# 브라우저에서 접속
http://localhost:8080

# API 테스트 예시
curl http://localhost:8080/ping
curl http://localhost:8080/users/123
curl "http://localhost:8080/search?q=hello"
curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{"name":"Jin","email":"jin@example.com"}'
```

## 학습 순서

### Gin 핵심 (01~08)
1. **01 - 기본 서버**: Gin 설치와 Hello World
2. **02 - 라우팅**: HTTP 메서드별 라우팅과 RESTful API
3. **03 - 파라미터 바인딩**: Path, Query, Body 파라미터 처리
4. **04 - Context**: gin.Context의 모든 기능
5. **05 - 미들웨어**: 인증, 로깅, CORS 등 미들웨어 구현
6. **06 - 라우트 그룹**: API 버저닝과 그룹화
7. **07 - 정적 파일**: 파일 업로드/다운로드, SPA 지원
8. **08 - 템플릿**: HTML 템플릿 렌더링과 웹 애플리케이션

### 에러 처리 & 로깅 (09~11)
9. **09 - HTTP 상태코드**: 표준 에러 응답 규약과 상태 코드 활용
10. **10 - 에러 핸들링**: 전역 에러 처리와 패닉 복구 미들웨어
11. **11 - 로깅**: 구조화된 로깅과 요청/응답 추적

### 구성 & 설정 (12~14)
12. **12 - 환경변수/설정파일**: Viper를 사용한 설정 관리와 환경별 구성
13. **13 - 의존성 주입**: 인터페이스 기반 설계와 팩토리 패턴
14. **14 - 실행 모드**: Debug/Release/Test 모드별 최적화

### 데이터 계층 (15~17)
15. **15 - GORM과 SQLite CRUD**: ORM을 사용한 데이터베이스 작업과 Repository 패턴
16. **16 - 마이그레이션과 시드 데이터**: 스키마 버전 관리와 테스트 데이터 생성
17. **17 - 트랜잭션과 컨텍스트 타임아웃**: ACID 트랜잭션과 동시성 제어

### 보안 (18~20)
18. **18 - CORS 설정**: Cross-Origin Resource Sharing 구성과 환경별 설정
19. **19 - JWT 인증 미들웨어**: Access/Refresh 토큰과 역할 기반 접근 제어
20. **20 - 입력 검증**: 구조체 태그 검증과 커스텀 검증자

### 테스트 & 품질 (21~23)
21. **21 - 핸들러 유닛 테스트**: httptest를 사용한 HTTP 핸들러 테스트
22. **22 - 통합 테스트**: 데이터베이스와 전체 스택을 포함한 통합 테스트
23. **23 - 린팅과 포맷팅**: golangci-lint와 코드 품질 도구

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