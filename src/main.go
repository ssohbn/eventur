package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "html/template"
)

func main() {
	// connect to the database
	DBclient, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	log.Println("connected to db!")

	defer func() {
		if err := DBclient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

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

		// try inserting into db
		err := createEvent(DBclient, event)
		if err != nil {
			log.Printf("failed to create event %s", err)
		}

		// If data binding is successful, return the user information
		c.JSON(http.StatusOK, gin.H{"message": "Event Created!", "event": event})
		log.Printf("%+v\n", event)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
      "title": "Main website",
      "isIndex":true,
    })
  })
  r.GET("/create", func(c *gin.Context) {
    c.HTML(http.StatusOK, "create.html", gin.H{
      "title": "Main website",
    })
  })
  r.GET("/profile", func(c *gin.Context) {
    c.HTML(http.StatusOK, "profile.html", gin.H{
      "title": "Main website",
    })
  })
  r.GET("/events", func(c *gin.Context) {
    c.HTML(http.StatusOK, "events.html", gin.H{
      "title": "Main website",
      "isEvents": true,
      "Events": []string{" "," "," "},
    })
  })

	//front end routes
	r.Static("/css", "src/static/css")
	r.Static("/imgs", "src/static/imgs")

	// run the server
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
