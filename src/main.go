package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
  
	"time"
)

type Event struct {
	Title string `form:"title" binding:"required"`
	Blurb string `form:"blurb"`

	// not required because tbd dates are allowed
	Date time.Time `form:"date" time_format:"2006-01-02"`

	// not required because tbd locations are allowed
	// maybe should be required for like city or something
	Location string `form:"location"`
}

func main() {
	r := gin.Default()

	r.POST("/api/createEvent", func(c *gin.Context) {
		log.Println("recv'd ")
		var event Event

		// bind form data (WHICH SHOULD BIND) to event and check if error is produced (not nil)
		// (sorta weird syntax)
		if err := c.ShouldBind(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// If data binding is successful, return the user information
		c.JSON(http.StatusOK, gin.H{"message": "Event Created!", "event": event})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//front end routes
  r.StaticFile("/", "src/static/index.html")
  r.StaticFile("/create", "src/static/create.html")
  r.StaticFile("/profile", "src/static/profile.html")
  r.StaticFile("/events", "src/static/events.html")

	// run the server
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
