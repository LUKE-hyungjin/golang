package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ========================================
// 에러 타입 정의
// ========================================

// AppError - 애플리케이션 에러 인터페이스
type AppError interface {
	error
	Status() int
	Code() string
	Details() interface{}
}

// HTTPError - HTTP 에러 구현체
type HTTPError struct {
	StatusCode int         `json:"status"`
	ErrorCode  string      `json:"code"`
	Message    string      `json:"message"`
	Detail     interface{} `json:"details,omitempty"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) Status() int {
	return e.StatusCode
}

func (e HTTPError) Code() string {
	return e.ErrorCode
}

func (e HTTPError) Details() interface{} {
	return e.Detail
}

// ValidationError - 검증 에러
type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
	return "Validation failed"
}

// BusinessError - 비즈니스 로직 에러
type BusinessError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e BusinessError) Error() string {
	return e.Message
}

// ErrorResponse - 통일된 에러 응답
type ErrorResponse struct {
	Success   bool       `json:"success"`
	Error     ErrorInfo  `json:"error"`
	RequestID string     `json:"request_id"`
	Timestamp time.Time  `json:"timestamp"`
}

type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Stack   string      `json:"stack,omitempty"` // 개발 모드에서만 표시
}

// ========================================
// 에러 핸들링 미들웨어
// ========================================

// ErrorHandlingMiddleware - 전역 에러 핸들링 미들웨어
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청 처리
		c.Next()

		// 에러가 있는지 확인
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

// RecoveryMiddleware - 패닉 복구 미들웨어
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 스택 트레이스 로깅
				stack := string(debug.Stack())
				fmt.Printf("PANIC RECOVERED: %v\n%s\n", err, stack)

				// 클라이언트 응답
				var message string
				switch v := err.(type) {
				case string:
					message = v
				case error:
					message = v.Error()
				default:
					message = "An unexpected error occurred"
				}

				respondWithError(c, http.StatusInternalServerError, "PANIC_ERROR", message, nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ValidationErrorMiddleware - 바인딩 에러 처리 미들웨어
func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 바인딩 에러 확인
		if c.Writer.Status() == http.StatusBadRequest {
			if len(c.Errors) > 0 {
				var fieldErrors []FieldError
				for _, err := range c.Errors {
					if strings.Contains(err.Error(), "binding") {
						// 바인딩 에러 파싱
						fieldErrors = append(fieldErrors, FieldError{
							Field:   "unknown",
							Message: err.Error(),
						})
					}
				}

				if len(fieldErrors) > 0 {
					respondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR",
						"Input validation failed", fieldErrors)
				}
			}
		}
	}
}

// CustomErrorMiddleware - 커스텀 에러 처리 미들웨어
func CustomErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Response Writer 래핑
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 에러 상태 코드 확인
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			// 이미 처리된 에러가 아닌 경우
			if blw.body.Len() == 0 {
				handleHTTPError(c, statusCode)
			}
		}
	}
}

// bodyLogWriter - Response Body 캡처용
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ========================================
// 에러 처리 헬퍼 함수
// ========================================

// handleError - 에러 타입별 처리
func handleError(c *gin.Context, err error) {
	// AppError 인터페이스 구현 확인
	var appErr AppError
	if errors.As(err, &appErr) {
		respondWithError(c, appErr.Status(), appErr.Code(),
			appErr.Error(), appErr.Details())
		return
	}

	// HTTPError 타입 확인
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		respondWithError(c, httpErr.StatusCode, httpErr.ErrorCode,
			httpErr.Message, httpErr.Detail)
		return
	}

	// ValidationError 타입 확인
	var valErr ValidationError
	if errors.As(err, &valErr) {
		respondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR",
			valErr.Error(), valErr.Errors)
		return
	}

	// BusinessError 타입 확인
	var bizErr BusinessError
	if errors.As(err, &bizErr) {
		respondWithError(c, http.StatusBadRequest, bizErr.Code,
			bizErr.Message, bizErr.Data)
		return
	}

	// 기본 에러 처리
	respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
		"An internal error occurred", nil)
}

// handleHTTPError - HTTP 상태 코드별 에러 처리
func handleHTTPError(c *gin.Context, statusCode int) {
	var code, message string

	switch statusCode {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
		message = "Invalid request"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
		message = "Authentication required"
	case http.StatusForbidden:
		code = "FORBIDDEN"
		message = "Access denied"
	case http.StatusNotFound:
		code = "NOT_FOUND"
		message = "Resource not found"
	case http.StatusMethodNotAllowed:
		code = "METHOD_NOT_ALLOWED"
		message = "Method not allowed"
	case http.StatusConflict:
		code = "CONFLICT"
		message = "Resource conflict"
	case http.StatusTooManyRequests:
		code = "TOO_MANY_REQUESTS"
		message = "Too many requests"
	case http.StatusInternalServerError:
		code = "INTERNAL_ERROR"
		message = "Internal server error"
	case http.StatusServiceUnavailable:
		code = "SERVICE_UNAVAILABLE"
		message = "Service temporarily unavailable"
	default:
		code = "UNKNOWN_ERROR"
		message = "An error occurred"
	}

	respondWithError(c, statusCode, code, message, nil)
}

// respondWithError - 통일된 에러 응답 전송
func respondWithError(c *gin.Context, status int, code string, message string, details interface{}) {
	requestID, _ := c.Get("RequestID")

	errorInfo := ErrorInfo{
		Code:    code,
		Message: message,
		Details: details,
	}

	// 개발 모드에서 스택 트레이스 추가
	if gin.Mode() == gin.DebugMode {
		errorInfo.Stack = string(debug.Stack())
	}

	response := ErrorResponse{
		Success:   false,
		Error:     errorInfo,
		RequestID: fmt.Sprintf("%v", requestID),
		Timestamp: time.Now(),
	}

	c.JSON(status, response)
}

// ========================================
// 에러 발생 헬퍼 함수
// ========================================

// NewHTTPError - HTTP 에러 생성
func NewHTTPError(status int, code string, message string) HTTPError {
	return HTTPError{
		StatusCode: status,
		ErrorCode:  code,
		Message:    message,
	}
}

// NewValidationError - 검증 에러 생성
func NewValidationError(errors []FieldError) ValidationError {
	return ValidationError{
		Errors: errors,
	}
}

// NewBusinessError - 비즈니스 에러 생성
func NewBusinessError(code string, message string) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}

// ========================================
// 데모 API 핸들러
// ========================================

// User - 사용자 모델
type User struct {
	ID    int    `json:"id" binding:"required,min=1"`
	Name  string `json:"name" binding:"required,min=3,max=50"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=18,max=120"`
}

