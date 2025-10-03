package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ========================================
// 설정 구조체 정의
// ========================================

// Config - 전체 설정 구조체
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
	Features FeatureFlags   `mapstructure:"features"`
	External ExternalAPIs   `mapstructure:"external"`
}

// ServerConfig - 서버 설정
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"` // debug, release, test
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	TrustedProxies  []string      `mapstructure:"trusted_proxies"`
}

// DatabaseConfig - 데이터베이스 설정
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig - Redis 설정
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// JWTConfig - JWT 설정
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	Issuer           string        `mapstructure:"issuer"`
	AccessExpiry     time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry    time.Duration `mapstructure:"refresh_expiry"`
	SigningAlgorithm string        `mapstructure:"signing_algorithm"`
}

// EmailConfig - 이메일 설정
type EmailConfig struct {
	SMTP     SMTPConfig `mapstructure:"smtp"`
	From     string     `mapstructure:"from"`
	FromName string     `mapstructure:"from_name"`
	ReplyTo  string     `mapstructure:"reply_to"`
}

// SMTPConfig - SMTP 설정
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
}

// StorageConfig - 스토리지 설정
type StorageConfig struct {
	Type      string `mapstructure:"type"` // local, s3, gcs
	LocalPath string `mapstructure:"local_path"`
	S3        S3Config `mapstructure:"s3"`
}

// S3Config - S3 설정
type S3Config struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Endpoint        string `mapstructure:"endpoint"`
}

