package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 구조체들
type HealthResponse struct {
	Status  string    `json:"status"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Version   string    `json:"api_version"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
	Version     string  `json:"api_version"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func main() {
	r := gin.Default()

	// 루트 엔드포인트
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Server with Route Groups and Versioning",
			"versions": []string{"v1", "v2"},
			"documentation": "/docs",
		})
	})

	// ========================================
	// 1. API 버저닝 - URL Path 방식
	// ========================================

	// API v1 그룹
	v1 := r.Group("/api/v1")
	{
		// 헬스체크
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, HealthResponse{
				Status:  "healthy",
				Version: "1.0",
				Time:    time.Now(),
			})
		})

		// 사용자 관련 라우트 그룹
		v1Users := v1.Group("/users")
		{
			v1Users.GET("", getUsersV1)
			v1Users.GET("/:id", getUserV1)
			v1Users.POST("", createUserV1)
			v1Users.PUT("/:id", updateUserV1)
			v1Users.DELETE("/:id", deleteUserV1)

			// 중첩 그룹 - 사용자 프로필
			v1Profile := v1Users.Group("/:id/profile")
			{
				v1Profile.GET("", getUserProfileV1)
				v1Profile.PUT("", updateUserProfileV1)
			}

			// 중첩 그룹 - 사용자 설정
			v1Settings := v1Users.Group("/:id/settings")
			{
				v1Settings.GET("", getUserSettings)
				v1Settings.PUT("", updateUserSettings)
			}
		}

		// 제품 관련 라우트 그룹
		v1Products := v1.Group("/products")
		{
			v1Products.GET("", getProductsV1)
			v1Products.GET("/:id", getProductV1)
		}
	}

	// API v2 그룹 (개선된 버전)
	v2 := r.Group("/api/v2")
	// v2 전용 미들웨어
	v2.Use(v2Middleware())
	{
		// 헬스체크
		v2.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, HealthResponse{
				Status:  "healthy",
				Version: "2.0",
				Time:    time.Now(),
			})
		})

		// v2 사용자 라우트 (개선된 응답 형식)
		v2Users := v2.Group("/users")
		{
			v2Users.GET("", getUsersV2)
			v2Users.GET("/:id", getUserV2)
			v2Users.POST("", createUserV2)

			// v2에서 추가된 기능
			v2Users.GET("/:id/activities", getUserActivities)
			v2Users.POST("/:id/follow", followUser)
		}

		// v2 제품 라우트 (필터링 기능 추가)
		v2Products := v2.Group("/products")
		{
			v2Products.GET("", getProductsV2)
			v2Products.GET("/search", searchProducts)
			v2Products.GET("/:id", getProductV2)
			v2Products.GET("/:id/reviews", getProductReviews)
		}
	}

	// ========================================
	// 2. 관리자 패널 라우트 그룹
	// ========================================
	admin := r.Group("/admin")
	admin.Use(adminAuthMiddleware())
	{
		// 대시보드
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin Dashboard",
				"stats": gin.H{
					"total_users":    100,
					"total_products": 50,
					"total_orders":   250,
				},
			})
		})

		// 사용자 관리
		adminUsers := admin.Group("/users")
		{
			adminUsers.GET("", getAllUsers)
			adminUsers.PUT("/:id/ban", banUser)
			adminUsers.PUT("/:id/unban", unbanUser)
			adminUsers.DELETE("/:id", deleteUser)
		}

		// 시스템 관리
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/logs", getSystemLogs)
			adminSystem.GET("/metrics", getSystemMetrics)
			adminSystem.POST("/maintenance", toggleMaintenance)
		}
	}

	// ========================================
	// 3. Public API (인증 불필요)
	// ========================================
	public := r.Group("/public")
	{
		public.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service": "Gin API Server",
				"status":  "operational",
				"time":    time.Now(),
			})
		})

		public.GET("/docs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"swagger": "/public/swagger",
				"postman": "/public/postman",
			})
		})
	}

	// ========================================
	// 4. Internal API (내부 서비스용)
	// ========================================
	internal := r.Group("/internal")
	internal.Use(internalAuthMiddleware())
	{
		internal.GET("/health/detailed", detailedHealthCheck)
		internal.POST("/cache/clear", clearCache)
		internal.POST("/jobs/trigger", triggerJob)
	}

	// ========================================
	// 5. Webhook 엔드포인트 그룹
	// ========================================
	webhooks := r.Group("/webhooks")
	{
		// 각 서비스별 webhook
		webhooks.POST("/github", handleGithubWebhook)
		webhooks.POST("/stripe", handleStripeWebhook)
		webhooks.POST("/slack", handleSlackWebhook)
	}

	// ========================================
	// 6. 헤더 기반 버저닝 예제
	// ========================================
	r.GET("/api/users", versionedHandler())

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("Available API versions: v1, v2")
	fmt.Println("Admin panel: /admin")
	fmt.Println("Public API: /public")
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

// ========================================
// V1 핸들러들
// ========================================

func getUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"users": []UserResponse{
			{
				ID:        "1",
				Username:  "user1",
				Email:     "user1@example.com",
				CreatedAt: time.Now().AddDate(0, -1, 0),
				Version:   "v1",
			},
		},
	})
}

func getUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, UserResponse{
		ID:        id,
		Username:  "user" + id,
		Email:     fmt.Sprintf("user%s@example.com", id),
		CreatedAt: time.Now(),
		Version:   "v1",
	})
}

func createUserV1(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "Invalid input",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"version": "v1",
		"message": "User created",
		"id":      "new-user-id",
	})
}

func updateUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "User updated",
		"id":      id,
	})
}

func deleteUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "User deleted",
		"id":      id,
	})
}

func getUserProfileV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"user_id": id,
		"profile": gin.H{
			"bio":      "User bio",
			"avatar":   "/avatars/default.png",
			"location": "Seoul, Korea",
		},
	})
}

func updateUserProfileV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "Profile updated",
		"user_id": id,
	})
}

func getProductsV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"products": []ProductResponse{
			{
				ID:       "1",
				Name:     "Product 1",
				Price:    99.99,
				Category: "Electronics",
				InStock:  true,
				Version:  "v1",
			},
		},
	})
}

func getProductV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       id,
		Name:     "Product " + id,
		Price:    99.99,
		Category: "Electronics",
		InStock:  true,
		Version:  "v1",
	})
}

// ========================================
// V2 핸들러들 (개선된 버전)
// ========================================

func getUsersV2(c *gin.Context) {
	// v2에서는 페이지네이션 지원
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"data": []UserResponse{
			{
				ID:        "1",
				Username:  "enhanced_user1",
				Email:     "user1@example.com",
				CreatedAt: time.Now(),
				Version:   "v2",
			},
		},
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      100,
			"total_pages": 10,
		},
	})
}

func getUserV2(c *gin.Context) {
	id := c.Param("id")
	// v2에서는 더 많은 정보 제공
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"data": gin.H{
			"id":         id,
			"username":   "user" + id,
			"email":      fmt.Sprintf("user%s@example.com", id),
			"created_at": time.Now(),
			"profile": gin.H{
				"bio":      "Enhanced bio",
				"avatar":   "/avatars/user.png",
				"verified": true,
			},
			"stats": gin.H{
				"posts":     42,
				"followers": 100,
				"following": 50,
			},
		},
	})
}

func createUserV2(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"version": "v2",
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "Validation failed",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"version": "v2",
		"data": gin.H{
			"id":         "new-user-id",
			"username":   input["username"],
			"created_at": time.Now(),
		},
		"message": "User successfully created",
	})
}

func getUserActivities(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"user_id": id,
		"activities": []gin.H{
			{
				"type":      "post_created",
				"timestamp": time.Now().Add(-1 * time.Hour),
				"details":   "Created a new post",
			},
			{
				"type":      "comment_added",
				"timestamp": time.Now().Add(-2 * time.Hour),
				"details":   "Commented on a post",
			},
		},
	})
}

func followUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"message": "Successfully followed user",
		"user_id": id,
	})
}

func getProductsV2(c *gin.Context) {
	// v2에서는 필터링 지원
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"filters": gin.H{
			"category":  category,
			"min_price": minPrice,
			"max_price": maxPrice,
		},
		"data": []ProductResponse{
			{
				ID:       "1",
				Name:     "Enhanced Product 1",
				Price:    99.99,
				Category: "Electronics",
				InStock:  true,
				Version:  "v2",
			},
		},
		"meta": gin.H{
			"total_count": 50,
			"filtered_count": 10,
		},
	})
}

func searchProducts(c *gin.Context) {
	query := c.Query("q")
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"query":   query,
		"results": []gin.H{
			{
				"id":       "1",
				"name":     "Product matching " + query,
				"score":    0.95,
				"category": "Electronics",
			},
		},
	})
}

func getProductV2(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"data": gin.H{
			"id":       id,
			"name":     "Enhanced Product " + id,
			"price":    99.99,
			"category": "Electronics",
			"in_stock": true,
			"images": []string{
				"/images/product1.jpg",
				"/images/product2.jpg",
			},
			"specifications": gin.H{
				"weight": "500g",
				"dimensions": "10x10x5 cm",
			},
		},
	})
}

func getProductReviews(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"version":    "v2",
		"product_id": id,
		"reviews": []gin.H{
			{
				"id":     "1",
				"rating": 5,
				"comment": "Great product!",
				"author": "user1",
				"date":   time.Now().AddDate(0, 0, -7),
			},
		},
		"average_rating": 4.5,
		"total_reviews":  42,
	})
}

// ========================================
// Admin 핸들러들
// ========================================

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"admin": true,
		"users": []gin.H{
			{"id": "1", "username": "user1", "status": "active"},
			{"id": "2", "username": "user2", "status": "banned"},
		},
	})
}

func banUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User banned",
		"user_id": id,
	})
}

func unbanUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User unbanned",
		"user_id": id,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNoContent, nil)
}

func getSystemLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"logs": []gin.H{
			{"timestamp": time.Now(), "level": "INFO", "message": "System started"},
			{"timestamp": time.Now().Add(-1 * time.Hour), "level": "WARN", "message": "High CPU usage"},
		},
	})
}

func getSystemMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"cpu_usage":    "45%",
		"memory_usage": "2.5GB",
		"disk_usage":   "120GB",
		"uptime":       "7 days",
	})
}

func toggleMaintenance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":           "Maintenance mode toggled",
		"maintenance_mode": true,
	})
}

// ========================================
// 기타 핸들러들
// ========================================

func getUserSettings(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": id,
		"settings": gin.H{
			"theme":         "dark",
			"language":      "en",
			"notifications": true,
		},
	})
}

func updateUserSettings(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Settings updated",
		"user_id": id,
	})
}

func detailedHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"checks": gin.H{
			"database":    "connected",
			"cache":       "connected",
			"message_queue": "connected",
		},
		"metrics": gin.H{
			"response_time": "15ms",
			"error_rate":    "0.01%",
		},
	})
}

func clearCache(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Cache cleared successfully",
		"cleared_keys": 1024,
	})
}

func triggerJob(c *gin.Context) {
	jobType := c.Query("type")
	c.JSON(http.StatusOK, gin.H{
		"message": "Job triggered",
		"job_id":  "job-123",
		"type":    jobType,
	})
}

func handleGithubWebhook(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	c.JSON(http.StatusOK, gin.H{
		"message": "GitHub webhook received",
		"event":   event,
	})
}

func handleStripeWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Stripe webhook received",
	})
}

func handleSlackWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Slack webhook received",
	})
}

// ========================================
// 미들웨어들
// ========================================

func v2Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", "2.0")
		c.Next()
	}
}

func adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		if token != "admin-secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Admin authentication required",
			})
			return
		}
		c.Next()
	}
}

func internalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Internal-API-Key")
		if apiKey != "internal-api-key-123" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Internal API key required",
			})
			return
		}
		c.Next()
	}
}

// 헤더 기반 버저닝
func versionedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader("API-Version")

		switch version {
		case "1.0":
			getUsersV1(c)
		case "2.0":
			getUsersV2(c)
		default:
			// 기본값은 최신 버전
			getUsersV2(c)
		}
	}
}