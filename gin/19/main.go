package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// 모델 정의
// ============================================================================

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"` // Never expose password
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token_id"`
	jwt.RegisteredClaims
}

// ============================================================================
// JWT 설정
// ============================================================================

type JWTConfig struct {
	SecretKey           string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
	Issuer              string
	Audience            []string
}

var jwtConfig = JWTConfig{
	SecretKey:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	AccessTokenExpiry:  15 * time.Minute,
	RefreshTokenExpiry: 7 * 24 * time.Hour,
	Issuer:             "gin-jwt-example",
	Audience:           []string{"gin-api"},
}

// Mock database
var users = map[string]*User{
	"admin@example.com": {
		ID:       1,
		Email:    "admin@example.com",
		Username: "admin",
		Password: "$2a$10$YKxocKxRyCeH9XeRT2sUDOaQJCwFO3V9PGq8f7kBBJKtTjOxPxF7e", // password: admin123
		Role:     "admin",
	},
	"user@example.com": {
		ID:       2,
		Email:    "user@example.com",
		Username: "user",
		Password: "$2a$10$mQ0W6yfCfqJWVWDy0BfQJeRAZTu3bBG7TfKFKxMgDfhLi7P8gJkLa", // password: user123
		Role:     "user",
	},
}

// Store for refresh tokens (in production use Redis/Database)
var refreshTokenStore = make(map[string]uint) // token -> userID

// ============================================================================
// JWT Functions
// ============================================================================

// GenerateTokenPair generates both access and refresh tokens
func GenerateTokenPair(user *User) (*TokenResponse, error) {
	// Generate Access Token
	accessToken, expiresAt, err := generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(jwtConfig.AccessTokenExpiry.Seconds()),
		ExpiresAt:    expiresAt,
	}, nil
}

// generateAccessToken creates a new access token
func generateAccessToken(user *User) (string, time.Time, error) {
	expiresAt := time.Now().Add(jwtConfig.AccessTokenExpiry)

	claims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        generateTokenID(),
			Audience:  jwtConfig.Audience,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))

	return tokenString, expiresAt, err
}

// generateRefreshToken creates a new refresh token
func generateRefreshToken(user *User) (string, error) {
	tokenID := generateTokenID()

	claims := RefreshClaims{
		UserID: user.ID,
		Token:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))

	if err == nil {
		// Store refresh token
		refreshTokenStore[tokenString] = user.ID
	}

	return tokenString, err
}

// ValidateToken validates and parses the token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Additional validations
	if claims.Issuer != jwtConfig.Issuer {
		return nil, errors.New("invalid issuer")
	}

	// Check audience
	validAudience := false
	for _, aud := range claims.Audience {
		for _, expectedAud := range jwtConfig.Audience {
			if aud == expectedAud {
				validAudience = true
				break
			}
		}
	}
	if !validAudience {
		return nil, errors.New("invalid audience")
	}

	return claims, nil
}

// ValidateRefreshToken validates refresh token
func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token exists in store
	if _, exists := refreshTokenStore[tokenString]; !exists {
		return nil, errors.New("refresh token not found or revoked")
	}

	return claims, nil
}

// RevokeRefreshToken removes refresh token from store
func RevokeRefreshToken(tokenString string) {
	delete(refreshTokenStore, tokenString)
}

// generateTokenID generates a unique token ID
func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ============================================================================
// Middleware
// ============================================================================

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userClaims := claims.(*Claims)
		hasRole := false

		for _, role := range requiredRoles {
			if userClaims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware allows both authenticated and unauthenticated access
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
				if claims, err := ValidateToken(bearerToken[1]); err == nil {
					c.Set("claims", claims)
					c.Set("user_id", claims.UserID)
					c.Set("authenticated", true)
				}
			}
		}
		c.Next()
	}
}

// ============================================================================
// Handlers
// ============================================================================

// Register handler
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if _, exists := users[req.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := &User{
		ID:       uint(len(users) + 1),
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	users[req.Email] = user

	// Generate tokens
	tokens, err := GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
		"tokens": tokens,
	})
}

// Login handler
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	user, exists := users[req.Email]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	tokens, err := GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
		"tokens": tokens,
	})
}

// Refresh token handler
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Find user
	var user *User
	for _, u := range users {
		if u.ID == claims.UserID {
			user = u
			break
		}
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Revoke old refresh token
	RevokeRefreshToken(req.RefreshToken)

	// Generate new tokens
	tokens, err := GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens":  tokens,
	})
}

// Logout handler
func Logout(c *gin.Context) {
	// Get refresh token from request
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req)

	// Revoke refresh token if provided
	if req.RefreshToken != "" {
		RevokeRefreshToken(req.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Protected route handlers
func GetProfile(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*Claims)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userClaims.UserID,
			"email":    userClaims.Email,
			"username": userClaims.Username,
			"role":     userClaims.Role,
		},
		"token_info": gin.H{
			"issued_at":  userClaims.IssuedAt.Time,
			"expires_at": userClaims.ExpiresAt.Time,
			"issuer":     userClaims.Issuer,
		},
	})
}

func AdminOnly(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin access granted",
		"data":    "Secret admin data",
	})
}

func PublicEndpoint(c *gin.Context) {
	authenticated, _ := c.Get("authenticated")

	response := gin.H{
		"message": "This is a public endpoint",
		"authenticated": authenticated == true,
	}

	if authenticated == true {
		claims, _ := c.Get("claims")
		userClaims := claims.(*Claims)
		response["user"] = userClaims.Username
	}

	c.JSON(http.StatusOK, response)
}

// ============================================================================
// Router Setup
// ============================================================================

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/register", Register)
		public.POST("/login", Login)
		public.POST("/refresh", RefreshToken)
		public.GET("/public", OptionalAuthMiddleware(), PublicEndpoint)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/logout", Logout)
		protected.GET("/profile", GetProfile)
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "This is a protected endpoint"})
		})
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(AuthMiddleware(), RoleMiddleware("admin"))
	{
		admin.GET("/users", func(c *gin.Context) {
			userList := []gin.H{}
			for _, user := range users {
				userList = append(userList, gin.H{
					"id":       user.ID,
					"email":    user.Email,
					"username": user.Username,
					"role":     user.Role,
				})
			}
			c.JSON(http.StatusOK, userList)
		})
		admin.GET("/dashboard", AdminOnly)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// JWT info endpoint
	router.GET("/jwt/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"issuer":               jwtConfig.Issuer,
			"audience":             jwtConfig.Audience,
			"access_token_expiry":  jwtConfig.AccessTokenExpiry.String(),
			"refresh_token_expiry": jwtConfig.RefreshTokenExpiry.String(),
			"algorithm":            "HS256",
		})
	})

	return router
}

// ============================================================================
// Utility Functions
// ============================================================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ============================================================================
// Main
// ============================================================================

func main() {
	router := setupRouter()

	log.Println("🚀 JWT Authentication Server starting on :8080")
	log.Println("🔑 JWT configuration loaded")
	log.Println("")
	log.Println("Test credentials:")
	log.Println("  Admin: admin@example.com / admin123")
	log.Println("  User:  user@example.com / user123")
	log.Println("")
	log.Println("Try:")
	log.Println(`  curl -X POST http://localhost:8080/api/v1/login \`)
	log.Println(`    -H "Content-Type: application/json" \`)
	log.Println(`    -d '{"email":"admin@example.com","password":"admin123"}'`)
	log.Println("")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

