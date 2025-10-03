package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 사용자 정보 구조체
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func main() {
	// 기본 미들웨어 없이 시작 (수동 설정)
	r := gin.New()

	// 1. 전역 미들웨어 - 모든 라우트에 적용
	r.Use(LoggerMiddleware())      // 커스텀 로거
	r.Use(gin.Recovery())           // 패닉 복구
	r.Use(CORSMiddleware())         // CORS 설정
	r.Use(RequestIDMiddleware())    // Request ID 추가

	// 공개 엔드포인트 (인증 불필요)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Welcome to Gin Middleware Example",
			"request_id": c.GetString("RequestID"),
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 2. 그룹 미들웨어 - 특정 그룹에만 적용
	// API v1 그룹
	v1 := r.Group("/api/v1")
	v1.Use(RateLimitMiddleware(10)) // 분당 10 요청 제한
	{
		// 공개 API
		v1.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Public API v1",
				"version": "1.0",
			})
		})

		// 인증이 필요한 API 그룹
		protected := v1.Group("/protected")
		protected.Use(AuthMiddleware()) // 인증 미들웨어
		{
			protected.GET("/profile", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(http.StatusOK, gin.H{
					"message": "Protected resource",
					"user":    user,
				})
			})

			// 관리자만 접근 가능
			admin := protected.Group("/admin")
			admin.Use(RequireRole("admin")) // 역할 체크 미들웨어
			{
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Admin only resource",
						"users": []User{
							{ID: "1", Username: "admin", Role: "admin"},
							{ID: "2", Username: "user1", Role: "user"},
						},
					})
				})

				admin.DELETE("/users/:id", func(c *gin.Context) {
					id := c.Param("id")
					c.JSON(http.StatusOK, gin.H{
						"message": "User deleted",
						"id":      id,
					})
				})
			}
		}
	}

	// 3. 라우트별 미들웨어 - 특정 라우트에만 적용
	r.GET("/slow", TimeoutMiddleware(2*time.Second), func(c *gin.Context) {
		// 느린 작업 시뮬레이션
		time.Sleep(1 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "Slow operation completed",
		})
	})

	// 4. 미들웨어 체인과 next() 흐름 데모
	r.GET("/middleware-chain",
		FirstMiddleware(),
		SecondMiddleware(),
		ThirdMiddleware(),
		func(c *gin.Context) {
			log.Println("4. Main Handler")
			c.JSON(http.StatusOK, gin.H{
				"message": "Handler executed",
				"flow":    c.GetStringSlice("flow"),
			})
			log.Println("5. Main Handler - After response")
		})

	// 5. 조건부 미들웨어
	r.GET("/conditional", ConditionalMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Conditional middleware passed",
		})
	})

	// 6. 에러 처리 미들웨어 데모
	r.GET("/error", ErrorHandlingMiddleware(), func(c *gin.Context) {
		// 의도적으로 에러 발생
		c.Error(fmt.Errorf("something went wrong"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})

	// 7. 데이터 변환 미들웨어
	r.POST("/transform", DataTransformMiddleware(), func(c *gin.Context) {
		data, _ := c.Get("transformedData")
		c.JSON(http.StatusOK, gin.H{
			"original":    c.GetString("originalData"),
			"transformed": data,
		})
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

// 1. 커스텀 로거 미들웨어
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log 작성
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v",
			clientIP,
			method,
			path,
			statusCode,
			latency,
		)
	}
}

// 2. CORS 미들웨어
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 3. Request ID 미들웨어
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// 4. Rate Limit 미들웨어 (간단한 구현)
func RateLimitMiddleware(limit int) gin.HandlerFunc {
	// 실제로는 Redis나 in-memory store 사용
	requests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// 1분 이내의 요청만 카운트
		windowStart := now.Add(-1 * time.Minute)

		// 이전 요청 기록 필터링
		var validRequests []time.Time
		for _, reqTime := range requests[clientIP] {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": "60 seconds",
			})
			return
		}

		// 현재 요청 기록
		requests[clientIP] = append(validRequests, now)
		c.Next()
	}
}

// 5. 인증 미들웨어
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Bearer token 체크
		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header",
			})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// 토큰 검증 (실제로는 JWT 검증 등)
		if token != "valid-token-123" && token != "admin-token-456" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// 사용자 정보 설정
		user := User{
			ID:       "1",
			Username: "testuser",
			Role:     "user",
		}

		if token == "admin-token-456" {
			user.Role = "admin"
			user.Username = "admin"
		}

		c.Set("user", user)
		c.Set("authenticated", true)
		c.Next()
	}
}

// 6. 역할 체크 미들웨어
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		user, ok := userInterface.(User)
		if !ok || user.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Requires %s role", requiredRole),
			})
			return
		}

		c.Next()
	}
}

// 7. 타임아웃 미들웨어
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 타임아웃 컨텍스트 생성
		ctx, cancel := c.Request.Context(), func() {}
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})

		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			// 정상 완료
		case <-time.After(timeout):
			// 타임아웃
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
		}
	}
}

// 8. 미들웨어 체인 데모
func FirstMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("1. First Middleware - Before")
		c.Set("flow", []string{"first-before"})
		c.Next()
		log.Println("6. First Middleware - After")
		flow := append(c.GetStringSlice("flow"), "first-after")
		c.Set("flow", flow)
	}
}

func SecondMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("2. Second Middleware - Before")
		flow := append(c.GetStringSlice("flow"), "second-before")
		c.Set("flow", flow)
		c.Next()
		log.Println("7. Second Middleware - After")
		flow = append(c.GetStringSlice("flow"), "second-after")
		c.Set("flow", flow)
	}
}

func ThirdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("3. Third Middleware - Before")
		flow := append(c.GetStringSlice("flow"), "third-before")
		c.Set("flow", flow)
		c.Next()
		log.Println("8. Third Middleware - After")
		flow = append(c.GetStringSlice("flow"), "third-after")
		c.Set("flow", flow)
	}
}

// 9. 조건부 미들웨어
func ConditionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 특정 조건 체크
		skipAuth := c.Query("skip_auth") == "true"

		if skipAuth {
			log.Println("Skipping authentication")
			c.Next()
			return
		}

		// 정상적인 인증 프로세스
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required (use ?skip_auth=true to bypass)",
			})
			return
		}

		c.Next()
	}
}

// 10. 에러 처리 미들웨어
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 에러가 있는지 확인
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error occurred: %v", err)

			// 에러 로깅, 알림 등 처리
			// 실제로는 Sentry, DataDog 등에 전송
		}
	}
}

// 11. 데이터 변환 미들웨어
func DataTransformMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err == nil {
			// 원본 저장
			c.Set("originalData", data)

			// 데이터 변환 (예: 모든 문자열을 대문자로)
			transformed := make(map[string]interface{})
			for k, v := range data {
				if str, ok := v.(string); ok {
					transformed[k] = strings.ToUpper(str)
				} else {
					transformed[k] = v
				}
			}
			c.Set("transformedData", transformed)
		}

		c.Next()
	}
}