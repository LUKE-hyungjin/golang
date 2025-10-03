# Lesson 16: 마이그레이션과 시드 데이터 🌱

> 데이터베이스 스키마 버전 관리와 테스트 데이터 생성 완벽 가이드

## 📌 이번 레슨에서 배우는 내용

데이터베이스 스키마는 애플리케이션과 함께 진화합니다. 마이그레이션은 이러한 변경사항을 체계적으로 관리하고, 시드 데이터는 개발과 테스트를 위한 샘플 데이터를 제공합니다. 이번 레슨에서는 프로덕션급 마이그레이션 시스템을 구현합니다.

### 핵심 학습 목표
- ✅ 버전별 마이그레이션 관리
- ✅ Up/Down 마이그레이션
- ✅ 시드 데이터 생성 전략
- ✅ 데이터 Import/Export
- ✅ 롤백 처리
- ✅ Faker를 활용한 더미 데이터 생성

## 🏗 마이그레이션 아키텍처

### 마이그레이션 흐름
```
1. 초기 테이블 생성 (001)
   ↓
2. 필드 추가 (002)
   ↓
3. 새 테이블 추가 (003)
   ↓
4. Soft Delete 추가 (004)
   ↓
5. 메트릭 필드 추가 (005)
```

### 마이그레이션 테이블
```sql
CREATE TABLE migrations (
    id INTEGER PRIMARY KEY,
    version TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    applied_at TIMESTAMP
);
```

## 🛠 구현된 기능

### 1. **마이그레이션 시스템**
- 버전별 마이그레이션 관리
- Up/Down 함수 지원
- 트랜잭션 기반 실행
- 마이그레이션 이력 추적

### 2. **시드 데이터 생성**
- Faker 라이브러리 활용
- 관계형 데이터 생성
- 랜덤 데이터 생성
- JSON Import/Export

### 3. **데이터 관리 도구**
- Clean: 모든 데이터 삭제
- Reset: Clean + Seed
- Export: JSON으로 내보내기
- Import: JSON에서 가져오기

## 🎯 주요 API 엔드포인트

### 마이그레이션 관리
```bash
GET  /migrations/status           # 마이그레이션 상태 조회
POST /migrations/run              # 마이그레이션 실행
POST /migrations/rollback/:version # 특정 버전으로 롤백
```

### 시드 데이터 관리
```bash
POST /seed/run      # 시드 데이터 생성
POST /seed/clean    # 모든 데이터 삭제
POST /seed/reset    # 데이터 리셋 (clean + seed)
POST /seed/export   # JSON으로 내보내기
POST /seed/import   # JSON에서 가져오기
```

### 정보 조회
```bash
GET  /info          # 데이터베이스 통계
GET  /health        # 헬스체크
```

## 💻 실습 가이드

### 1. 설치 및 실행
```bash
# 의존성 설치
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u github.com/go-faker/faker/v4

# 실행 (자동으로 마이그레이션 실행)
cd gin/16
go run main.go
```

### 2. 마이그레이션 관리

#### 마이그레이션 상태 확인
```bash
curl http://localhost:8080/migrations/status | jq

# 응답 예시
{
  "applied_migrations": [
    {
      "id": 1,
      "version": "001_create_users_table",
      "name": "Create users table",
      "applied_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "version": "002_add_user_fields",
      "name": "Add name, bio, avatar fields to users",
      "applied_at": "2024-01-15T10:30:01Z"
    }
  ],
  "total": 5
}
```

#### 마이그레이션 실행
```bash
# 모든 pending 마이그레이션 실행
curl -X POST http://localhost:8080/migrations/run

# 응답
{
  "message": "Migrations completed successfully"
}
```

#### 롤백
```bash
# 특정 버전으로 롤백
curl -X POST http://localhost:8080/migrations/rollback/003_create_posts_categories_tags

# 응답
{
  "message": "Rolled back to 003_create_posts_categories_tags"
}
```

