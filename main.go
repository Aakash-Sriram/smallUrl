package main

import (
	"fmt"

	"os"

	"math/big"

	"github.com/gin-gonic/gin"

	"crypto/rand"

	"github.com/go-redis/redis/v8"
)

func main() {

	base := os.Getenv("BASE_URL")
	if base == "" {
		panic("BASE_URL not set in environment variables")
	}
	//==============================================
	//==============================================

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "No clue what am doing",
		})
	})

	r.POST("/shorten", tinyUrl)

	r.POST("/:id", redirectHandler)
	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
	//==============================================
	//==============================================

}

type Request struct {
	URL string `json:"url" binding:"required"`
}

func tinyUrl(c *gin.Context) {
	var body Request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	slug, err := SecureRandomString(8)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate slug"})
		return
	}
	c.JSON(200, gin.H{"shortUrl": slug})

}

// ==============================================
// ==============================================
func redirectHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"url": "shortened_url",
	})
}
func setUpRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:8080", // Redis server address
		DB:   0,                // Default DB
	})
	return rdb
}
func SecureRandomString(n int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, n)

	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}

	return string(result), nil
}
