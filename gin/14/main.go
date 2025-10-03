package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 실행 모드 타입 정의
// ============================================================================

type RunMode string

const (
	DebugMode   RunMode = "debug"
	ReleaseMode RunMode = "release"
	TestMode    RunMode = "test"
)

// ============================================================================
// 모드별 설정 구조체
// ============================================================================

type ModeConfig struct {
	Mode            RunMode
	LogLevel        string
	LogOutput       io.Writer
	EnableProfiling bool
	EnableMetrics   bool
	EnableSwagger   bool
	ErrorDetails    bool
	PanicRecovery   bool
	RequestLogging  bool
	ResponseLogging bool
	ColoredOutput   bool
	MaxMemory       int64 // bytes
	MaxCPU          int
	RateLimit       int
	Timeout         time.Duration
}

// 모드별 기본 설정
func GetModeConfig(mode RunMode) *ModeConfig {
	switch mode {
	case DebugMode:
		return &ModeConfig{
			Mode:            DebugMode,
			LogLevel:        "debug",
			LogOutput:       os.Stdout,
			EnableProfiling: true,
			EnableMetrics:   true,
			EnableSwagger:   true,
			ErrorDetails:    true,
			PanicRecovery:   true,
			RequestLogging:  true,
			ResponseLogging: true,
			ColoredOutput:   true,
			MaxMemory:       0, // unlimited
			MaxCPU:          0, // use all cores
			RateLimit:       0, // no limit
			Timeout:         30 * time.Second,
		}
	case ReleaseMode:
		return &ModeConfig{
			Mode:            ReleaseMode,
			LogLevel:        "info",
			LogOutput:       os.Stdout, // or file
			EnableProfiling: false,
			EnableMetrics:   true,
			EnableSwagger:   false,
			ErrorDetails:    false,
			PanicRecovery:   true,
			RequestLogging:  true,
			ResponseLogging: false,
			ColoredOutput:   false,
			MaxMemory:       1 << 30, // 1GB
			MaxCPU:          runtime.NumCPU(),
			RateLimit:       100, // requests per minute
			Timeout:         15 * time.Second,
		}
	case TestMode:
		return &ModeConfig{
			Mode:            TestMode,
			LogLevel:        "error",
			LogOutput:       io.Discard, // suppress logs
			EnableProfiling: false,
			EnableMetrics:   false,
			EnableSwagger:   false,
			ErrorDetails:    true,
			PanicRecovery:   false, // let tests catch panics
			RequestLogging:  false,
			ResponseLogging: false,
			ColoredOutput:   false,
			MaxMemory:       1 << 28, // 256MB
			MaxCPU:          2,
			RateLimit:       0,
			Timeout:         5 * time.Second,
		}
	default:
		return GetModeConfig(ReleaseMode)
	}
}

// ============================================================================
// 모드별 미들웨어
// ============================================================================

// Debug 모드 전용 미들웨어
func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청 정보 상세 로깅
		log.Printf("[DEBUG] %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		log.Printf("[DEBUG] Headers: %v", c.Request.Header)

		// 메모리 사용량 추적
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.Printf("[DEBUG] Memory - Alloc: %v MB, Sys: %v MB, NumGC: %v",
			m.Alloc/1024/1024, m.Sys/1024/1024, m.NumGC)

		// 요청 처리 시간 측정
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		log.Printf("[DEBUG] Request completed in %v", latency)
	}
}

// Release 모드 전용 미들웨어
func ReleaseMiddleware(config *ModeConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 보안 헤더 추가
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Rate limiting
		if config.RateLimit > 0 {
			// 실제 구현에서는 Redis 등을 사용한 분산 rate limiting
			// 여기서는 간단한 예시만
		}

		c.Next()

		// 민감한 정보 제거
		c.Header("Server", "")
	}
}

// Test 모드 전용 미들웨어
func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 테스트용 헤더 추가
		c.Header("X-Test-Mode", "true")

		// 테스트 요청 ID 생성
		testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
		c.Set("test_id", testID)

		c.Next()
	}
}

// 프로파일링 미들웨어
func ProfilingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/debug/pprof" {
			pprof.Index(c.Writer, c.Request)
			c.Abort()
		} else if c.Request.URL.Path == "/debug/pprof/cmdline" {
			pprof.Cmdline(c.Writer, c.Request)
			c.Abort()
		} else if c.Request.URL.Path == "/debug/pprof/profile" {
			pprof.Profile(c.Writer, c.Request)
			c.Abort()
		} else if c.Request.URL.Path == "/debug/pprof/symbol" {
			pprof.Symbol(c.Writer, c.Request)
			c.Abort()
		} else if c.Request.URL.Path == "/debug/pprof/trace" {
			pprof.Trace(c.Writer, c.Request)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// ============================================================================
// 모드별 에러 핸들러
// ============================================================================

func DebugErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 상세한 에러 정보 포함
			c.JSON(c.Writer.Status(), gin.H{
				"error":       err.Error(),
				"type":        fmt.Sprintf("%T", err.Err),
				"meta":        err.Meta,
				"stack_trace": string(debug.Stack()),
				"request_id":  c.GetString("request_id"),
				"path":        c.Request.URL.Path,
				"method":      c.Request.Method,
			})
		}
	}
}

func ReleaseErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// 일반적인 에러 메시지만 노출
			status := c.Writer.Status()
			message := "Internal server error"

			switch status {
			case 400:
				message = "Bad request"
			case 401:
				message = "Unauthorized"
			case 403:
				message = "Forbidden"
			case 404:
				message = "Not found"
			}

			c.JSON(status, gin.H{
				"error":      message,
				"request_id": c.GetString("request_id"),
			})
		}
	}
}

func TestErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 테스트에 유용한 정보 포함
			c.JSON(c.Writer.Status(), gin.H{
				"error":   err.Error(),
				"test_id": c.GetString("test_id"),
			})
		}
	}
}

// ============================================================================
// 모드별 라우터 설정
// ============================================================================

func SetupDebugRouter(config *ModeConfig) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	// Debug 모드 미들웨어
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(DebugMiddleware())
	router.Use(DebugErrorHandler())

	// 프로파일링 활성화
	if config.EnableProfiling {
		router.Use(ProfilingMiddleware())
		router.GET("/debug/vars", expvarHandler)
		router.GET("/debug/gc", gcHandler)
		router.GET("/debug/mem", memStatsHandler)
	}

	// Swagger UI
	if config.EnableSwagger {
		router.GET("/swagger/*any", func(c *gin.Context) {
			c.JSON(200, gin.H{"swagger": "UI would be here"})
		})
	}

	return router
}

func SetupReleaseRouter(config *ModeConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Release 모드 미들웨어
	router.Use(gin.Recovery())
	router.Use(ReleaseMiddleware(config))
	router.Use(ReleaseErrorHandler())

	// 메트릭스만 활성화
	if config.EnableMetrics {
		router.GET("/metrics", metricsHandler)
	}

	return router
}

func SetupTestRouter(config *ModeConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Test 모드 미들웨어
	router.Use(TestMiddleware())
	router.Use(TestErrorHandler())

	return router
}

// ============================================================================
// 디버그 핸들러
// ============================================================================

func expvarHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"goroutines": runtime.NumGoroutine(),
		"cpu":        runtime.NumCPU(),
		"memory": gin.H{
			"alloc":   runtime.MemStats{}.Alloc,
			"total":   runtime.MemStats{}.TotalAlloc,
			"sys":     runtime.MemStats{}.Sys,
			"numGC":   runtime.MemStats{}.NumGC,
		},
	})
}

func gcHandler(c *gin.Context) {
	runtime.GC()
	c.JSON(200, gin.H{"message": "GC triggered"})
}

func memStatsHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(200, gin.H{
		"alloc":         m.Alloc,
		"total_alloc":   m.TotalAlloc,
		"sys":           m.Sys,
		"lookups":       m.Lookups,
		"mallocs":       m.Mallocs,
		"frees":         m.Frees,
		"heap_alloc":    m.HeapAlloc,
		"heap_sys":      m.HeapSys,
		"heap_idle":     m.HeapIdle,
		"heap_inuse":    m.HeapInuse,
		"heap_released": m.HeapReleased,
		"heap_objects":  m.HeapObjects,
		"num_gc":        m.NumGC,
		"gc_cpu_fraction": m.GCCPUFraction,
	})
}

func metricsHandler(c *gin.Context) {
	// Prometheus 형식의 메트릭스
	c.String(200, `# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 142

# HELP http_request_duration_seconds HTTP request duration
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 120
http_request_duration_seconds_bucket{le="0.5"} 135
http_request_duration_seconds_bucket{le="1"} 140

# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines %d`, runtime.NumGoroutine())
}

// ============================================================================
// 애플리케이션 설정
// ============================================================================

type Application struct {
	Router *gin.Engine
	Config *ModeConfig
	Mode   RunMode
}

func NewApplication(mode RunMode) *Application {
	config := GetModeConfig(mode)

	// 리소스 제한 설정
	if config.MaxMemory > 0 {
		debug.SetMemoryLimit(config.MaxMemory)
	}

	if config.MaxCPU > 0 {
		runtime.GOMAXPROCS(config.MaxCPU)
	}

	// 라우터 설정
	var router *gin.Engine
	switch mode {
	case DebugMode:
		router = SetupDebugRouter(config)
	case ReleaseMode:
		router = SetupReleaseRouter(config)
	case TestMode:
		router = SetupTestRouter(config)
	default:
		router = SetupReleaseRouter(config)
	}

	app := &Application{
		Router: router,
		Config: config,
		Mode:   mode,
	}

	// 공통 라우트 설정
	app.setupRoutes()

	return app
}