### 3. 시드 데이터 생성

#### 시드 실행
```bash
curl -X POST http://localhost:8080/seed/run

# 응답
{
  "message": "Database seeded successfully"
}

# 로그 출력
🌱 Starting seed process...
✅ Seeded 8 categories
✅ Seeded 15 tags
✅ Seeded users (total: 11)
✅ Seeded posts (total: 35)
✅ Seed process completed!
```

#### 데이터 리셋
```bash
# 모든 데이터 삭제 후 다시 시드
curl -X POST http://localhost:8080/seed/reset

# 응답
{
  "message": "Database reset successfully"
}
```

#### 데이터 정리
```bash
# 모든 데이터 삭제 (테이블은 유지)
curl -X POST http://localhost:8080/seed/clean

# 응답
{
  "message": "Database cleaned successfully"
}
```

### 4. Import/Export

#### 데이터 내보내기
```bash
# 현재 데이터를 JSON 파일로 내보내기
curl -X POST "http://localhost:8080/seed/export?file=backup.json"

# 응답
{
  "message": "Data exported successfully",
  "file": "backup.json"
}

# 파일 확인
cat backup.json | jq '.users[0]'
```

#### 데이터 가져오기
```bash
# JSON 파일에서 데이터 가져오기
curl -X POST "http://localhost:8080/seed/import?file=backup.json"

# 응답
{
  "message": "Data imported successfully",
  "file": "backup.json"
}
```

### 5. 데이터베이스 정보 조회
```bash
curl http://localhost:8080/info | jq

# 응답 예시
{
  "database": "SQLite",
  "file": "blog.db",
  "stats": {
    "users": 11,
    "posts": 35,
    "categories": 8,
    "tags": 15
  }
}
```

## 🔍 코드 하이라이트

### 마이그레이션 정의
```go
type MigrationFunc struct {
    Version string
    Name    string
    Up      func(*gorm.DB) error
    Down    func(*gorm.DB) error
}

// 마이그레이션 예시
{
    Version: "002_add_user_fields",
    Name:    "Add name, bio, avatar fields to users",
    Up: func(db *gorm.DB) error {
        if !db.Migrator().HasColumn(&User{}, "name") {
            db.Migrator().AddColumn(&User{}, "name")
        }
        // ... 다른 필드들
        return nil
    },
    Down: func(db *gorm.DB) error {
        db.Migrator().DropColumn(&User{}, "name")
        // ... 다른 필드들
        return nil
    },
}
```

### 트랜잭션 기반 마이그레이션
```go
func (m *Migrator) Migrate() error {
    for _, migration := range m.migrations {
        // 이미 적용된 마이그레이션 체크
        var count int64
        m.db.Model(&Migration{}).Where("version = ?", migration.Version).Count(&count)

        if count == 0 {
            // 트랜잭션으로 실행
            err := m.db.Transaction(func(tx *gorm.DB) error {
                // 마이그레이션 실행
                if err := migration.Up(tx); err != nil {
                    return err
                }

                // 마이그레이션 기록
                record := Migration{
                    Version:   migration.Version,
                    Name:      migration.Name,
                    AppliedAt: time.Now(),
                }
                return tx.Create(&record).Error
            })

            if err != nil {
                return fmt.Errorf("migration %s failed: %w", migration.Version, err)
            }
        }
    }
    return nil
}
```

### Faker를 사용한 시드 데이터
```go
func (s *Seeder) seedUsers() error {
    for i := 0; i < 10; i++ {
        user := User{
            Email:    faker.Email(),
            Username: faker.Username(),
            Name:     faker.Name(),
            Bio:      faker.Sentence(),
            Avatar:   fmt.Sprintf("https://i.pravatar.cc/150?img=%d", i+1),
            IsActive: true,
        }

        if err := s.db.Create(&user).Error; err != nil {
            log.Printf("Failed to create user: %v", err)
            continue
        }
    }
    return nil
}
```

