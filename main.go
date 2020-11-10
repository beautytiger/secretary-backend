package main

import (
	"daocloud.io/secretary-backend/pkg/meeting"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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
		c.File("/tmp/meeting.jpg")
	})
	meet.POST("/:id/record", meeting.PutRecords)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
