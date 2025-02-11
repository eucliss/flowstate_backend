package flowstate

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	AddAuth(router, logger)
	InitUserDB(logger, "test")
	// user := User{Username: "admin", Password: "password"}
	// user.Create()
	return router
}

func TestLogin(t *testing.T) {
	router := setupAuthTestRouter()

	tests := []struct {
		name           string
		username       string
		password       string
		expectedCode   int
		expectedCookie bool
	}{
		{
			name:           "Valid Admin Login",
			username:       "admin",
			password:       "password",
			expectedCode:   http.StatusSeeOther,
			expectedCookie: true,
		},
		{
			name:           "Invalid Credentials",
			username:       "wrong",
			password:       "wrong",
			expectedCode:   http.StatusUnauthorized,
			expectedCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := User{
				Username: tt.username,
				Password: tt.password,
			}
			req, _ := http.NewRequest("POST", "/login", createJSONReader(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCookie {
				assert.Contains(t, w.Header().Get("Set-Cookie"), "Authorization")
			}
		})
	}
}

func TestLogout(t *testing.T) {
	router := setupAuthTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logout", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Set-Cookie"), "Authorization=;")
}

func TestCreateJWT(t *testing.T) {
	token, err := CreateJWT("testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token
	parsedToken, err := verifyToken(token)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "testuser", claims["sub"])
	assert.Equal(t, "flowstate-api", claims["iss"])
}

func TestAuthenticateMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(AuthenticateMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "protected")
	})

	tests := []struct {
		name         string
		setupAuth    func() string
		expectedCode int
	}{
		{
			name: "Valid Token",
			setupAuth: func() string {
				token, _ := CreateJWT("testuser")
				return token
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "No Token",
			setupAuth: func() string {
				return ""
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Token",
			setupAuth: func() string {
				return "invalid.token.here"
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)

			if token := tt.setupAuth(); token != "" {
				req.AddCookie(&http.Cookie{
					Name:  "Authorization",
					Value: token,
				})
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

// Helper function to create JSON reader
func createJSONReader(v interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(v)
	return &buf
}
