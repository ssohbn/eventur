package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"html/template"
)

func accounts(db *mongo.Client) gin.Accounts {
	accounts := make(gin.Accounts)
	users, err := allUsers(db)
	if err != nil {
		log.Panicf("failed to access acounts! this is very bad. %s\n", err)
	}

	for _, user := range users {
		accounts[user.Username] = user.Password
	}

	return accounts
}

func usernameFromAuthorization(c *gin.Context) (string, error) {
	header := c.Request.Header["Authorization"][0]
	combo := strings.Split(header, " ")[1]

	authbytes, err := base64.StdEncoding.DecodeString(combo)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to decode authorization header: %s\n", err))
	}

	log.Printf("authbytes: %s, authorizationheader: %s", authbytes, header)
	message := string(authbytes)
	log.Printf("message decoded: %s", message)
	username := strings.Split(message, ":")[0]

	return username, nil
}

func hasAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // skip this god awful broken middleware
		header, ok := c.Request.Header["Authorization"]
		log.Printf("header gunk: %s\n", header)
		if !ok {
			c.Redirect(http.StatusFound, "/signup")
		}
		// c.Next()
	}
}

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
	r.POST("/api/createEvent", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {

		username, err := usernameFromAuthorization(c)
		if err != nil {
			log.Printf("failed to get username in createEvent: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(fmt.Sprintf("recv'd createEvent from %s from %s", username, username))
		var event Event

		// bind form data (WHICH SHOULD BIND) to event and check if error is produced (not nil)
		// (sorta weird syntax)
		if err := c.ShouldBind(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		event.Director = username

		// try inserting into db
		err = createEvent(DBclient, event)
		if err != nil {
			log.Printf("failed to create event %s", err)
		}

		// If data binding is successful, return the user information
		c.JSON(http.StatusOK, gin.H{"message": "Event Created!", "event": event})
		log.Printf("%+v\n", event)
	})

	r.GET("/api/events", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		c.JSON(http.StatusOK, getEvents(DBclient))
	})

	r.GET("/api/users", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
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

	r.GET("/", hasAuth(), gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Main website",
			"isIndex":    true,
			"eventsList": getEvents(DBclient),
		})
	})

	r.GET("/create", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/profile", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		name, err := usernameFromAuthorization(c)
		if err != nil {
			log.Printf("failed to get username in profile: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// redirect to /profile/:name
		c.Redirect(http.StatusFound, "/profile/"+name)
	})

  r.GET("/RSVP/:event", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		c.HTML(http.StatusOK, "rsvp.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/profile/:name", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		name := c.Param("name")
		user, err := findUser(DBclient, name)
		if err != nil {
			log.Printf("failed to find user: %s", err)
		}
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title":   "Main website",
			"img_url": user.Img_url,
			"name":    user.Username,
			"bio":     user.Bio,
			"events":  getEventsByDirector(DBclient, user.Username),
		})
	})

	r.GET("/events", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		username, err := usernameFromAuthorization(c)
		if err != nil {
			log.Printf("failed to get username in getEvent: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "events.html", gin.H{
			"title":    "Main website",
			"isEvents": true,
			"Events":   getInterestedEvents(DBclient, username),
		})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"title": "Signup Page",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
		// c.HTML(http.StatusOK, "login.html", gin.H{
		// 	"title": "Login Page",
		// })
	})

	r.GET("/filter", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		c.HTML(http.StatusOK, "filter.html", gin.H{
			"title": "Main website",
		})
	})

	r.POST("/api/interested", gin.BasicAuth(accounts(DBclient)), func(c *gin.Context) {
		type FormData struct {
			EventName string `form:"eventName"`
		}
		var form FormData
		username, err := usernameFromAuthorization(c)
		if err != nil {
			log.Printf("failed to get username in interested: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = c.ShouldBindJSON(&form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Printf("errored:%s\n", err)
			return
		}
		log.Println("event", form.EventName)
		addInterest(DBclient, Interest{Username: username, Event: form.EventName})
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

		data := fmt.Sprintf("%s:%s", user.Username, user.Password)
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		header := fmt.Sprintf("Basic %s", encoded)

		log.Printf("data %s, bytes %v\n", data, header)
		log.Printf("%+v\n", user)

		// If data binding is successful, return the user information
		// WE DIRELY NEED TO ACCEPT THESE HEADERS IN JAVASCRIPT
		// THE ENTIRE PROGRAM IS SOFTLOCKED UNTIL THIS IS ACCEPTED.
		// DO THIS IN SIGNUP ON THE RESPONSE FROM FETCH REQUEST TO THIS API ENDPOINT
		// FIX
		// FIX
		// FIX
		// FIX
		// FIX
		// FIX
		// FIX
		// FIX
		c.Header("WWW-Authenticate", `Basic realm="dear god help"`)
		// c.Request.SetBasicAuth(user.Username, user.Password)
		c.JSON(http.StatusOK, gin.H{"message": "user Created!", "user": user, "Authorization": header})
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
