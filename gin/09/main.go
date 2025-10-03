package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ========================================
// 에러 응답 구조체 정의
// ========================================

// StandardError - 표준 에러 응답 구조
type StandardError struct {
	Code      int         `json:"code"`               // HTTP 상태 코드
	Message   string      `json:"message"`            // 사용자에게 보여줄 메시지
	ErrorCode string      `json:"error_code"`         // 내부 에러 코드
	Details   interface{} `json:"details,omitempty"`  // 상세 정보 (옵셔널)
	Timestamp time.Time   `json:"timestamp"`          // 에러 발생 시간
	Path      string      `json:"path"`               // 요청 경로
	RequestID string      `json:"request_id"`         // 요청 추적 ID
}

// ValidationError - 입력 검증 에러
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ErrorResponse - API 에러 응답 래퍼
type ErrorResponse struct {
	Success bool           `json:"success"`
	Error   *StandardError `json:"error"`
}

// SuccessResponse - 성공 응답 래퍼
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ========================================
// 커스텀 에러 타입들
// ========================================

// BusinessError - 비즈니스 로직 에러
type BusinessError struct {
	Code    string
	Message string
	Status  int
}

func (e BusinessError) Error() string {
	return e.Message
}

// 에러 코드 상수
const (
	// 클라이언트 에러 (4xx)
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrMethodNotAllowed    = "METHOD_NOT_ALLOWED"
	ErrConflict            = "CONFLICT"
	ErrValidation          = "VALIDATION_ERROR"
	ErrTooManyRequests     = "TOO_MANY_REQUESTS"

	// 서버 에러 (5xx)
	ErrInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
	ErrDatabaseConnection  = "DATABASE_ERROR"
	ErrExternalService     = "EXTERNAL_SERVICE_ERROR"
)

// ========================================
// 에러 응답 헬퍼 함수들
// ========================================

// NewErrorResponse - 에러 응답 생성
func NewErrorResponse(c *gin.Context, status int, code string, message string, details interface{}) {
	requestID, _ := c.Get("RequestID")

	errorResp := ErrorResponse{
		Success: false,
		Error: &StandardError{
			Code:      status,
			Message:   message,
			ErrorCode: code,
			Details:   details,
			Timestamp: time.Now(),
			Path:      c.Request.URL.Path,
			RequestID: fmt.Sprintf("%v", requestID),
		},
	}

	c.JSON(status, errorResp)
}

// NewSuccessResponse - 성공 응답 생성
func NewSuccessResponse(c *gin.Context, status int, data interface{}, meta interface{}) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	c.JSON(status, response)
}

// ========================================
// 상태 코드별 헬퍼 함수들
// ========================================

// BadRequest - 400
func BadRequest(c *gin.Context, message string, details interface{}) {
	NewErrorResponse(c, http.StatusBadRequest, ErrBadRequest, message, details)
}

// Unauthorized - 401
func Unauthorized(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusUnauthorized, ErrUnauthorized, message, nil)
}

// Forbidden - 403
func Forbidden(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusForbidden, ErrForbidden, message, nil)
}

// NotFound - 404
func NotFound(c *gin.Context, resource string) {
	message := fmt.Sprintf("%s not found", resource)
	NewErrorResponse(c, http.StatusNotFound, ErrNotFound, message, nil)
}

// Conflict - 409
func Conflict(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusConflict, ErrConflict, message, nil)
}

// InternalServerError - 500
func InternalServerError(c *gin.Context, message string) {
	NewErrorResponse(c, http.StatusInternalServerError, ErrInternalServer, message, nil)
}

// ValidationFailed - 422
func ValidationFailed(c *gin.Context, errors []ValidationError) {
	NewErrorResponse(c, http.StatusUnprocessableEntity, ErrValidation, "Validation failed", errors)
}

