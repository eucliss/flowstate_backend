package main

import (
	"flowstate/flowstate"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type helloWorld struct {
	Message string `json:"message" form:"message"`
}

var currentState = helloWorld{Message: "Hello, World!"}

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logger.Info(
			"Request",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration,
		)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flowstate.InitUserDB(logger, "prod")

	// admin := flowstate.User{Username: "admin", Password: "password"}
	// flowstate.CreateUser(db, &admin)
	a := flowstate.User{Username: "admin"}
	user := a.Get()
	fmt.Println(user)

	return

	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(LoggerMiddleware(logger))

	authMiddleware := flowstate.AddAuth(router, logger)

	router.GET("/", authMiddleware, func(c *gin.Context) {
		c.Bind(&currentState)
		c.JSON(http.StatusOK, currentState)
	})

	router.POST("/set", authMiddleware, func(c *gin.Context) {
		c.Bind(&currentState)
		fmt.Println("Setting state:", currentState)
		c.JSON(http.StatusOK, currentState)
	})

	router.POST("/reset", authMiddleware, func(c *gin.Context) {
		currentState = helloWorld{Message: "Hello, World!"}
		c.JSON(http.StatusOK, currentState)
	})

	router.Run(":8080")
}