### 관계형 데이터 시드
```go
func (s *Seeder) seedPosts() error {
    var users []User
    var categories []Category
    var tags []Tag

    s.db.Find(&users)
    s.db.Find(&categories)
    s.db.Find(&tags)

    for _, user := range users {
        for i := 0; i < rand.Intn(5)+1; i++ {
            post := Post{
                Title:      faker.Sentence(),
                Content:    faker.Paragraph(),
                UserID:     user.ID,
                CategoryID: &categories[rand.Intn(len(categories))].ID,
            }

            s.db.Create(&post)

            // 랜덤 태그 추가
            numTags := rand.Intn(3) + 1
            selectedTags := tags[:numTags]
            s.db.Model(&post).Association("Tags").Append(selectedTags)
        }
    }
    return nil
}
```

## 🎨 마이그레이션 전략

### 1. **버전 네이밍 규칙**
```
001_create_users_table       # 테이블 생성
002_add_user_fields          # 필드 추가
003_rename_column            # 컬럼 이름 변경
004_add_index               # 인덱스 추가
005_drop_unused_table       # 테이블 삭제
```

### 2. **안전한 마이그레이션**
```go
// 컬럼 존재 여부 체크
if !db.Migrator().HasColumn(&User{}, "new_field") {
    db.Migrator().AddColumn(&User{}, "new_field")
}

// 인덱스 존재 여부 체크
if !db.Migrator().HasIndex(&User{}, "idx_email") {
    db.Migrator().CreateIndex(&User{}, "Email")
}
```

### 3. **데이터 마이그레이션**
```go
// 스키마 변경과 데이터 변환
Up: func(db *gorm.DB) error {
    // 1. 새 컬럼 추가
    db.Migrator().AddColumn(&User{}, "full_name")

    // 2. 기존 데이터 변환
    var users []User
    db.Find(&users)
    for _, user := range users {
        user.FullName = user.FirstName + " " + user.LastName
        db.Save(&user)
    }

    // 3. 기존 컬럼 삭제
    db.Migrator().DropColumn(&User{}, "first_name")
    db.Migrator().DropColumn(&User{}, "last_name")

    return nil
}
```

## 📝 베스트 프랙티스

### 1. **마이그레이션 원칙**
- 항상 Up과 Down 함수 작성
- 트랜잭션으로 실행
- 멱등성 보장 (여러 번 실행해도 안전)
- 프로덕션 데이터 백업

### 2. **시드 데이터 관리**
- 개발/테스트 환경에서만 사용
- 현실적인 데이터 생성
- 관계 무결성 유지
- 성능 테스트용 대량 데이터 옵션

### 3. **롤백 전략**
- 데이터 손실 최소화
- 단계별 롤백 지원
- 롤백 전 백업 필수
- 테스트 환경에서 검증

## 🚀 프로덕션 체크리스트

- [ ] 모든 마이그레이션이 테스트되었는가?
- [ ] 롤백 스크립트가 준비되었는가?
- [ ] 데이터베이스 백업이 되었는가?
- [ ] 다운타임이 필요한가?
- [ ] 마이그레이션 순서가 올바른가?
- [ ] 대용량 데이터에서 테스트했는가?
- [ ] 인덱스가 적절히 설정되었는가?

## 📚 추가 학습 자료

- [GORM Migration Guide](https://gorm.io/docs/migration.html)
- [Database Migration Best Practices](https://www.prisma.io/dataguide/types/relational/migration-strategies)
- [Faker Documentation](https://github.com/go-faker/faker)
- [Schema Evolution](https://martinfowler.com/articles/evodb.html)

## 🎯 다음 레슨 예고

**Lesson 17: 트랜잭션과 컨텍스트 타임아웃**
- ACID 트랜잭션 처리
- 동시성 제어
- 컨텍스트 기반 타임아웃
- 데드락 방지

마이그레이션으로 데이터베이스를 안전하게 진화시키세요! 🌱