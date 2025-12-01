package main

import (
	"fmt"

	"time"

	"os"

	"math/big"

	"github.com/gin-gonic/gin"

	"crypto/rand"

	"github.com/go-redis/redis/v8"
)

var base = os.Getenv("BASE_URL")
var port = os.Getenv("PORT")

var addrRedis = os.Getenv("REDIS_ADDR")
var rdb *redis.Client

func main() {
	// Read configuration from environment with sensible defaults for local development.
	if addrRedis == "" {
		addrRedis = "localhost:6379"
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: addrRedis, // Redis server address
		DB:   0,         // Default DB
	})
	if port == "" {
		port = "9808"
	}
	if base == "" {
		base = fmt.Sprintf("http://localhost:%s", port)
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
	r.GET("/:id", redirectHandler)

	// Allow port override with PORT env var (useful in Docker)
	addr := fmt.Sprintf(":%s", port)
	err := r.Run(addr)
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
	if err1 := c.ShouldBindJSON(&body); err1 != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	slug, err2 := SecureRandomString(8)
	if err2 != nil {
		c.JSON(500, gin.H{"error": "Failed to generate slug"})
		return
	}
	ctx := c.Request.Context()
	ttl := time.Hour * 24
	err3 := rdb.Set(ctx, slug, body.URL, ttl)
	if err3.Err() != nil {
		c.JSON(500, gin.H{"error": "Failed to add to redis"})
		return
	}
	c.JSON(200, gin.H{"shortUrl": fmt.Sprintf("%s/%s", base, slug)})

}

// ==============================================
// ==============================================
func redirectHandler(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("id")
	res, err := rdb.Get(ctx, slug).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{"NotFound": "missing url"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}
	c.Redirect(302, res)
}

func SecureRandomString(n int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}

	return string(result), nil
}
