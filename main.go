package main

import (
	"daocloud.io/secretary-backend/pkg/meeting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	meeting.LoadLogFromDisk()
	SetupCloseHandler()
	r := gin.Default()
	r.Use(cors.Default())
	r.Static("/front", "./dist")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	meet := r.Group("/v1/meeting")
	meet.GET("/", meeting.List)
	meet.POST("/", meeting.Add)
	meet.GET("/:id", meeting.Detail)
	meet.GET("/:id/record", meeting.GetRecords)
	meet.GET("/:id/image", func(c *gin.Context) {
		c.File("./meeting.jpg")
	})
	meet.POST("/:id/record", meeting.PutRecords)
	meet.POST("/:id/speaker", meeting.PutSpeaker)
	meet.POST("/:id/words", meeting.PutWords)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("ticker come")
				meeting.PersistentLogToDisk()
			case <-c:
				log.Println("ticker done")
				log.Println("program cleaning up")
				meeting.PersistentLogToDisk()
				os.Exit(0)
			}
		}
	}()
}
