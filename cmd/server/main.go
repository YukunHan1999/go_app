package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/myapp/config"
	"github.com/myapp/internal/db"
	"github.com/myapp/internal/router"
)

func main() {

	// load config
	cfg := config.LoadConfig()

	// init db
	database := db.Connect(cfg.DbName)

	// register router
	r := router.RegisterRoute(database)

	log.Println("Server running on :8080")
	r.Use(MyMiddleware())
	r.Use(Recovery())
	r.GET("/ping", MyMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}

func MyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)

		log.Printf("[%d] Request %s %s took %v", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
