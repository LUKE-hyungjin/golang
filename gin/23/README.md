# Lesson 23: 린팅과 포맷팅 (golangci-lint & go fmt) 🎨

> 코드 품질 도구로 깨끗하고 일관된 Go 코드 작성하기

## 📌 이번 레슨에서 배우는 내용

코드 품질은 프로젝트의 유지보수성과 직결됩니다. golangci-lint와 go fmt를 활용하여 일관되고 깨끗한 코드를 작성하는 방법을 학습합니다.

### 핵심 학습 목표
- ✅ golangci-lint 설정 및 사용
- ✅ 코드 포맷팅 자동화
- ✅ 커스텀 린터 규칙 설정
- ✅ CI/CD 파이프라인 통합
- ✅ Makefile을 통한 자동화
- ✅ 코드 품질 메트릭

## 🏗 코드 품질 도구 구조

### 도구 체인
```
┌──────────────────┐
│   Source Code    │
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│    go fmt        │ → 코드 포맷팅
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│   goimports      │ → import 정리
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│  golangci-lint   │ → 종합 린팅
└────────┬─────────┘
         ↓
┌────────┴─────────┐
│   Quality Gate   │ → 품질 기준 통과
└──────────────────┘
```

### golangci-lint 구성 요소
```
1. Linters
   - errcheck: 에러 체크
   - gosimple: 코드 단순화
   - govet: Go vet
   - ineffassign: 비효율적 할당
   - staticcheck: 정적 분석
   - unused: 사용되지 않는 코드

2. Formatters
   - gofmt: Go 표준 포맷
   - goimports: import 관리
   - gofumpt: 더 엄격한 포맷

3. Security
   - gosec: 보안 취약점
   - secrets: 하드코딩된 시크릿
```

## 🛠 구현된 기능

### 1. **golangci-lint 설정**
- 포괄적인 린터 설정
- 커스텀 규칙 정의
- 예외 처리 설정
- 심각도 레벨 구성

### 2. **Makefile 자동화**
- 린팅 실행
- 포맷팅 체크
- 테스트 및 커버리지
- 도구 설치

### 3. **코드 스타일 가이드**
- 네이밍 컨벤션
- 에러 처리 패턴
- 주석 작성 규칙
- 패키지 구조

### 4. **CI/CD 통합**
- GitHub Actions 설정
- GitLab CI 설정
- 자동 체크
- PR 검증

## 💻 실습 가이드

### 1. golangci-lint 설치
```bash
# macOS
brew install golangci-lint

# Linux/Windows
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 버전 확인
golangci-lint --version
```

### 2. 프로젝트 초기 설정
```bash
cd gin/23

# golangci-lint 설정 파일 생성
touch .golangci.yml

# Makefile 생성
touch Makefile

# 의존성 설치
go mod init lint-example
go get -u github.com/gin-gonic/gin
```

### 3. 린트 실행
```bash
# 기본 실행
golangci-lint run

# 특정 디렉토리
golangci-lint run ./...

# 자동 수정
golangci-lint run --fix

# 특정 린터만 실행
golangci-lint run --enable-only gofmt,goimports

# 새로운 이슈만 표시
golangci-lint run --new-from-rev=HEAD~1
```

### 4. Makefile 사용
```bash
# 도구 설치
make install-tools

# 린트 실행
make lint

# 자동 수정
make lint-fix

# 포맷 체크
make fmt-check

# 포맷 적용
make fmt

# 전체 체크
make check

# CI 파이프라인
make ci
```

## 🎯 주요 린터 설명

### 1. errcheck - 에러 처리 검사
```go
// ❌ Bad: 에러 무시
file, _ := os.Open("file.txt")
defer file.Close()

// ✅ Good: 에러 처리
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
```

### 2. gosimple - 코드 단순화
```go
// ❌ Bad: 복잡한 표현
if err != nil {
    return err
} else {
    return nil
}

// ✅ Good: 단순한 표현
return err
```

### 3. ineffassign - 비효율적 할당
```go
// ❌ Bad: 사용되지 않는 할당
func process() error {
    err := doSomething()
    err = doAnother() // 이전 err 값이 사용되지 않음
    return err
}

// ✅ Good: 올바른 에러 처리
func process() error {
    if err := doSomething(); err != nil {
        return err
    }
    return doAnother()
}
```

### 4. staticcheck - 정적 분석
```go
// ❌ Bad: 잘못된 Printf 형식
fmt.Printf("%s", 123) // 형식 불일치

// ✅ Good: 올바른 형식
fmt.Printf("%d", 123)
```

### 5. gosec - 보안 검사
```go
// ❌ Bad: SQL Injection 취약점
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)
db.Query(query)

// ✅ Good: Prepared statement
db.Query("SELECT * FROM users WHERE id = ?", userInput)
```

## 🔍 golangci-lint 설정 상세