func (app *Application) setupRoutes() {
	// 헬스체크
	app.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"mode":   string(app.Mode),
			"time":   time.Now().Unix(),
		})
	})

	// 모드 정보
	app.Router.GET("/mode", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"mode":            string(app.Mode),
			"debug":           app.Mode == DebugMode,
			"release":         app.Mode == ReleaseMode,
			"test":            app.Mode == TestMode,
			"profiling":       app.Config.EnableProfiling,
			"metrics":         app.Config.EnableMetrics,
			"swagger":         app.Config.EnableSwagger,
			"error_details":   app.Config.ErrorDetails,
			"request_logging": app.Config.RequestLogging,
			"colored_output":  app.Config.ColoredOutput,
		})
	})

	// 모드 전환 (개발 환경에서만)
	if app.Mode == DebugMode {
		app.Router.POST("/mode/:mode", func(c *gin.Context) {
			newMode := RunMode(c.Param("mode"))

			if newMode != DebugMode && newMode != ReleaseMode && newMode != TestMode {
				c.JSON(400, gin.H{"error": "Invalid mode"})
				return
			}

			c.JSON(200, gin.H{
				"message":      "Mode change requires restart",
				"current_mode": string(app.Mode),
				"new_mode":     string(newMode),
			})
		})
	}

	// 샘플 API 엔드포인트
	api := app.Router.Group("/api")
	{
		// 정상 응답
		api.GET("/users", func(c *gin.Context) {
			users := []gin.H{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			}
			c.JSON(200, users)
		})

		// 에러 테스트
		api.GET("/error", func(c *gin.Context) {
			c.Error(fmt.Errorf("test error"))
			c.Status(500)
		})

		// 패닉 테스트
		api.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		// 느린 응답 테스트
		api.GET("/slow", func(c *gin.Context) {
			time.Sleep(2 * time.Second)
			c.JSON(200, gin.H{"message": "slow response"})
		})

		// 메모리 사용 테스트
		api.GET("/memory", func(c *gin.Context) {
			// 10MB 할당
			data := make([]byte, 10*1024*1024)
			for i := range data {
				data[i] = byte(i % 256)
			}
			c.JSON(200, gin.H{
				"allocated": "10MB",
				"checksum":  len(data),
			})
		})
	}

	// 모드별 특수 엔드포인트
	switch app.Mode {
	case DebugMode:
		app.Router.GET("/debug/config", func(c *gin.Context) {
			c.JSON(200, app.Config)
		})

		app.Router.GET("/debug/routes", func(c *gin.Context) {
			routes := app.Router.Routes()
			c.JSON(200, routes)
		})

		app.Router.GET("/debug/env", func(c *gin.Context) {
			c.JSON(200, os.Environ())
		})

	case TestMode:
		app.Router.POST("/test/reset", func(c *gin.Context) {
			// 테스트 데이터 리셋
			c.JSON(200, gin.H{"message": "Test data reset"})
		})

		app.Router.POST("/test/seed", func(c *gin.Context) {
			// 테스트 데이터 시딩
			c.JSON(200, gin.H{"message": "Test data seeded"})
		})
	}
}

// ============================================================================
// 서버 실행
// ============================================================================

func (app *Application) Run(addr string) error {
	log.Printf("🚀 Starting server in %s mode on %s", app.Mode, addr)
	log.Printf("📊 Configuration:")
	log.Printf("  - Profiling: %v", app.Config.EnableProfiling)
	log.Printf("  - Metrics: %v", app.Config.EnableMetrics)
	log.Printf("  - Request Logging: %v", app.Config.RequestLogging)
	log.Printf("  - Max CPU: %d", app.Config.MaxCPU)

	// 타임아웃 설정
	server := &http.Server{
		Addr:         addr,
		Handler:      app.Router,
		ReadTimeout:  app.Config.Timeout,
		WriteTimeout: app.Config.Timeout,
		IdleTimeout:  app.Config.Timeout * 2,
	}

	return server.ListenAndServe()
}

// ============================================================================
// Main
// ============================================================================

func main() {
	// 환경변수에서 모드 읽기
	mode := RunMode(os.Getenv("GIN_MODE"))
	if mode == "" {
		mode = DebugMode
	}

	// 애플리케이션 생성
	app := NewApplication(mode)

	// 모드별 안내 메시지
	switch mode {
	case DebugMode:
		log.Println("🔍 Debug Mode - All features enabled")
		log.Println("📝 Profiling available at /debug/pprof")
		log.Println("🎨 Colored output enabled")
	case ReleaseMode:
		log.Println("🚀 Release Mode - Optimized for production")
		log.Println("🔒 Security headers enabled")
		log.Println("📊 Metrics available at /metrics")
	case TestMode:
		log.Println("🧪 Test Mode - Optimized for testing")
		log.Println("🔇 Logging suppressed")
	}

	// 서버 시작
	if err := app.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}