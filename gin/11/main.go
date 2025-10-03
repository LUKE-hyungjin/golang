package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

// ========================================
// 로그 레벨 정의
// ========================================

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// ========================================
// 로그 엔트리 구조체
// ========================================

// LogEntry - 구조화된 로그 엔트리
type LogEntry struct {
	Level      string                 `json:"level"`
	Timestamp  string                 `json:"timestamp"`
	RequestID  string                 `json:"request_id"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	StatusCode int                    `json:"status_code"`
	Latency    string                 `json:"latency"`
	ClientIP   string                 `json:"client_ip"`
	UserAgent  string                 `json:"user_agent"`
	Error      string                 `json:"error,omitempty"`
	Request    *RequestLog           `json:"request,omitempty"`
	Response   *ResponseLog          `json:"response,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// RequestLog - 요청 로그 정보
type RequestLog struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Query   string            `json:"query,omitempty"`
	Form    map[string]string `json:"form,omitempty"`
}

// ResponseLog - 응답 로그 정보
type ResponseLog struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Size    int               `json:"size"`
}

// ========================================
// 커스텀 로거 인터페이스
// ========================================

// Logger - 로거 인터페이스
type Logger interface {
	Debug(entry LogEntry)
	Info(entry LogEntry)
	Warn(entry LogEntry)
	Error(entry LogEntry)
	Fatal(entry LogEntry)
}

// JSONLogger - JSON 형식 로거
type JSONLogger struct {
	output *log.Logger
	level  LogLevel
}

func NewJSONLogger(level LogLevel) *JSONLogger {
	return &JSONLogger{
		output: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func (l *JSONLogger) log(level LogLevel, entry LogEntry) {
	if level < l.level {
		return
	}

	entry.Level = level.String()
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}

	data, _ := json.Marshal(entry)
	l.output.Println(string(data))
}

func (l *JSONLogger) Debug(entry LogEntry) { l.log(DEBUG, entry) }
func (l *JSONLogger) Info(entry LogEntry)  { l.log(INFO, entry) }
func (l *JSONLogger) Warn(entry LogEntry)  { l.log(WARN, entry) }
func (l *JSONLogger) Error(entry LogEntry) { l.log(ERROR, entry) }
func (l *JSONLogger) Fatal(entry LogEntry) { l.log(FATAL, entry) }

// FileLogger - 파일 로거
type FileLogger struct {
	logger   Logger
	file     *os.File
	maxSize  int64
	filePath string
}

func NewFileLogger(filePath string, maxSize int64, level LogLevel) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	jsonLogger := &JSONLogger{
		output: log.New(file, "", 0),
		level:  level,
	}

	return &FileLogger{
		logger:   jsonLogger,
		file:     file,
		maxSize:  maxSize,
		filePath: filePath,
	}, nil
}

func (f *FileLogger) checkRotation() {
	info, err := f.file.Stat()
	if err != nil {
		return
	}

	if info.Size() >= f.maxSize {
		f.rotate()
	}
}

func (f *FileLogger) rotate() {
	f.file.Close()

	// 기존 파일 백업
	backupPath := fmt.Sprintf("%s.%s", f.filePath, time.Now().Format("20060102150405"))
	os.Rename(f.filePath, backupPath)

	// 새 파일 생성
	file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	f.file = file
	if jsonLogger, ok := f.logger.(*JSONLogger); ok {
		jsonLogger.output = log.New(file, "", 0)
	}
}

func (f *FileLogger) Debug(entry LogEntry) { f.checkRotation(); f.logger.Debug(entry) }
func (f *FileLogger) Info(entry LogEntry)  { f.checkRotation(); f.logger.Info(entry) }
func (f *FileLogger) Warn(entry LogEntry)  { f.checkRotation(); f.logger.Warn(entry) }
func (f *FileLogger) Error(entry LogEntry) { f.checkRotation(); f.logger.Error(entry) }
func (f *FileLogger) Fatal(entry LogEntry) { f.checkRotation(); f.logger.Fatal(entry) }

func (f *FileLogger) Close() {
	if f.file != nil {
		f.file.Close()
	}
}

// ========================================
// Response Writer 래퍼
// ========================================

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ========================================
// 로깅 미들웨어
// ========================================

// BasicLoggingMiddleware - 기본 로깅 미들웨어
func BasicLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %v | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
	})
}

