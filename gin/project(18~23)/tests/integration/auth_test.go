package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"project-security/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	router       *gin.Engine
	testServer   *TestServer
	testUser     *TestUser
	accessToken  string
	refreshToken string
}

type TestUser struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func (suite *AuthTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Initialize test server
	server, err := NewTestServer()
	suite.Require().NoError(err)
	suite.testServer = server
	suite.router = server.Router

	// Create test user
	suite.testUser = &TestUser{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "Test123!@#",
		FirstName: "Test",
		LastName:  "User",
	}
}

func (suite *AuthTestSuite) TearDownSuite() {
	suite.testServer.Cleanup()
}

func (suite *AuthTestSuite) TestRegister_Success() {
	req := validator.RegisterRequest{
		Username:        "newuser",
		Email:           "newuser@example.com",
		Password:        "NewUser123!@#",
		ConfirmPassword: "NewUser123!@#",
		FirstName:       "New",
		LastName:        "User",
	}

	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/api/v1/auth/register", body, nil)

	suite.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "user")
	suite.Contains(response, "access_token")
	suite.Contains(response, "refresh_token")
}

func (suite *AuthTestSuite) TestRegister_ValidationError() {
	testCases := []struct {
		name     string
		request  validator.RegisterRequest
		expected int
	}{
		{
			name: "Weak password",
			request: validator.RegisterRequest{
				Username:        "user1",
				Email:           "user1@example.com",
				Password:        "weak",
				ConfirmPassword: "weak",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Password mismatch",
			request: validator.RegisterRequest{
				Username:        "user2",
				Email:           "user2@example.com",
				Password:        "Strong123!@#",
				ConfirmPassword: "Different123!@#",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			request: validator.RegisterRequest{
				Username:        "user3",
				Email:           "invalid-email",
				Password:        "Strong123!@#",
				ConfirmPassword: "Strong123!@#",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body, _ := json.Marshal(tc.request)
			w := suite.performRequest("POST", "/api/v1/auth/register", body, nil)
			suite.Equal(tc.expected, w.Code)
		})
	}
}

func (suite *AuthTestSuite) TestLogin_Success() {
	// First register a user
	suite.registerTestUser()

	// Now test login
	req := validator.LoginRequest{
		Email:    suite.testUser.Email,
		Password: suite.testUser.Password,
	}

	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/api/v1/auth/login", body, nil)

	suite.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "user")
	suite.Contains(response, "access_token")
	suite.Contains(response, "refresh_token")

	// Store tokens for other tests
	suite.accessToken = response["access_token"].(string)
	suite.refreshToken = response["refresh_token"].(string)
}

func (suite *AuthTestSuite) TestLogin_InvalidCredentials() {
	req := validator.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "WrongPassword123!",
	}

	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/api/v1/auth/login", body, nil)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *AuthTestSuite) TestRefreshToken_Success() {
	// Login first to get tokens
	suite.registerTestUser()
	suite.login()

	req := validator.RefreshTokenRequest{
		RefreshToken: suite.refreshToken,
	}

	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/api/v1/auth/refresh", body, nil)

	suite.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "access_token")
	suite.Contains(response, "refresh_token")
}

func (suite *AuthTestSuite) TestProtectedRoute_WithToken() {
	suite.registerTestUser()
	suite.login()

	headers := map[string]string{
		"Authorization": "Bearer " + suite.accessToken,
	}

	w := suite.performRequest("GET", "/api/v1/users/profile", nil, headers)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *AuthTestSuite) TestProtectedRoute_WithoutToken() {
	w := suite.performRequest("GET", "/api/v1/users/profile", nil, nil)
	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *AuthTestSuite) TestProtectedRoute_InvalidToken() {
	headers := map[string]string{
		"Authorization": "Bearer invalid-token",
	}

	w := suite.performRequest("GET", "/api/v1/users/profile", nil, headers)
	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *AuthTestSuite) TestRoleBasedAccess() {
	// This would require setting up admin user
	// For now, test that non-admin can't access admin routes
	suite.registerTestUser()
	suite.login()

	headers := map[string]string{
		"Authorization": "Bearer " + suite.accessToken,
	}

	w := suite.performRequest("GET", "/api/v1/users", nil, headers)
	suite.Equal(http.StatusForbidden, w.Code)
}

// Helper methods

func (suite *AuthTestSuite) performRequest(method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *AuthTestSuite) registerTestUser() {
	req := validator.RegisterRequest{
		Username:        suite.testUser.Username,
		Email:           suite.testUser.Email,
		Password:        suite.testUser.Password,
		ConfirmPassword: suite.testUser.Password,
		FirstName:       suite.testUser.FirstName,
		LastName:        suite.testUser.LastName,
	}

	body, _ := json.Marshal(req)
	suite.performRequest("POST", "/api/v1/auth/register", body, nil)
}

func (suite *AuthTestSuite) login() {
	req := validator.LoginRequest{
		Email:    suite.testUser.Email,
		Password: suite.testUser.Password,
	}

	body, _ := json.Marshal(req)
	w := suite.performRequest("POST", "/api/v1/auth/login", body, nil)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	suite.accessToken = response["access_token"].(string)
	suite.refreshToken = response["refresh_token"].(string)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// Benchmark tests

func BenchmarkLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)
	server, _ := NewTestServer()
	defer server.Cleanup()

	// Register user once
	registerReq := validator.RegisterRequest{
		Username:        "benchuser",
		Email:           "bench@example.com",
		Password:        "Bench123!@#",
		ConfirmPassword: "Bench123!@#",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	// Benchmark login
	loginReq := validator.LoginRequest{
		Email:    "bench@example.com",
		Password: "Bench123!@#",
	}
	loginBody, _ := json.Marshal(loginReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
	}
}

func BenchmarkTokenValidation(b *testing.B) {
	gin.SetMode(gin.TestMode)
	server, _ := NewTestServer()
	defer server.Cleanup()

	// Get a valid token
	// ... (register and login to get token)

	token := "valid-token-here"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
	}
}