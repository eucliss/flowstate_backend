package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type helloWorld struct {
	Message string `json:"message" form:"message"`
}

type user struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

var currentState = helloWorld{Message: "Hello, World!"}
var secretKey = []byte("testing-secret-key")
var loggedInUsers = make(map[string]user)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Duration: %v", c.Request.Method, c.Writer.Status(), duration)
	}
}

func createToken(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,        // subject
		"iss": "flowstate-api", // issuer
		// "aud": "getRole(username)", // audience
		"iat": time.Now().Unix(),                     // issued at
		"exp": time.Now().Add(time.Hour * 24).Unix(), // expires at
	})
	fmt.Printf("Token claims added: %v\n", claims)

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func authenticateMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	token, err := verifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	username, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
		c.Abort()
		return
	}

	c.Set("username", username)
	c.Next()
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if its valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(LoggerMiddleware())
	router.GET("/", authenticateMiddleware, func(c *gin.Context) {
		c.Bind(&currentState)
		c.JSON(http.StatusOK, currentState)
	})

	router.GET("/logout", func(c *gin.Context) {
		c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	})

	router.POST("/login", func(c *gin.Context) {
		var u user
		c.Bind(&u)

		if u.Username == "admin" && u.Password == "password" || u.Username == "user" && u.Password == "password" {
			fmt.Println("Logging in user:", u.Username)
			token, err := createToken(u.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error creating token": err.Error()})
				return
			}

			loggedInUsers[u.Username] = u

			c.SetCookie("Authorization", token, 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusSeeOther, "/")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		}
	})

	router.POST("/set", authenticateMiddleware, func(c *gin.Context) {
		c.Bind(&currentState)
		fmt.Println("Setting state:", currentState)
		c.JSON(http.StatusOK, currentState)
	})

	router.POST("/reset", authenticateMiddleware, func(c *gin.Context) {
		currentState = helloWorld{Message: "Hello, World!"}
		c.JSON(http.StatusOK, currentState)
	})

	router.Run(":8080")
}
