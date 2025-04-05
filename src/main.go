package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"

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
	r.GET("/api/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, getEvents(DBclient))
	})

	r.GET("/api/users", func(c *gin.Context) {
		// pray no error.
		users, _ := allUsers(DBclient)

		c.JSON(http.StatusOK, users)
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
			"eventsList": getEvents(DBclient),
		})
	})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/profile/:name", func(c *gin.Context) {
		name := c.Param("name")
		user, err := findUser(DBclient, name)
		if err != nil {
			log.Printf("failed to find user: %s", err)
		}
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title":   "Main website",
			"img_url": "https://upload.wikimedia.org/wikipedia/commons/thumb/5/52/Spider-Man.jpg/1200px-Spider-Man.jpg",
			"name":    user.Username,
			"bio":     user.Bio,
			"events":  getEventsByDirector(DBclient,user.Username),
		})
	})

	r.GET("/events", func(c *gin.Context) {
		c.HTML(http.StatusOK, "events.html", gin.H{
			"title":    "Main website",
			"isEvents": true,
			"Events":   getEvents(DBclient),
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