func main() {
	r := gin.Default()

	// Request ID 미들웨어
	r.Use(func(c *gin.Context) {
		c.Set("RequestID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
		c.Next()
	})

	// ========================================
	// 1. 정상 응답 예제 (2xx)
	// ========================================

	// 200 OK - 성공적인 GET 요청
	r.GET("/api/users", func(c *gin.Context) {
		users := []gin.H{
			{"id": 1, "name": "John", "email": "john@example.com"},
			{"id": 2, "name": "Jane", "email": "jane@example.com"},
		}

		NewSuccessResponse(c, http.StatusOK, users, gin.H{
			"total": 2,
			"page":  1,
		})
	})

	// 201 Created - 리소스 생성 성공
	r.POST("/api/users", func(c *gin.Context) {
		var user map[string]interface{}

		if err := c.ShouldBindJSON(&user); err != nil {
			BadRequest(c, "Invalid JSON format", err.Error())
			return
		}

		// 이메일 중복 체크 시뮬레이션
		if user["email"] == "duplicate@example.com" {
			Conflict(c, "Email already exists")
			return
		}

		// 사용자 생성 성공
		user["id"] = 123
		user["created_at"] = time.Now()

		NewSuccessResponse(c, http.StatusCreated, user, nil)
	})

	// 204 No Content - 성공했지만 응답 본문 없음
	r.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// ID가 999면 없는 것으로 처리
		if id == "999" {
			NotFound(c, "User")
			return
		}

		// 삭제 성공
		c.Status(http.StatusNoContent)
	})

	// ========================================
	// 2. 클라이언트 에러 (4xx)
	// ========================================

	// 400 Bad Request - 잘못된 요청
	r.GET("/api/bad-request", func(c *gin.Context) {
		BadRequest(c, "Missing required parameters", gin.H{
			"required": []string{"name", "email"},
			"provided": []string{},
		})
	})

	// 401 Unauthorized - 인증 필요
	r.GET("/api/protected", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			Unauthorized(c, "Authentication required")
			return
		}

		if token != "Bearer valid-token" {
			Unauthorized(c, "Invalid or expired token")
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"message": "Access granted",
		}, nil)
	})

	// 403 Forbidden - 권한 없음
	r.DELETE("/api/admin/users", func(c *gin.Context) {
		// 관리자가 아닌 경우
		Forbidden(c, "Admin access required")
	})

	// 404 Not Found - 리소스 없음
	r.GET("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		if id == "999" {
			NotFound(c, "User")
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"id":   id,
			"name": "John Doe",
		}, nil)
	})

	// 405 Method Not Allowed
	r.GET("/api/method-not-allowed", func(c *gin.Context) {
		NewErrorResponse(c, http.StatusMethodNotAllowed, ErrMethodNotAllowed,
			"Method not allowed", gin.H{
				"allowed_methods": []string{"POST", "PUT"},
			})
	})

	// 409 Conflict - 충돌
	r.POST("/api/conflict", func(c *gin.Context) {
		Conflict(c, "Resource already exists with the same identifier")
	})

	// 422 Unprocessable Entity - 검증 실패
	r.POST("/api/validate", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Age      int    `json:"age"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			BadRequest(c, "Invalid JSON", err.Error())
			return
		}

		var errors []ValidationError

		// 이메일 검증
		if input.Email == "" {
			errors = append(errors, ValidationError{
				Field:   "email",
				Message: "Email is required",
			})
		} else if input.Email == "invalid" {
			errors = append(errors, ValidationError{
				Field:   "email",
				Message: "Invalid email format",
				Value:   input.Email,
			})
		}

		// 패스워드 검증
		if len(input.Password) < 6 {
			errors = append(errors, ValidationError{
				Field:   "password",
				Message: "Password must be at least 6 characters",
			})
		}

		// 나이 검증
		if input.Age < 18 {
			errors = append(errors, ValidationError{
				Field:   "age",
				Message: "Must be 18 or older",
				Value:   fmt.Sprintf("%d", input.Age),
			})
		}

		if len(errors) > 0 {
			ValidationFailed(c, errors)
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"message": "Validation passed",
		}, nil)
	})

	// 429 Too Many Requests
	r.GET("/api/rate-limited", func(c *gin.Context) {
		NewErrorResponse(c, http.StatusTooManyRequests, ErrTooManyRequests,
			"Rate limit exceeded", gin.H{
				"limit":       100,
				"remaining":   0,
				"reset_after": "60 seconds",
			})
	})

	// ========================================
	// 3. 서버 에러 (5xx)
	// ========================================

	// 500 Internal Server Error
	r.GET("/api/error", func(c *gin.Context) {
		// 에러 시뮬레이션
		errorType := c.Query("type")

		switch errorType {
		case "db":
			NewErrorResponse(c, http.StatusInternalServerError, ErrDatabaseConnection,
				"Database connection failed", gin.H{
					"retry_after": "30 seconds",
				})
		case "panic":
			// 패닉 시뮬레이션 (리커버리 미들웨어가 처리)
			panic("Something went terribly wrong!")
		default:
			InternalServerError(c, "An unexpected error occurred")
		}
	})

	// 502 Bad Gateway
	r.GET("/api/external", func(c *gin.Context) {
		NewErrorResponse(c, http.StatusBadGateway, ErrExternalService,
			"External service is not responding", gin.H{
				"service": "payment-gateway",
				"timeout": "30s",
			})
	})

	// 503 Service Unavailable
	r.GET("/api/maintenance", func(c *gin.Context) {
		NewErrorResponse(c, http.StatusServiceUnavailable, ErrServiceUnavailable,
			"Service is under maintenance", gin.H{
				"retry_after": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			})
	})

	// ========================================
	// 4. 비즈니스 로직 에러 처리
	// ========================================

	r.POST("/api/transfer", func(c *gin.Context) {
		var transfer struct {
			From   string  `json:"from"`
			To     string  `json:"to"`
			Amount float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&transfer); err != nil {
			BadRequest(c, "Invalid request body", err.Error())
			return
		}

		// 비즈니스 규칙 검증
		if transfer.Amount <= 0 {
			err := BusinessError{
				Code:    "INVALID_AMOUNT",
				Message: "Transfer amount must be positive",
				Status:  http.StatusBadRequest,
			}
			NewErrorResponse(c, err.Status, err.Code, err.Message, gin.H{
				"amount": transfer.Amount,
			})
			return
		}

		if transfer.Amount > 10000 {
			err := BusinessError{
				Code:    "AMOUNT_LIMIT_EXCEEDED",
				Message: "Transfer amount exceeds daily limit",
				Status:  http.StatusBadRequest,
			}
			NewErrorResponse(c, err.Status, err.Code, err.Message, gin.H{
				"amount": transfer.Amount,
				"limit":  10000,
			})
			return
		}

		// 잔액 부족 시뮬레이션
		if transfer.From == "poor-account" {
			err := BusinessError{
				Code:    "INSUFFICIENT_FUNDS",
				Message: "Insufficient funds in source account",
				Status:  http.StatusBadRequest,
			}
			NewErrorResponse(c, err.Status, err.Code, err.Message, gin.H{
				"available": 100,
				"requested": transfer.Amount,
			})
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"transaction_id": "txn_" + fmt.Sprintf("%d", time.Now().Unix()),
			"status":         "completed",
			"amount":         transfer.Amount,
		}, nil)
	})

	// ========================================
	// 5. 파일 업로드 에러 처리
	// ========================================

	r.POST("/api/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			BadRequest(c, "No file uploaded", err.Error())
			return
		}

		// 파일 크기 체크 (5MB 제한)
		if file.Size > 5*1024*1024 {
			NewErrorResponse(c, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
				"File size exceeds maximum allowed size", gin.H{
					"max_size":     "5MB",
					"uploaded_size": fmt.Sprintf("%.2fMB", float64(file.Size)/(1024*1024)),
				})
			return
		}

		// 파일 타입 체크
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}

		if !allowedTypes[file.Header.Get("Content-Type")] {
			NewErrorResponse(c, http.StatusUnsupportedMediaType, "INVALID_FILE_TYPE",
				"File type not supported", gin.H{
					"allowed_types": []string{"image/jpeg", "image/png", "image/gif"},
					"uploaded_type": file.Header.Get("Content-Type"),
				})
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"filename": file.Filename,
			"size":     file.Size,
			"url":      "/uploads/" + file.Filename,
		}, nil)
	})

	// ========================================
	// 6. 페이지네이션 에러 처리
	// ========================================

	r.GET("/api/paginated", func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		var pageNum, limitNum int
		if _, err := fmt.Sscanf(page, "%d", &pageNum); err != nil || pageNum < 1 {
			BadRequest(c, "Invalid page parameter", gin.H{
				"page":         page,
				"valid_format": "positive integer",
			})
			return
		}

		if _, err := fmt.Sscanf(limit, "%d", &limitNum); err != nil || limitNum < 1 || limitNum > 100 {
			BadRequest(c, "Invalid limit parameter", gin.H{
				"limit":        limit,
				"valid_range":  "1-100",
			})
			return
		}

		// 페이지가 범위를 벗어난 경우
		totalItems := 50
		totalPages := (totalItems + limitNum - 1) / limitNum
		if pageNum > totalPages {
			NewErrorResponse(c, http.StatusBadRequest, "PAGE_OUT_OF_RANGE",
				"Page number exceeds total pages", gin.H{
					"requested_page": pageNum,
					"total_pages":    totalPages,
				})
			return
		}

		NewSuccessResponse(c, http.StatusOK, gin.H{
			"items": []string{"item1", "item2"},
		}, gin.H{
			"page":        pageNum,
			"limit":       limitNum,
			"total_items": totalItems,
			"total_pages": totalPages,
		})
	})

	// ========================================
	// 7. API 버전 에러
	// ========================================

	r.Any("/api/*path", func(c *gin.Context) {
		version := c.GetHeader("API-Version")

		if version != "" && version < "2.0" {
			NewErrorResponse(c, http.StatusGone, "API_VERSION_DEPRECATED",
				"This API version is no longer supported", gin.H{
					"requested_version": version,
					"minimum_version":   "2.0",
					"current_version":   "3.0",
				})
			return
		}

		// 존재하지 않는 엔드포인트
		NotFound(c, "Endpoint")
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("Test endpoints:")
	fmt.Println("  - GET  /api/users          (200 OK)")
	fmt.Println("  - POST /api/users          (201 Created or 409 Conflict)")
	fmt.Println("  - GET  /api/users/999      (404 Not Found)")
	fmt.Println("  - GET  /api/protected      (401 Unauthorized)")
	fmt.Println("  - POST /api/validate       (422 Validation Error)")
	fmt.Println("  - GET  /api/error?type=db  (500 Internal Server Error)")

	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}