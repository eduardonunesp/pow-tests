package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bwesterb/go-pow"
	"github.com/gin-gonic/gin"
)

func testPOW(difficulty uint32) (string, error) {
	start := time.Now()

	token := make([]byte, 16)
	rand.Read(token)

	req := pow.NewRequest(uint32(difficulty), token)
	proof, _ := pow.Fulfil(req, token)
	ok, err := pow.Check(req, proof, token)

	if !ok {
		return "", fmt.Errorf("failed to check the pow")
	}

	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}

	return fmt.Sprintf("%.6fms", float64(time.Since(start))/1e6), nil
}

func httpServer() {
	r := gin.Default()
	r.GET("/work", func(c *gin.Context) {
		difficultyParam := c.Query("difficulty")

		if len(difficultyParam) == 0 {
			difficultyParam = "10"
		}

		log.Println("Difficulty", difficultyParam)

		difficulty, err := strconv.ParseInt(difficultyParam, 10, 32)

		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}

		result, err := testPOW(uint32(difficulty))

		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}

		c.JSON(200, gin.H{
			"message": "done",
			"time":    result,
		})
	})

	r.Run("0.0.0.0:8080")
}

func commandLine(difficultyParam string) {
	if len(difficultyParam) == 0 {
		difficultyParam = "10"
	}

	log.Println("Difficulty", difficultyParam)

	difficulty, err := strconv.ParseInt(difficultyParam, 10, 32)

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	result, err := testPOW(uint32(difficulty))

	log.Println("time: ", result)
}

func main() {
	difficulty, ok := os.LookupEnv("DIFFICULTY")

	if !ok {
		httpServer()
	} else {
		commandLine(difficulty)
	}
}
