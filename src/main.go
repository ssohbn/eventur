package main

import (
	"net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  //routes
	r.StaticFile("/", "src/static/index.html")
  r.StaticFile("/create", "src/static/create.html")
  r.StaticFile("/profile", "src/static/profile.html")
  r.StaticFile("/events", "src/static/events.html")

  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  // run the server
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}