package flowstate

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Change this
var secretKey = []byte("your-256-bit-secret")

var authLogger *slog.Logger

func AddAuth(router *gin.Engine, l *slog.Logger) func(c *gin.Context) {
	authLogger = l
	router.POST("/login", login)
	router.GET("/logout", logout)
	return AuthenticateMiddleware
}

func login(c *gin.Context) {
	var u User
	c.Bind(&u)
	authLogger.Info("Logging in user", "username", u.Username)
	if u.Exists() && u.LoginSuccess() {
		authLogger.Info("User exists and login successful")
		token, err := CreateJWT(u.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error creating token": err.Error()})
			return
		}

		c.SetCookie("Authorization", token, 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}

func logout(c *gin.Context) {
	// get username from claims
	authLogger.Info("Logging out user. - revoking token.")
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Create a JWT token for a user
func CreateJWT(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,        // subject
		"iss": "flowstate-api", // issuer
		// "aud": "getRole(username)", // audience
		"iat": time.Now().Unix(),                     // issued at
		"exp": time.Now().Add(time.Hour * 24).Unix(), // expires at
	})

	authLogger.Info("Token claims added", "claims", claims)

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Authenticate a user
func AuthenticateMiddleware(c *gin.Context) {
	// Get the token from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// Get the username from the claims
	username, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
		c.Abort()
		return
	}

	// Set the username in the context
	c.Set("username", username)

	// Continue to the next handler
	c.Next()
}

// Verify a JWT token
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if its valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Return the token
	return token, nil
}
