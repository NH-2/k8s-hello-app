package main

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		generate_load() // artificial load to simulate request processing
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Healthy": true,
		})
	})
	r.GET("/healthz/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Ready": true,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func generate_load() {
	// Number of logical CPUs on the local machine
	cores := runtime.NumCPU()

	// Generate load on each of the cores
	for i := 0; i < cores; i++ {
		go func() {
			for {
				// The tight loop generates the Load like crazy!
			}
		}()
	}

	// Let it run for a certain amount of time, e.g., 10 milli-seconds.
	time.Sleep(10 * time.Millisecond)
}
