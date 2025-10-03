# Lesson 15: GORM과 SQLite를 사용한 CRUD 작업 📊

> Go의 대표적인 ORM인 GORM과 SQLite를 활용한 데이터베이스 작업 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

데이터베이스는 웹 애플리케이션의 핵심입니다. GORM은 Go의 가장 인기 있는 ORM(Object-Relational Mapping) 라이브러리로, 데이터베이스 작업을 간편하게 만들어줍니다. 이번 레슨에서는 SQLite를 사용하여 CRUD 작업을 구현합니다.

### 핵심 학습 목표
- ✅ GORM 기본 설정과 모델 정의
- ✅ CRUD (Create, Read, Update, Delete) 작업
- ✅ 관계 설정 (1:N, N:M)
- ✅ 페이지네이션과 필터링
- ✅ Repository 패턴 구현
- ✅ 소프트 삭제와 하드 삭제

## 🏗 데이터베이스 구조

### 모델 관계도
```
User (1) ──── (N) Post
  │                 │
  │                 ├── (N) Comment
  │                 │
  └── (N) Comment   └── (N:M) Tag
                    └── (N:1) Category
```

### 테이블 구조
```sql
-- Users 테이블
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    name TEXT,
    age INTEGER,
    bio TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Posts 테이블
CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    slug TEXT UNIQUE NOT NULL,
    published BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    user_id INTEGER REFERENCES users(id),
    category_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🛠 구현된 기능

### 1. **모델 정의**
```go
type User struct {
    Base
    Email    string `gorm:"uniqueIndex;not null"`
    Username string `gorm:"uniqueIndex;not null;size:50"`
    Name     string `gorm:"size:100"`
    Posts    []Post `gorm:"foreignKey:UserID"`
}
```

### 2. **Repository 패턴**
- UserRepository: 사용자 CRUD 작업
- PostRepository: 포스트 CRUD 작업
- 관심사 분리와 테스트 용이성

### 3. **고급 쿼리 기능**
- 연관 데이터 Preload
- 조건부 필터링
- 페이지네이션
- 검색 기능
- 집계 쿼리

## 🎯 주요 API 엔드포인트

### 사용자 관리
```bash
POST   /users          # 사용자 생성
GET    /users          # 사용자 목록 (페이지네이션)
GET    /users/:id      # 사용자 상세 조회
PUT    /users/:id      # 사용자 수정
DELETE /users/:id      # 사용자 삭제
```

### 포스트 관리
```bash
POST   /posts          # 포스트 생성
GET    /posts          # 포스트 목록 (필터링, 페이지네이션)
GET    /posts/:id      # 포스트 상세 조회
GET    /posts/slug/:slug # Slug로 포스트 조회
PUT    /posts/:id      # 포스트 수정
DELETE /posts/:id      # 포스트 삭제
```

### 검색 및 필터
```bash
GET    /search?q=keyword     # 포스트 검색
GET    /popular?limit=10     # 인기 포스트
```

## 💻 실습 가이드

### 1. 설치 및 실행
```bash
# 의존성 설치
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite

# 실행
cd gin/15
go run main.go

# SQLite 데이터베이스 파일이 자동 생성됨
ls -la blog.db
```

### 2. 사용자 CRUD

#### 사용자 생성
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "username": "alice",
    "name": "Alice Kim",
    "age": 28,
    "bio": "Software Engineer"
  }'

# 여러 사용자 생성
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "username": "bob",
    "name": "Bob Lee",
    "age": 32
  }'
```

#### 사용자 조회
```bash
# 단일 사용자 조회 (연관 데이터 포함)
curl http://localhost:8080/users/1 | jq

# 사용자 목록 (페이지네이션)
curl "http://localhost:8080/users?page=1&page_size=10" | jq

# 응답 예시
{
  "users": [...],
  "total": 15,
  "page": 1,
  "page_size": 10,
  "total_pages": 2
}
```

#### 사용자 수정
```bash
# 특정 필드만 수정
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Park",
    "bio": "Senior Software Engineer"
  }'
```

#### 사용자 삭제
```bash
# 소프트 삭제 (기본)
curl -X DELETE http://localhost:8080/users/1

# 하드 삭제 (완전 삭제)
curl -X DELETE "http://localhost:8080/users/1?hard=true"
```

### 3. 포스트 CRUD

#### 포스트 생성
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with GORM",
    "content": "GORM is a fantastic ORM library for Go...",
    "user_id": 1,
    "published": true
  }'

# Slug는 자동 생성됨
```

#### 포스트 조회
```bash
# ID로 조회 (조회수 자동 증가)
curl http://localhost:8080/posts/1 | jq

# Slug로 조회
curl http://localhost:8080/posts/slug/getting-started-with-gorm-1234567890 | jq

# 필터링된 목록
curl "http://localhost:8080/posts?published=true&user_id=1&page=1" | jq
```

### 4. 검색 기능

#### 키워드 검색
```bash
# 제목 또는 내용에서 검색
curl "http://localhost:8080/search?q=GORM" | jq

