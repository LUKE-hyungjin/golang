package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/example/banking-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req_%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests (Lesson 11)
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body if needed
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		c.Next()

		// Log after processing
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Mask sensitive data in request body
		maskedBody := maskSensitiveData(string(requestBody))

		log.WithFields(map[string]interface{}{
			"request_id":  c.GetString("request_id"),
			"method":      method,
			"path":        path,
			"status":      statusCode,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   clientIP,
			"user_agent":  c.Request.UserAgent(),
			"body_size":   len(requestBody),
			"body":        maskedBody,
			"error":       c.Errors.String(),
		}).Info("HTTP Request")

		// Audit log for important operations
		if shouldAudit(method, path, statusCode) {
			// Audit logging would be implemented here
		}
	}
}

// RecoveryMiddleware recovers from panics (Lesson 10)
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				stackTrace := string(buf[:n])

				log.WithFields(map[string]interface{}{
					"request_id": c.GetString("request_id"),
					"error":      err,
					"stack":      stackTrace,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				// Return error response
				c.JSON(500, gin.H{
					"error":      "Internal server error",
					"request_id": c.GetString("request_id"),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ErrorHandler middleware handles errors globally (Lesson 10)
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Determine status code
			status := c.Writer.Status()
			if status == 200 {
				status = 500
			}

			// Create error response based on error type
			response := gin.H{
				"error":      err.Error(),
				"request_id": c.GetString("request_id"),
				"timestamp":  time.Now().Unix(),
			}

			// Add details in debug mode
			if gin.Mode() == gin.DebugMode {
				response["details"] = err.Meta
				response["type"] = fmt.Sprintf("%T", err.Err)
			}

			c.JSON(status, response)
		}
	}
}

// SecurityHeaders adds security headers (Lesson 14 - Release mode)
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

// RateLimiter implements rate limiting (Lesson 14 - Release mode)
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
}

type visitor struct {
	lastSeen time.Time
	count    int
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func RateLimiter(limit int) gin.HandlerFunc {
	// Clean up old visitors every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			limiter.mu.Lock()
			for ip, v := range limiter.visitors {
				if time.Since(v.lastSeen) > time.Minute {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		v, exists := limiter.visitors[ip]
		if !exists {
			limiter.visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			limiter.mu.Unlock()
			c.Next()
			return
		}

		// Reset count if last seen was more than a minute ago
		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
		} else {
			v.count++
			if v.count > limit {
				limiter.mu.Unlock()
				c.JSON(429, gin.H{
					"error": "Too many requests",
					"retry_after": 60,
				})
				c.Abort()
				return
			}
		}
		v.lastSeen = time.Now()
		limiter.mu.Unlock()

		c.Next()
	}
}

// DebugMiddleware for debug mode (Lesson 14)
func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log detailed request information
		fmt.Printf("[DEBUG] %s %s from %s\n", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		fmt.Printf("[DEBUG] Headers: %v\n", c.Request.Header)

		// Memory stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("[DEBUG] Memory - Alloc: %v MB, Sys: %v MB, NumGC: %v\n",
			m.Alloc/1024/1024, m.Sys/1024/1024, m.NumGC)

		// Timing
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		fmt.Printf("[DEBUG] Request completed in %v\n", latency)
	}
}

// TimeoutMiddleware sets request timeout (Lesson 17)
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to track if handler completes
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Handler completed normally
		case <-ctx.Done():
			// Timeout occurred
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
				"request_id": c.GetString("request_id"),
			})
			c.Abort()
		}
	}
}

// CORS middleware
func CORS(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, o := range origins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper functions
func maskSensitiveData(data string) string {
	// Mask passwords
	data = maskField(data, "password")
	data = maskField(data, "pin")
	data = maskField(data, "cvv")
	data = maskField(data, "card_number")
	data = maskField(data, "account_number")
	return data
}

func maskField(data, field string) string {
	// Simple masking - in production use proper JSON parsing
	if strings.Contains(data, field) {
		// This is a simplified version
		return strings.ReplaceAll(data, field, field+"_masked")
	}
	return data
}

func shouldAudit(method, path string, statusCode int) bool {
	// Audit important operations
	importantPaths := []string{"/transfer", "/withdraw", "/deposit", "/users", "/accounts"}
	for _, p := range importantPaths {
		if strings.Contains(path, p) && (method == "POST" || method == "PUT" || method == "DELETE") {
			return true
		}
	}
	return false
}