### 기본 설정 구조
```yaml
# .golangci.yml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor
    - testdata

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true

  gofmt:
    simplify: true

  gocyclo:
    min-complexity: 15

linters:
  disable-all: true
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gofmt
    - goimports
    - gosec
    - misspell
```

### 커스텀 규칙 설정
```yaml
issues:
  exclude-rules:
    # 테스트 파일에서 일부 린터 제외
    - path: _test\.go
      linters:
        - gomnd
        - funlen
        - dupl

    # 특정 파일 제외
    - path: internal/generated/
      linters:
        - all
```

## 📊 코드 품질 메트릭

### 1. 복잡도 측정
```bash
# Cyclomatic complexity
golangci-lint run --enable gocyclo

# Cognitive complexity
golangci-lint run --enable gocognit
```

### 2. 코드 중복 검사
```bash
# 중복 코드 찾기
golangci-lint run --enable dupl
```

### 3. 코드 커버리지
```bash
# 커버리지 측정
go test -cover ./...

# HTML 리포트
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🎨 코드 포맷팅

### gofmt 사용
```bash
# 포맷 체크
gofmt -d .

# 포맷 적용
gofmt -s -w .

# 단순화와 함께
gofmt -s -w .
```

### goimports 사용
```bash
# 설치
go install golang.org/x/tools/cmd/goimports@latest

# 실행
goimports -w .

# 로컬 패키지 그룹화
goimports -local github.com/yourusername -w .
```

### gofumpt (더 엄격한 포맷)
```bash
# 설치
go install mvdan.cc/gofumpt@latest

# 실행
gofumpt -w .
```

## 🚀 CI/CD 통합

### GitHub Actions
```yaml
# .github/workflows/lint.yml
name: Lint

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  golangci:
    name: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=5m
```

### GitLab CI
```yaml
# .gitlab-ci.yml
lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run --timeout 5m
  only:
    - merge_requests
    - main
```

### Pre-commit Hook
```bash
# .git/hooks/pre-commit
#!/bin/sh
echo "Running linters..."
golangci-lint run --fix

if [ $? -ne 0 ]; then
    echo "Linting failed. Please fix the issues and try again."
    exit 1
fi

echo "Linting passed!"
```

## 📝 베스트 프랙티스

### 1. **점진적 도입**
```yaml
# 처음에는 필수 린터만
linters:
  enable:
    - errcheck
    - govet
    - staticcheck

# 점진적으로 추가
linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - gosimple
    - ineffassign
    - unused
```

### 2. **팀 규칙 합의**
```yaml
# 팀에서 합의한 규칙만 적용
linters-settings:
  funlen:
    lines: 100      # 팀 합의
    statements: 50  # 팀 합의

  lll:
    line-length: 120  # 팀 합의
```

### 3. **예외 처리 문서화**
```go
// nolint:errcheck // 이 에러는 무시해도 안전함
defer file.Close()

// nolint:gosec // 이미 검증된 입력
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", validatedID)
```

### 4. **정기적 업데이트**
```bash
# golangci-lint 업데이트
brew upgrade golangci-lint

# 새 린터 확인
golangci-lint linters

# 설정 검증
golangci-lint run --print-config
```

## 🔒 보안 린팅

### gosec 규칙
```yaml
linters-settings:
  gosec:
    excludes:
      - G104  # Unhandled errors
      - G304  # File path traversal
    config:
      G301: "0600"  # File permissions
```

### 민감 정보 검사
```bash
# 하드코딩된 시크릿 찾기
golangci-lint run --enable gosec

# Git secrets 사용
git secrets --scan
```

## 📊 품질 리포트

### SonarQube 통합
```bash
# 리포트 생성
golangci-lint run --out-format checkstyle > report.xml

# SonarQube에 전송
sonar-scanner \
  -Dsonar.go.golangci-lint.reportPaths=report.xml
```

### 품질 대시보드
```bash
# 메트릭 수집
go test -cover -json > test.json
golangci-lint run --out-format json > lint.json

# 시각화 도구 사용
# - GolangCI-Lint Web UI
# - SonarQube
# - CodeClimate
```

## 🎯 체크리스트

- [ ] golangci-lint 설치
- [ ] .golangci.yml 설정
- [ ] Makefile 작성
- [ ] Pre-commit hook 설정
- [ ] CI/CD 파이프라인 구성
- [ ] 팀 코딩 스타일 가이드 작성
- [ ] 정기적 린터 업데이트
- [ ] 품질 메트릭 모니터링

## 📚 추가 학습 자료

- [golangci-lint](https://golangci-lint.run/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

## 🎯 실습 과제

1. **프로젝트에 golangci-lint 적용**
   - 설정 파일 작성
   - 모든 이슈 해결
   - CI 파이프라인 구성

2. **커스텀 린터 규칙**
   - 팀 규칙 정의
   - 예외 처리 설정
   - 문서화

3. **자동화 구축**
   - Makefile 작성
   - Pre-commit hook
   - 자동 수정 스크립트

깨끗한 코드로 프로젝트 품질을 높이세요! 🎨