// StructuredLoggingMiddleware - 구조화된 로깅 미들웨어
func StructuredLoggingMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Request ID 생성
		requestID := fmt.Sprintf("req-%d", start.UnixNano())
		c.Set("RequestID", requestID)

		// Request body 캡처 (필요한 경우)
		var requestBody []byte
		if c.Request.Body != nil && shouldLogBody(c.Request.Header.Get("Content-Type")) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
		}

		// Response writer 래핑
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 요청 처리
		c.Next()

		// 로그 엔트리 생성
		latency := time.Since(start)
		entry := LogEntry{
			Timestamp:  start.Format(time.RFC3339),
			RequestID:  requestID,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			Latency:    latency.String(),
			ClientIP:   c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		// Request 정보 추가
		if gin.Mode() == gin.DebugMode {
			entry.Request = &RequestLog{
				Headers: getHeaders(c.Request.Header),
				Query:   c.Request.URL.RawQuery,
			}

			if len(requestBody) > 0 && len(requestBody) < 1024 { // 1KB 이하만 로깅
				entry.Request.Body = string(requestBody)
			}
		}

		// Response 정보 추가
		if gin.Mode() == gin.DebugMode && blw.body.Len() > 0 && blw.body.Len() < 1024 {
			entry.Response = &ResponseLog{
				Headers: getHeaders(c.Writer.Header()),
				Body:    blw.body.String(),
				Size:    blw.body.Len(),
			}
		}

		// 에러 정보 추가
		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
		}

		// 로그 레벨 결정 및 로깅
		switch {
		case c.Writer.Status() >= 500:
			logger.Error(entry)
		case c.Writer.Status() >= 400:
			logger.Warn(entry)
		case c.Writer.Status() >= 300:
			logger.Info(entry)
		default:
			logger.Info(entry)
		}
	}
}

// AccessLoggingMiddleware - 접근 로그 미들웨어
func AccessLoggingMiddleware(accessLog *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Apache Combined Log Format
		accessLog.Printf(`%s - - [%s] "%s %s %s" %d %d "%s" "%s" %v`,
			c.ClientIP(),
			start.Format("02/Jan/2006:15:04:05 -0700"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			c.Writer.Size(),
			c.Request.Referer(),
			c.Request.UserAgent(),
			latency,
		)
	}
}

// SlowRequestLoggingMiddleware - 느린 요청 로깅
func SlowRequestLoggingMiddleware(threshold time.Duration, logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		if latency > threshold {
			entry := LogEntry{
				RequestID:  c.GetString("RequestID"),
				Method:     c.Request.Method,
				Path:       c.Request.URL.Path,
				StatusCode: c.Writer.Status(),
				Latency:    latency.String(),
				ClientIP:   c.ClientIP(),
				Extra: map[string]interface{}{
					"threshold": threshold.String(),
					"exceeded":  (latency - threshold).String(),
				},
			}
			logger.Warn(entry)
		}
	}
}

// ErrorLoggingMiddleware - 에러 전용 로깅
func ErrorLoggingMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				entry := LogEntry{
					RequestID:  c.GetString("RequestID"),
					Method:     c.Request.Method,
					Path:       c.Request.URL.Path,
					StatusCode: c.Writer.Status(),
					ClientIP:   c.ClientIP(),
					Error:      err.Error(),
					Extra: map[string]interface{}{
						"type": err.Type,
						"meta": err.Meta,
					},
				}
				logger.Error(entry)
			}
		}
	}
}

// AuditLoggingMiddleware - 감사 로그 미들웨어
func AuditLoggingMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 사용자 정보 (인증 미들웨어에서 설정된다고 가정)
		userID, _ := c.Get("UserID")
		userRole, _ := c.Get("UserRole")

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// 중요한 작업에 대한 감사 로그
		if isAuditRequired(c.Request.Method, c.Request.URL.Path) {
			entry := LogEntry{
				Timestamp:  start.Format(time.RFC3339),
				RequestID:  c.GetString("RequestID"),
				Method:     c.Request.Method,
				Path:       c.Request.URL.Path,
				StatusCode: c.Writer.Status(),
				Latency:    latency.String(),
				ClientIP:   c.ClientIP(),
				Extra: map[string]interface{}{
					"user_id":   userID,
					"user_role": userRole,
					"action":    getActionType(c.Request.Method),
					"resource":  getResourceType(c.Request.URL.Path),
				},
			}

			// Request body 로깅 (민감한 정보 제외)
			if c.Request.Method != "GET" {
				body, _ := c.GetRawData()
				if len(body) > 0 {
					sanitized := sanitizeBody(body)
					entry.Extra["request_body"] = sanitized
				}
			}

			logger.Info(entry)
		}
	}
}

// ========================================
// 헬퍼 함수들
// ========================================

func getHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		// 민감한 헤더 제외
		if !isSensitiveHeader(key) {
			result[key] = strings.Join(values, ",")
		}
	}
	return result
}

func isSensitiveHeader(header string) bool {
	sensitive := []string{"Authorization", "Cookie", "X-Api-Key", "X-Auth-Token"}
	header = strings.ToLower(header)
	for _, s := range sensitive {
		if strings.ToLower(s) == header {
			return true
		}
	}
	return false
}

func shouldLogBody(contentType string) bool {
	// JSON, XML, form data만 로깅
	allowedTypes := []string{"application/json", "application/xml", "application/x-www-form-urlencoded"}
	for _, t := range allowedTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}
	return false
}

