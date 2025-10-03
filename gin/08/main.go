package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 데이터 모델
type User struct {
	ID       int
	Name     string
	Email    string
	IsActive bool
	Role     string
	JoinedAt time.Time
}

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	InStock     bool
	Category    string
	Image       string
}

type PageData struct {
	Title       string
	Message     string
	CurrentYear int
	User        *User
	IsLoggedIn  bool
}

// 커스텀 템플릿 함수들
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatCurrency(price float64) string {
	return fmt.Sprintf("₩%.2f", price)
}

func add(a, b int) int {
	return a + b
}

func isEven(n int) bool {
	return n%2 == 0
}

func main() {
	r := gin.Default()

	// 커스텀 템플릿 함수 등록
	r.SetFuncMap(template.FuncMap{
		"formatDate":     formatDate,
		"formatCurrency": formatCurrency,
		"add":            add,
		"isEven":         isEven,
	})

	// 템플릿 파일 로드
	r.LoadHTMLGlob("08/templates/*")

	// 정적 파일 서빙
	r.Static("/static", "./08/static")

	// ========================================
	// 1. 기본 템플릿 렌더링
	// ========================================

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Gin Template Example",
			"message": "Gin 템플릿 엔진을 사용한 HTML 렌더링",
			"year":    time.Now().Year(),
		})
	})

	// ========================================
	// 2. 데이터와 함께 템플릿 렌더링
	// ========================================

	r.GET("/users", func(c *gin.Context) {
		users := []User{
			{ID: 1, Name: "홍길동", Email: "hong@example.com", IsActive: true, Role: "admin", JoinedAt: time.Now().AddDate(0, -6, 0)},
			{ID: 2, Name: "김철수", Email: "kim@example.com", IsActive: true, Role: "user", JoinedAt: time.Now().AddDate(0, -3, 0)},
			{ID: 3, Name: "이영희", Email: "lee@example.com", IsActive: false, Role: "user", JoinedAt: time.Now().AddDate(0, -1, 0)},
			{ID: 4, Name: "박민수", Email: "park@example.com", IsActive: true, Role: "moderator", JoinedAt: time.Now()},
		}

		c.HTML(http.StatusOK, "users.html", gin.H{
			"title":      "사용자 목록",
			"users":      users,
			"totalUsers": len(users),
			"activeUsers": func() int {
				count := 0
				for _, u := range users {
					if u.IsActive {
						count++
					}
				}
				return count
			}(),
		})
	})

	// ========================================
	// 3. 개별 사용자 페이지
	// ========================================

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 실제로는 DB에서 조회
		user := User{
			ID:       1,
			Name:     "홍길동",
			Email:    "hong@example.com",
			IsActive: true,
			Role:     "admin",
			JoinedAt: time.Now().AddDate(-1, 0, 0),
		}

		c.HTML(http.StatusOK, "user-detail.html", gin.H{
			"title": fmt.Sprintf("%s의 프로필", user.Name),
			"user":  user,
			"stats": gin.H{
				"posts":     42,
				"comments":  128,
				"followers": 256,
			},
		})
	})

	// ========================================
	// 4. 제품 카탈로그
	// ========================================

	r.GET("/products", func(c *gin.Context) {
		products := []Product{
			{ID: 1, Name: "노트북", Price: 1500000, Description: "고성능 게이밍 노트북", InStock: true, Category: "전자제품", Image: "/static/laptop.jpg"},
			{ID: 2, Name: "마우스", Price: 50000, Description: "무선 마우스", InStock: true, Category: "악세서리", Image: "/static/mouse.jpg"},
			{ID: 3, Name: "키보드", Price: 100000, Description: "기계식 키보드", InStock: false, Category: "악세서리", Image: "/static/keyboard.jpg"},
			{ID: 4, Name: "모니터", Price: 400000, Description: "4K 모니터", InStock: true, Category: "전자제품", Image: "/static/monitor.jpg"},
		}

		categories := make(map[string]int)
		for _, p := range products {
			categories[p.Category]++
		}

		c.HTML(http.StatusOK, "products.html", gin.H{
			"title":      "제품 카탈로그",
			"products":   products,
			"categories": categories,
		})
	})

	// ========================================
	// 5. 폼 처리
	// ========================================

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.html", gin.H{
			"title": "문의하기",
		})
	})

	r.POST("/contact", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		message := c.PostForm("message")

		// 실제로는 이메일 전송이나 DB 저장
		c.HTML(http.StatusOK, "contact-success.html", gin.H{
			"title":   "문의 완료",
			"name":    name,
			"email":   email,
			"message": message,
		})
	})

	// ========================================
	// 6. 로그인/로그아웃 (세션 시뮬레이션)
	// ========================================

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "로그인",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 간단한 검증 (실제로는 DB 확인)
		if username == "admin" && password == "1234" {
			// 쿠키 설정 (실제로는 세션 사용)
			c.SetCookie("user", username, 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/dashboard")
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "로그인",
				"error": "잘못된 사용자명 또는 비밀번호입니다",
			})
		}
	})

	r.GET("/dashboard", func(c *gin.Context) {
		// 로그인 확인
		username, err := c.Cookie("user")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":    "대시보드",
			"username": username,
			"stats": gin.H{
				"totalSales":  "₩5,234,000",
				"newOrders":   42,
				"customers":   1234,
				"growthRate":  "+15.3%",
			},
			"recentActivities": []string{
				"새 주문 #1234 접수",
				"고객 문의 답변 완료",
				"재고 업데이트",
				"월간 리포트 생성",
			},
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("user", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	})

	// ========================================
	// 7. 에러 페이지
	// ========================================

	r.GET("/error", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "오류 발생",
			"message": "서버에서 오류가 발생했습니다",
			"code":    500,
		})
	})

	// 404 처리
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "페이지를 찾을 수 없습니다",
			"path":  c.Request.URL.Path,
		})
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("Visit http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}