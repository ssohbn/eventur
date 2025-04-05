package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the database
	DBclient := connectDB()
	if DBclient == nil {
		log.Fatal("Failed to connect to MongoDB")
	}

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
		log.Printf("%+v\n", event)
		createEvent(DBclient, event)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//front end routes
	r.StaticFile("/", "src/static/index.html")
	r.StaticFile("/index", "src/static/index.html")
	r.StaticFile("/create", "src/static/create.html")
	r.StaticFile("/profile", "src/static/profile.html")
	r.StaticFile("/events", "src/static/events.html")

	// run the server
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