func isAuditRequired(method, path string) bool {
	// POST, PUT, DELETE 요청은 감사 대상
	if method == "POST" || method == "PUT" || method == "DELETE" {
		return true
	}
	// 특정 경로는 GET도 감사
	auditPaths := []string{"/api/admin", "/api/users", "/api/sensitive"}
	for _, p := range auditPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func getActionType(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	case "GET":
		return "READ"
	default:
		return method
	}
}

func getResourceType(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func sanitizeBody(body []byte) string {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return string(body)
	}

	// 민감한 필드 마스킹
	sensitiveFields := []string{"password", "token", "secret", "credit_card"}
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***MASKED***"
		}
	}

	sanitized, _ := json.Marshal(data)
	return string(sanitized)
}

// ========================================
// 메인 함수 및 데모 핸들러
// ========================================

func main() {
	// 로거 초기화
	jsonLogger := NewJSONLogger(INFO)

	// 파일 로거 초기화
	logDir := "./logs"
	os.MkdirAll(logDir, 0755)

	fileLogger, err := NewFileLogger(
		filepath.Join(logDir, "app.log"),
		10*1024*1024, // 10MB
		INFO,
	)
	if err != nil {
		panic("Failed to create file logger: " + err.Error())
	}
	defer fileLogger.Close()

	// Access 로그 파일
	accessLogFile, _ := os.OpenFile(
		filepath.Join(logDir, "access.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	defer accessLogFile.Close()
	accessLogger := log.New(accessLogFile, "", 0)

	// Gin 설정
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	// 로깅 미들웨어 적용
	r.Use(StructuredLoggingMiddleware(jsonLogger))        // 구조화된 로깅
	r.Use(AccessLoggingMiddleware(accessLogger))          // 접근 로그
	r.Use(SlowRequestLoggingMiddleware(100*time.Millisecond, fileLogger)) // 느린 요청 로깅
	r.Use(ErrorLoggingMiddleware(fileLogger))             // 에러 로깅
	r.Use(AuditLoggingMiddleware(fileLogger))             // 감사 로그

	// 인증 시뮬레이션 미들웨어
	r.Use(func(c *gin.Context) {
		// 토큰에서 사용자 정보 추출 (시뮬레이션)
		if token := c.GetHeader("Authorization"); token != "" {
			c.Set("UserID", "user123")
			c.Set("UserRole", "admin")
		}
		c.Next()
	})

	// ========================================
	// 테스트 엔드포인트
	// ========================================

	// 1. 정상 요청
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// 2. 느린 요청
	r.GET("/api/slow", func(c *gin.Context) {
		time.Sleep(200 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"message": "This was a slow request",
		})
	})

	// 3. 에러 발생
	r.GET("/api/error", func(c *gin.Context) {
		c.Error(fmt.Errorf("this is a test error"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})

	// 4. 사용자 생성 (감사 로그 대상)
	r.POST("/api/users", func(c *gin.Context) {
		var user map[string]interface{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":      "user-123",
			"message": "User created successfully",
			"user":    user,
		})
	})

	// 5. 민감한 데이터 처리
	r.POST("/api/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid credentials format",
			})
			return
		}

		// 로그에는 password가 마스킹되어 기록됨
		c.JSON(http.StatusOK, gin.H{
			"token": "jwt-token-here",
			"user": gin.H{
				"id":       "user-123",
				"username": credentials.Username,
			},
		})
	})

	// 6. 대용량 응답
	r.GET("/api/large", func(c *gin.Context) {
		data := make([]gin.H, 1000)
		for i := range data {
			data[i] = gin.H{
				"id":    i,
				"value": fmt.Sprintf("Item %d", i),
			}
		}
		c.JSON(http.StatusOK, data)
	})

	// 7. 파일 다운로드
	r.GET("/api/download", func(c *gin.Context) {
		c.Header("Content-Disposition", "attachment; filename=test.txt")
		c.Data(http.StatusOK, "text/plain", []byte("This is a test file"))
	})

	// 8. 다양한 상태 코드
	r.GET("/api/status/:code", func(c *gin.Context) {
		code := c.Param("code")
		switch code {
		case "400":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		case "401":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		case "403":
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		case "404":
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		case "500":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		default:
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
	})

	// 9. 관리자 작업 (감사 로그)
	r.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("User %s deleted", id),
		})
	})

	// 10. 배치 작업
	r.POST("/api/batch", func(c *gin.Context) {
		// 여러 작업 수행 시뮬레이션
		for i := 0; i < 5; i++ {
			time.Sleep(50 * time.Millisecond)
		}
		c.JSON(http.StatusOK, gin.H{
			"processed": 5,
			"message":   "Batch processing completed",
		})
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("Logs are being written to ./logs/")
	fmt.Println("\nTest endpoints:")
	fmt.Println("  GET  /api/health     - Normal request")
	fmt.Println("  GET  /api/slow       - Slow request (triggers slow log)")
	fmt.Println("  GET  /api/error      - Error request")
	fmt.Println("  POST /api/users      - Create user (audit log)")
	fmt.Println("  POST /api/login      - Login (sensitive data masking)")
	fmt.Println("  GET  /api/status/:code - Various status codes")

	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}