func main() {
	// Gin 인스턴스 생성 (기본 미들웨어 없이)
	r := gin.New()

	// Request ID 미들웨어
	r.Use(func(c *gin.Context) {
		c.Set("RequestID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
		c.Next()
	})

	// 에러 핸들링 미들웨어 적용
	r.Use(RecoveryMiddleware())        // 패닉 복구
	r.Use(ErrorHandlingMiddleware())   // 전역 에러 핸들링
	r.Use(ValidationErrorMiddleware()) // 검증 에러 처리

	// ========================================
	// 테스트 엔드포인트
	// ========================================

	// 1. 정상 응답
	r.GET("/api/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    "Operation successful",
		})
	})

	// 2. HTTPError 발생
	r.GET("/api/http-error", func(c *gin.Context) {
		err := NewHTTPError(http.StatusNotFound, "USER_NOT_FOUND",
			"The requested user does not exist")
		c.Error(err)
		c.Abort()
	})

	// 3. ValidationError 발생
	r.POST("/api/validation-error", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			// 수동 검증 에러 생성
			fieldErrors := []FieldError{
				{Field: "email", Message: "Invalid email format", Value: "not-an-email"},
				{Field: "age", Message: "Must be 18 or older", Value: "16"},
			}
			c.Error(NewValidationError(fieldErrors))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
	})

	// 4. BusinessError 발생
	r.POST("/api/business-error", func(c *gin.Context) {
		// 비즈니스 로직 검증
		amount := c.Query("amount")
		if amount == "0" {
			err := NewBusinessError("INVALID_AMOUNT",
				"Transaction amount must be greater than zero")
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction processed",
		})
	})

	// 5. 패닉 발생 (RecoveryMiddleware가 처리)
	r.GET("/api/panic", func(c *gin.Context) {
		panic("This is a deliberate panic!")
	})

	// 6. 일반 에러 발생
	r.GET("/api/generic-error", func(c *gin.Context) {
		err := errors.New("something went wrong")
		c.Error(err)
		c.Abort()
	})

	// 7. 체인된 에러 처리
	r.GET("/api/chained-error", func(c *gin.Context) {
		// 데이터베이스 에러 시뮬레이션
		dbErr := errors.New("connection refused")
		wrappedErr := fmt.Errorf("database error: %w", dbErr)

		httpErr := HTTPError{
			StatusCode: http.StatusServiceUnavailable,
			ErrorCode:  "DATABASE_ERROR",
			Message:    "Unable to connect to database",
			Detail:     wrappedErr.Error(),
		}

		c.Error(httpErr)
		c.Abort()
	})

	// 8. 다중 에러 처리
	r.POST("/api/multiple-errors", func(c *gin.Context) {
		// 여러 검증 에러 누적
		c.Error(errors.New("First error"))
		c.Error(errors.New("Second error"))
		c.Error(NewHTTPError(http.StatusBadRequest, "MULTIPLE_ERRORS",
			"Multiple validation failures"))
		c.Abort()
	})

	// 9. 파일 업로드 에러
	r.POST("/api/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			httpErr := NewHTTPError(http.StatusBadRequest, "FILE_REQUIRED",
				"No file was uploaded")
			c.Error(httpErr)
			c.Abort()
			return
		}

		// 파일 크기 체크
		if file.Size > 5*1024*1024 { // 5MB
			httpErr := HTTPError{
				StatusCode: http.StatusRequestEntityTooLarge,
				ErrorCode:  "FILE_TOO_LARGE",
				Message:    "File size exceeds maximum allowed",
				Detail: map[string]interface{}{
					"max_size":      "5MB",
					"uploaded_size": fmt.Sprintf("%.2fMB", float64(file.Size)/(1024*1024)),
				},
			}
			c.Error(httpErr)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"filename": file.Filename,
			"size":     file.Size,
		})
	})

	// 10. 타임아웃 에러
	r.GET("/api/timeout", func(c *gin.Context) {
		select {
		case <-time.After(5 * time.Second):
			c.JSON(http.StatusOK, gin.H{"success": true})
		case <-c.Request.Context().Done():
			err := NewHTTPError(http.StatusRequestTimeout, "REQUEST_TIMEOUT",
				"Request processing timed out")
			c.Error(err)
			c.Abort()
		}
	})

	// 11. 권한 체크 에러
	r.GET("/api/admin", func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			err := NewHTTPError(http.StatusUnauthorized, "AUTH_REQUIRED",
				"Authentication token is required")
			c.Error(err)
			c.Abort()
			return
		}

		if token != "Bearer admin-token" {
			err := NewHTTPError(http.StatusForbidden, "INSUFFICIENT_PRIVILEGES",
				"You don't have permission to access this resource")
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Welcome admin!",
		})
	})

	// 12. 외부 서비스 에러
	r.GET("/api/external", func(c *gin.Context) {
		// 외부 API 호출 시뮬레이션
		resp, err := http.Get("http://nonexistent-service.local/api")
		if err != nil {
			httpErr := HTTPError{
				StatusCode: http.StatusBadGateway,
				ErrorCode:  "EXTERNAL_SERVICE_ERROR",
				Message:    "Failed to connect to external service",
				Detail: map[string]string{
					"service": "payment-gateway",
					"error":   err.Error(),
				},
			}
			c.Error(httpErr)
			c.Abort()
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    string(body),
		})
	})

	// 13. 데이터 처리 에러
	r.POST("/api/process", func(c *gin.Context) {
		var data json.RawMessage
		if err := c.ShouldBindJSON(&data); err != nil {
			httpErr := NewHTTPError(http.StatusBadRequest, "INVALID_JSON",
				"Failed to parse JSON data")
			c.Error(httpErr)
			c.Abort()
			return
		}

		// 처리 중 에러 시뮬레이션
		if string(data) == `{"trigger":"error"}` {
			err := NewBusinessError("PROCESSING_FAILED",
				"Unable to process the provided data")
			c.Error(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data processed successfully",
		})
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("\nTest endpoints:")
	fmt.Println("  GET  /api/success          - Success response")
	fmt.Println("  GET  /api/http-error       - HTTP error example")
	fmt.Println("  POST /api/validation-error - Validation error")
	fmt.Println("  POST /api/business-error   - Business logic error")
	fmt.Println("  GET  /api/panic           - Panic recovery")
	fmt.Println("  GET  /api/generic-error   - Generic error")
	fmt.Println("  GET  /api/admin           - Authorization error")
	fmt.Println("  POST /api/upload          - File upload error")

	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}