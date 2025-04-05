package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"
	
	"time"

	"html/template"
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

	// var jwtKey = []byte("my_secret_key")
	// var tokens []string
	// 
	// type Claims struct {
	// 	Username string `json:"username"`
	// 	jwt.RegisteredClaims
	// }

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

  // front end routes
  r.SetFuncMap(template.FuncMap{
    "dict": func(values ...interface{}) map[string]interface{} {
        m := make(map[string]interface{})
        for i := 0; i < len(values); i += 2 {
            key := values[i].(string)
            m[key] = values[i+1]
        }
        return m
    },
  })
	r.LoadHTMLGlob("src/templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Main website",
			"isIndex":    true,
			"eventsList": []Event{{Title: "Event 1", Blurb: "This is the first event.", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Location: "Location 1" , Img_url: "https://as2.ftcdn.net/v2/jpg/04/96/15/83/1000_F_496158338_SgDd7OQQC2QVfN7U5Qijl2muktM0LjjG.jpg"}, {Title: "Event 2", Blurb: "This is the second event.", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Location: "Location 2", Img_url: "https://images.pexels.com/photos/1105666/pexels-photo-1105666.jpeg"}},
		})
	})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title":   "Main website",
			"img_url": "https://upload.wikimedia.org/wikipedia/commons/thumb/5/52/Spider-Man.jpg/1200px-Spider-Man.jpg",
			"name":    "Spider-man",
			"bio":     "I am a superhero from New York City. I have spider-like abilities and I fight crime.",
			"events":  []Event{{Title: "Event 1", Blurb: "This is the first event.", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Location: "Location 1"}, {Title: "Event 2", Blurb: "This is the second event.", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Location: "Location 2"}},
		})
	})

	r.GET("/events", func(c *gin.Context) {
		c.HTML(http.StatusOK, "events.html", gin.H{
			"title":    "Main website",
			"isEvents": true,
			"Events":   []Event{{Title: "Event 1", Blurb: "This is the first event.", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Location: "Location 1"}, {Title: "Event 2", Blurb: "This is the second event.", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Location: "Location 2"}},
		})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"title": "Signup Page",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login Page",
		})
	})

	r.POST("/api/signup", func(c *gin.Context) {
		log.Println("recv'd ")

		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Printf("errored:%s\n", err)
			return
		}

		// try inserting into db
		err := createUser(DBclient, user)
		if err != nil {
			log.Printf("failed to create user: %s", err)
		}

		// If data binding is successful, return the user information
		c.JSON(http.StatusOK, gin.H{"message": "user Created!", "user": user})
		log.Printf("%+v\n", user)
	})

	r.POST("/api/login", func(c *gin.Context) {
		log.Println("recv'd ")

		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Printf("errored:%s\n", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user Created!", "user": user})
		log.Printf("%+v\n", user)
	})



	//front end routes
	r.Static("/js", "src/static/js")
	r.Static("/css", "src/static/css")
	r.Static("/imgs", "src/static/imgs")

	// run the server
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