# 응답 예시
{
  "keyword": "GORM",
  "results": [
    {
      "id": 1,
      "title": "Getting Started with GORM",
      "content": "...",
      "user": {...}
    }
  ],
  "count": 1
}
```

#### 인기 포스트
```bash
# 조회수 기준 상위 10개
curl http://localhost:8080/popular?limit=10 | jq
```

### 5. 고급 쿼리 예제
```bash
# GORM 고급 쿼리 예제 확인
curl http://localhost:8080/examples/queries | jq
```

## 🔍 코드 하이라이트

### GORM 연결 설정
```go
func NewDatabase(debug bool) (*Database, error) {
    // SQLite 연결
    db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    // 자동 마이그레이션
    db.AutoMigrate(&User{}, &Post{}, &Category{}, &Tag{}, &Comment{})

    return &Database{db}, nil
}
```

### Repository 패턴 구현
```go
type UserRepository struct {
    db *Database
}

func (r *UserRepository) FindByID(id uint) (*User, error) {
    var user User
    // Preload로 연관 데이터 로딩
    err := r.db.Preload("Posts").
           Preload("Comments").
           First(&user, id).Error
    return &user, err
}
```

### 페이지네이션 구현
```go
func (r *UserRepository) FindAll(offset, limit int) ([]User, int64, error) {
    var users []User
    var total int64

    // 전체 개수
    r.db.Model(&User{}).Count(&total)

    // 페이지네이션 적용
    err := r.db.Offset(offset).
           Limit(limit).
           Find(&users).Error

    return users, total, err
}
```

### 관계 설정과 Preload
```go
func (r *PostRepository) FindByID(id uint) (*Post, error) {
    var post Post
    err := r.db.Preload("User").              // 작성자
           Preload("Tags").                   // 태그들
           Preload("Category").               // 카테고리
           Preload("Comments.User").          // 댓글과 댓글 작성자
           First(&post, id).Error

    // 조회수 증가
    r.db.Model(&post).Update("view_count", post.ViewCount+1)

    return &post, err
}
```

### 트랜잭션 처리
```go
func (s *BlogService) CreatePostWithTags(post *Post, tagIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 포스트 생성
        if err := tx.Create(post).Error; err != nil {
            return err
        }

        // 2. 태그 연결
        for _, tagID := range tagIDs {
            var tag Tag
            if err := tx.First(&tag, tagID).Error; err != nil {
                return err
            }
            if err := tx.Model(post).Association("Tags").Append(&tag); err != nil {
                return err
            }
        }

        return nil
    })
}
```

## 🎨 GORM 고급 기능

### 1. **Hooks (생명주기 콜백)**
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 생성 전 실행
    u.Username = strings.ToLower(u.Username)
    return nil
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
    // 생성 후 실행
    log.Printf("New post created: %s", p.Title)
    return nil
}
```

### 2. **Scopes (재사용 가능한 쿼리)**
```go
func Published(db *gorm.DB) *gorm.DB {
    return db.Where("published = ?", true)
}

func Popular(db *gorm.DB) *gorm.DB {
    return db.Where("view_count > ?", 100)
}

// 사용
db.Scopes(Published, Popular).Find(&posts)
```

### 3. **Association Mode**
```go
// Many to Many 관계 처리
db.Model(&post).Association("Tags").Append(&tag1, &tag2)
db.Model(&post).Association("Tags").Delete(&tag1)
db.Model(&post).Association("Tags").Clear()
db.Model(&post).Association("Tags").Count()
```

### 4. **Raw SQL**
```go
// Raw 쿼리 실행
var users []User
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// Exec으로 직접 실행
db.Exec("UPDATE users SET age = age + 1")
```

## 📝 베스트 프랙티스

### 1. **모델 설계 원칙**
```go
// Base 모델 사용
type Base struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 임베딩으로 재사용
type User struct {
    Base
    // ... 필드들
}
```

### 2. **인덱스 설정**
```go
type User struct {
    Email    string `gorm:"uniqueIndex"`
    Username string `gorm:"uniqueIndex;size:50"`
    Age      int    `gorm:"index"`
}
```

### 3. **에러 처리**
```go
if err := db.First(&user, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // 레코드 없음
        return nil, ErrUserNotFound
    }
    // 다른 에러
    return nil, err
}
```

### 4. **성능 최적화**
```go
// N+1 문제 방지
db.Preload("Posts").Find(&users)

// 필요한 필드만 선택
db.Select("id", "email", "name").Find(&users)

// 배치 작업
db.CreateInBatches(users, 100)
```

## 🚀 성능 팁

- **연결 풀 설정**: 프로덕션 환경에서 적절한 연결 풀 설정
- **인덱스 활용**: 자주 조회되는 컬럼에 인덱스 추가
- **Preload 최적화**: 필요한 경우만 Preload 사용
- **캐싱**: 자주 조회되는 데이터 캐싱
- **배치 처리**: 대량 데이터는 배치로 처리

## 📚 추가 학습 자료

- [GORM 공식 문서](https://gorm.io/docs/)
- [SQLite 공식 문서](https://www.sqlite.org/docs.html)
- [Database/SQL 인터페이스](https://golang.org/pkg/database/sql/)
- [SQL 쿼리 최적화](https://use-the-index-luke.com/)

## 🎯 다음 레슨 예고

**Lesson 16: 마이그레이션과 시드 데이터**
- 스키마 버전 관리
- 마이그레이션 전략
- 시드 데이터 생성
- 롤백 처리

GORM으로 데이터베이스 작업이 이렇게 쉬워집니다! 📊