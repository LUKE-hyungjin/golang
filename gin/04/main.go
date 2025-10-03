package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User 구조체
type User struct {
	ID       string    `json:"id" binding:"required"`
	Username string    `json:"username" binding:"required,min=3,max=20"`
	Email    string    `json:"email" binding:"required,email"`
	Age      int       `json:"age" binding:"min=1,max=120"`
	Role     string    `json:"role" binding:"required,oneof=admin user guest"`
	Created  time.Time `json:"created"`
}

// LoginRequest 구조체
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func main() {
	r := gin.Default()

	// 1. Context의 기본 메서드들
	r.GET("/context/basic/:name", func(c *gin.Context) {
		// Path 파라미터 가져오기
		name := c.Param("name")

		// Query 파라미터 가져오기
		age := c.Query("age")
		city := c.DefaultQuery("city", "Seoul")

		// Query 파라미터 존재 확인
		page, exists := c.GetQuery("page")
		if !exists {
			page = "1"
		}

		// 헤더 가져오기
		userAgent := c.GetHeader("User-Agent")
		contentType := c.Request.Header.Get("Content-Type")

		c.JSON(http.StatusOK, gin.H{
			"name":         name,
			"age":          age,
			"city":         city,
			"page":         page,
			"user_agent":   userAgent,
			"content_type": contentType,
		})
	})

	// 2. Request Body 바인딩
	r.POST("/context/bind", func(c *gin.Context) {
		var user User

		// JSON 바인딩 시도
		if err := c.ShouldBindJSON(&user); err != nil {
			// 바인딩 에러 시 BadRequest 응답
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
			return
		}

		// 타임스탬프 추가
		user.Created = time.Now()

		// 성공 응답
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	})

	// 3. 다양한 응답 포맷
	r.GET("/context/response/:format", func(c *gin.Context) {
		format := c.Param("format")
		data := gin.H{
			"message": "Hello from Gin",
			"time":    time.Now().Format(time.RFC3339),
		}

		switch format {
		case "json":
			// JSON 응답
			c.JSON(http.StatusOK, data)
		case "xml":
			// XML 응답
			c.XML(http.StatusOK, data)
		case "yaml":
			// YAML 응답
			c.YAML(http.StatusOK, data)
		case "string":
			// 문자열 응답
			c.String(http.StatusOK, "Message: %s, Time: %s",
				data["message"], data["time"])
		case "data":
			// 바이너리 데이터 응답
			c.Data(http.StatusOK, "text/plain",
				[]byte(fmt.Sprintf("Raw data: %v", data)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported format. Use: json, xml, yaml, string, or data",
			})
		}
	})

	// 4. Context에 값 저장하고 가져오기
	r.POST("/context/login", func(c *gin.Context) {
		var login LoginRequest
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 간단한 인증 체크 (실제로는 DB 조회 필요)
		if login.Username == "admin" && login.Password == "password123" {
			// Context에 사용자 정보 저장
			c.Set("username", login.Username)
			c.Set("role", "admin")
			c.Set("authenticated", true)

			// 다음 핸들러로 진행
			processAfterLogin(c)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
		}
	})

	// 5. 파일 업로드/다운로드
	r.POST("/context/upload", func(c *gin.Context) {
		// 단일 파일 받기
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file uploaded",
			})
			return
		}
		defer file.Close()

		// 파일 정보
		c.JSON(http.StatusOK, gin.H{
			"filename": header.Filename,
			"size":     header.Size,
			"headers":  header.Header,
		})
	})

	// 파일 다운로드
	r.GET("/context/download", func(c *gin.Context) {
		// 파일 경로 (실제로는 동적으로 결정)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=example.txt")
		c.Header("Content-Type", "application/octet-stream")
		c.String(http.StatusOK, "This is the file content")
	})

	// 6. 쿠키 처리
	r.GET("/context/cookie/set", func(c *gin.Context) {
		// 쿠키 설정
		c.SetCookie(
			"session_id",     // name
			"abc123xyz",      // value
			3600,             // maxAge (초)
			"/",              // path
			"localhost",      // domain
			false,            // secure
			true,             // httpOnly
		)
		c.JSON(http.StatusOK, gin.H{
			"message": "Cookie set successfully",
		})
	})

	r.GET("/context/cookie/get", func(c *gin.Context) {
		// 쿠키 읽기
		cookie, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cookie not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"session_id": cookie,
		})
	})

	// 7. Request 정보 가져오기
	r.GET("/context/request-info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method":      c.Request.Method,
			"url":         c.Request.URL.String(),
			"proto":       c.Request.Proto,
			"remote_addr": c.Request.RemoteAddr,
			"client_ip":   c.ClientIP(),
			"host":        c.Request.Host,
			"referer":     c.Request.Referer(),
			"user_agent":  c.Request.UserAgent(),
			"headers":     c.Request.Header,
		})
	})

	// 8. 스트림 응답
	r.GET("/context/stream", func(c *gin.Context) {
		chanStream := make(chan int, 10)
		go func() {
			defer close(chanStream)
			for i := 0; i < 5; i++ {
				chanStream <- i
				time.Sleep(time.Second)
			}
		}()

		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// 9. Context 복사 (고루틴에서 사용)
	r.GET("/context/async", func(c *gin.Context) {
		// Context 복사 (고루틴에서 안전하게 사용)
		cCp := c.Copy()
		go func() {
			// 비동기 작업 시뮬레이션
			time.Sleep(2 * time.Second)
			fmt.Printf("Async request from: %s\n", cCp.ClientIP())
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": "Request is being processed asynchronously",
		})
	})

	// 10. Abort와 에러 처리
	r.GET("/context/abort", func(c *gin.Context) {
		// 조건 체크
		token := c.GetHeader("Authorization")
		if token == "" {
			// 요청 중단하고 에러 응답
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return // 이 이후 코드는 실행되지 않음
		}

		// 정상 처리
		c.JSON(http.StatusOK, gin.H{
			"message": "Authorized",
			"token":   token,
		})
	})

	// 11. Redirect
	r.GET("/context/redirect", func(c *gin.Context) {
		// HTTP 리다이렉트
		c.Redirect(http.StatusMovedPermanently, "https://gin-gonic.com")
	})

	// 내부 리다이렉트 (라우터 내에서)
	r.GET("/old-endpoint", func(c *gin.Context) {
		// 새 엔드포인트로 요청 전달
		c.Request.URL.Path = "/new-endpoint"
		r.HandleContext(c)
	})

	r.GET("/new-endpoint", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is the new endpoint",
		})
	})

	// 12. Content Negotiation
	r.GET("/context/negotiate", func(c *gin.Context) {
		// Accept 헤더에 따라 다른 형식으로 응답
		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered: []string{gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML},
			HTMLData: gin.H{
				"message": "This is HTML response",
			},
			JSONData: gin.H{
				"message": "This is JSON response",
			},
			XMLData: gin.H{
				"message": "This is XML response",
			},
			YAMLData: gin.H{
				"message": "This is YAML response",
			},
		})
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

// Context에서 값 가져와서 처리하는 함수
func processAfterLogin(c *gin.Context) {
	// Context에서 값 가져오기
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Username not found in context",
		})
		return
	}

	role, _ := c.Get("role")
	authenticated, _ := c.Get("authenticated")

	// MustGet은 값이 없으면 panic
	// username := c.MustGet("username").(string)

	// 타입 assertion과 함께 안전하게 가져오기
	if auth, ok := authenticated.(bool); ok && auth {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"username":      username,
			"role":          role,
			"authenticated": auth,
		})
	}
}