// LoggingConfig - 로깅 설정
type LoggingConfig struct {
	Level      string `mapstructure:"level"` // debug, info, warn, error
	Format     string `mapstructure:"format"` // json, text
	Output     string `mapstructure:"output"` // stdout, file
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`     // days
}

// SecurityConfig - 보안 설정
type SecurityConfig struct {
	CORS          CORSConfig      `mapstructure:"cors"`
	RateLimit     RateLimitConfig `mapstructure:"rate_limit"`
	AllowedHosts  []string        `mapstructure:"allowed_hosts"`
	SSLRedirect   bool            `mapstructure:"ssl_redirect"`
	CSRFProtection bool           `mapstructure:"csrf_protection"`
}

// CORSConfig - CORS 설정
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// RateLimitConfig - Rate Limiting 설정
type RateLimitConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	BurstSize     int  `mapstructure:"burst_size"`
}

// FeatureFlags - 기능 플래그
type FeatureFlags struct {
	NewDashboard     bool `mapstructure:"new_dashboard"`
	BetaFeatures     bool `mapstructure:"beta_features"`
	MaintenanceMode  bool `mapstructure:"maintenance_mode"`
	DebugMode        bool `mapstructure:"debug_mode"`
	EnableMetrics    bool `mapstructure:"enable_metrics"`
	EnableProfiling  bool `mapstructure:"enable_profiling"`
}

// ExternalAPIs - 외부 API 설정
type ExternalAPIs struct {
	PaymentGateway APIConfig `mapstructure:"payment_gateway"`
	Analytics      APIConfig `mapstructure:"analytics"`
	Notification   APIConfig `mapstructure:"notification"`
}

// APIConfig - API 설정
type APIConfig struct {
	BaseURL string            `mapstructure:"base_url"`
	APIKey  string            `mapstructure:"api_key"`
	Timeout time.Duration     `mapstructure:"timeout"`
	Retry   int               `mapstructure:"retry"`
	Headers map[string]string `mapstructure:"headers"`
}

// ========================================
// 설정 로더
// ========================================

// ConfigLoader - 설정 로더 인터페이스
type ConfigLoader interface {
	Load() (*Config, error)
	Watch(callback func(*Config))
	Get(key string) interface{}
	Set(key string, value interface{})
}

// ViperConfigLoader - Viper 기반 설정 로더
type ViperConfigLoader struct {
	viper *viper.Viper
	config *Config
}

// NewConfigLoader - 새 설정 로더 생성
func NewConfigLoader(configPath string) ConfigLoader {
	v := viper.New()

	// 설정 파일 경로 설정
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 기본 설정 경로
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
	}

	// 환경 변수 설정
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 기본값 설정
	setDefaults(v)

	return &ViperConfigLoader{
		viper: v,
	}
}

// Load - 설정 로드
func (cl *ViperConfigLoader) Load() (*Config, error) {
	// 설정 파일 읽기
	if err := cl.viper.ReadInConfig(); err != nil {
		// 설정 파일이 없어도 환경변수와 기본값으로 동작
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		log.Println("Config file not found, using defaults and environment variables")
	} else {
		log.Printf("Using config file: %s", cl.viper.ConfigFileUsed())
	}

	// 환경별 설정 오버라이드
	env := cl.viper.GetString("app_env")
	if env != "" {
		envConfigPath := fmt.Sprintf("config.%s", env)
		cl.viper.SetConfigName(envConfigPath)
		if err := cl.viper.MergeInConfig(); err == nil {
			log.Printf("Merged environment config: %s", envConfigPath)
		}
	}

	// 구조체로 언마샬
	var config Config
	if err := cl.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 설정 검증
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cl.config = &config
	return &config, nil
}

// Watch - 설정 파일 변경 감시
func (cl *ViperConfigLoader) Watch(callback func(*Config)) {
	cl.viper.WatchConfig()
	cl.viper.OnConfigChange(func(e fsnotify.ConfigChangeEvent) {
		log.Printf("Config file changed: %s", e.Name)

		var config Config
		if err := cl.viper.Unmarshal(&config); err != nil {
			log.Printf("Failed to reload config: %v", err)
			return
		}

		if err := validateConfig(&config); err != nil {
			log.Printf("Config validation failed after reload: %v", err)
			return
		}

		cl.config = &config
		callback(&config)
	})
}

// Get - 설정 값 가져오기
func (cl *ViperConfigLoader) Get(key string) interface{} {
	return cl.viper.Get(key)
}

// Set - 설정 값 설정
func (cl *ViperConfigLoader) Set(key string, value interface{}) {
	cl.viper.Set(key, value)
}

// ========================================
// 헬퍼 함수들
// ========================================

// setDefaults - 기본값 설정
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.shutdown_timeout", "30s")
	v.SetDefault("server.max_header_bytes", 1<<20) // 1MB

	// Database defaults
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.dial_timeout", "5s")

	// JWT defaults
	v.SetDefault("jwt.signing_algorithm", "HS256")
	v.SetDefault("jwt.access_expiry", "15m")
	v.SetDefault("jwt.refresh_expiry", "7d")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.max_size", 100)
	v.SetDefault("logging.max_backups", 3)
	v.SetDefault("logging.max_age", 7)

	// Security defaults
	v.SetDefault("security.cors.enabled", true)
	v.SetDefault("security.cors.allow_origins", []string{"*"})
	v.SetDefault("security.cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("security.cors.allow_headers", []string{"Authorization", "Content-Type"})
	v.SetDefault("security.cors.max_age", 86400)

	v.SetDefault("security.rate_limit.enabled", true)
	v.SetDefault("security.rate_limit.requests_per_minute", 60)
	v.SetDefault("security.rate_limit.burst_size", 10)

	// Feature flags defaults
	v.SetDefault("features.new_dashboard", false)
	v.SetDefault("features.beta_features", false)
	v.SetDefault("features.maintenance_mode", false)
	v.SetDefault("features.enable_metrics", true)
}

// validateConfig - 설정 검증
func validateConfig(config *Config) error {
	// 필수 설정 검증
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Driver != "" && config.Database.Host == "" {
		return fmt.Errorf("database host is required when driver is set")
	}

	if config.JWT.Secret == "" && config.Server.Mode == "release" {
		return fmt.Errorf("JWT secret is required in release mode")
	}

	// 모드 검증
	validModes := []string{"debug", "release", "test"}
	validMode := false
	for _, mode := range validModes {
		if config.Server.Mode == mode {
			validMode = true
			break
		}
	}
	if !validMode {
		return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
	}

	return nil
}

// GetDatabaseDSN - 데이터베이스 연결 문자열 생성
func GetDatabaseDSN(config *DatabaseConfig) string {
	switch config.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database)
	default:
		return ""
	}
}

// ========================================
// 메인 함수 및 데모
// ========================================

func main() {
	// 설정 파일 경로 (명령행 인자나 환경변수로 받을 수 있음)
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" && len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// 설정 로더 생성
	configLoader := NewConfigLoader(configPath)

	// 설정 로드
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 설정 변경 감시 (옵션)
	configLoader.Watch(func(newConfig *Config) {
		log.Println("Configuration reloaded")
		// 여기서 필요한 재설정 작업 수행
	})

	// Gin 모드 설정
	gin.SetMode(config.Server.Mode)

	// Gin 라우터 생성
	r := gin.Default()

	// ========================================
	// 설정 정보 엔드포인트
	// ========================================

	// 1. 현재 설정 조회 (민감한 정보 제외)
	r.GET("/api/config", func(c *gin.Context) {
		// 민감한 정보 제거
		safeConfig := map[string]interface{}{
			"server": map[string]interface{}{
				"host": config.Server.Host,
				"port": config.Server.Port,
				"mode": config.Server.Mode,
			},
			"database": map[string]interface{}{
				"driver": config.Database.Driver,
				"host":   config.Database.Host,
				"port":   config.Database.Port,
			},
			"features": config.Features,
			"logging": map[string]interface{}{
				"level":  config.Logging.Level,
				"format": config.Logging.Format,
				"output": config.Logging.Output,
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"config":      safeConfig,
			"environment": os.Getenv("APP_ENV"),
			"config_file": viper.ConfigFileUsed(),
		})
	})

	// 2. 환경 변수 조회
	r.GET("/api/env", func(c *gin.Context) {
		envVars := make(map[string]string)
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "APP_") {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					// 민감한 정보 마스킹
					if strings.Contains(strings.ToLower(parts[0]), "password") ||
						strings.Contains(strings.ToLower(parts[0]), "secret") ||
						strings.Contains(strings.ToLower(parts[0]), "key") {
						envVars[parts[0]] = "***MASKED***"
					} else {
						envVars[parts[0]] = parts[1]
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"environment_variables": envVars,
		})
	})

	// 3. 기능 플래그 확인
	r.GET("/api/features", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"features": config.Features,
		})
	})

	// 4. 기능 플래그별 엔드포인트
	r.GET("/api/dashboard", func(c *gin.Context) {
		if !config.Features.NewDashboard {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "New dashboard is not enabled",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "New dashboard is available",
			"version": "2.0",
		})
	})

	// 5. 베타 기능
	r.GET("/api/beta", func(c *gin.Context) {
		if !config.Features.BetaFeatures {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Beta features are not enabled",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Beta features activated",
			"features": []string{"feature1", "feature2", "feature3"},
		})
	})

	// 6. 유지보수 모드
	r.Use(func(c *gin.Context) {
		if config.Features.MaintenanceMode && !strings.HasPrefix(c.Request.URL.Path, "/api/health") {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Service is under maintenance",
				"message": "We'll be back soon!",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// 7. 헬스 체크
	r.GET("/api/health", func(c *gin.Context) {
		health := gin.H{
			"status": "healthy",
			"server": gin.H{
				"mode": config.Server.Mode,
				"port": config.Server.Port,
			},
			"timestamp": time.Now(),
		}

		// 데이터베이스 연결 상태 (실제로는 ping 수행)
		if config.Database.Host != "" {
			health["database"] = gin.H{
				"connected": true, // 실제로는 DB 연결 확인
				"driver":    config.Database.Driver,
			}
		}

		// Redis 연결 상태
		if config.Redis.Host != "" {
			health["redis"] = gin.H{
				"connected": true, // 실제로는 Redis 연결 확인
				"host":      config.Redis.Host,
			}
		}

		c.JSON(http.StatusOK, health)
	})

	// 8. 설정 다시 로드 (관리자용)
	r.POST("/api/admin/reload-config", func(c *gin.Context) {
		// 실제로는 인증 확인 필요
		newConfig, err := configLoader.Load()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to reload config: %v", err),
			})
			return
		}

		config = newConfig
		c.JSON(http.StatusOK, gin.H{
			"message": "Configuration reloaded successfully",
		})
	})

	// 9. 런타임 설정 변경
	r.PUT("/api/admin/config", func(c *gin.Context) {
		var update map[string]interface{}
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 설정 업데이트
		for key, value := range update {
			configLoader.Set(key, value)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Configuration updated",
			"updated": update,
		})
	})

	// 10. 환경별 응답
	r.GET("/api/info", func(c *gin.Context) {
		info := gin.H{
			"environment": os.Getenv("APP_ENV"),
			"version":     os.Getenv("APP_VERSION"),
			"mode":        config.Server.Mode,
		}

		// 디버그 모드에서만 상세 정보 표시
		if config.Server.Mode == "debug" || config.Features.DebugMode {
			info["detailed"] = gin.H{
				"go_version":  runtime.Version(),
				"num_cpu":     runtime.NumCPU(),
				"num_goroutine": runtime.NumGoroutine(),
			}
		}

		c.JSON(http.StatusOK, info)
	})

	// 서버 시작
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	fmt.Println("===========================================")
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Printf("Mode: %s\n", config.Server.Mode)
	fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
	fmt.Println("===========================================")
	fmt.Println("\nEndpoints:")
	fmt.Println("  GET /api/config    - View current configuration")
	fmt.Println("  GET /api/env       - View environment variables")
	fmt.Println("  GET /api/features  - View feature flags")
	fmt.Println("  GET /api/health    - Health check")
	fmt.Println("  GET /api/info      - Server information")

	// 타임아웃 설정이 있는 서버 생성
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}