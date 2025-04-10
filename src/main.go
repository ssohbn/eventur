package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
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

// TODO: deprecate
// username should be passed thru gin context 
// should be fairly elegant to set thru authorization middleware
// need to write the middleware...
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

type Api struct {
	gemini_key string
	mongo_uri  string
	dbclient   *mongo.Client
}

func (api *Api) createEventRoute(c *gin.Context) {
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

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", api.gemini_key)
	prompt := fmt.Sprintf(`
		Generate a simple list of tags for a new event listing. give only the tags separated by commas. 
		Title: %s
		Blurb: %s
		Description: %s
	`, event.Title, event.Blurb, event.Description)

	jsonData := fmt.Sprintf(`{
		"contents": [
			{"parts": [
				{"text": %s}
			]}
		]}`, prompt)

	log.Printf("key: %s\n", string(api.gemini_key))
	log.Printf("url: %s\n", string(url))
	log.Printf("jsondata: %s\n", string(jsonData))

	resp, err := http.Post(url, "application/json", strings.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Println("Status Code:", resp.StatusCode)
	log.Println("Response Body:", string(body))

	// good paste
	var data struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	tags := data.Candidates[0].Content.Parts[0].Text
	event.Tags = tags

	log.Println(tags)

	err = createEvent(api.dbclient, event)
	if err != nil {
		log.Printf("failed to create event: %+v, err: %s", event, err)
		return
	}
}

func (api *Api) events(c *gin.Context) {
	c.JSON(http.StatusOK, getEvents(api.dbclient))
}

func (api *Api) listUsers(c *gin.Context) {
	// pray no error.
	users, _ := allUsers(api.dbclient)

	c.JSON(http.StatusOK, users)
}

func (api *Api) signup(c *gin.Context) {
	log.Println("recv'd ")

	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("errored:%s\n", err)
		return
	}

	// try inserting into db
	err := createUser(api.dbclient, user)
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
}
func (api *Api) interested(c *gin.Context) {
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
	addInterest(api.dbclient, Interest{Username: username, Event: form.EventName})
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
			log.Fatalf("failed to disconnect from mongodb client cleanly: %s", err)
		}
	}()

	godotenv.Load()
	api := Api{
		gemini_key: os.Getenv("GEMINI_APIKEY"),
		mongo_uri:  os.Getenv("MONGOURI"),
		dbclient:   DBclient,
	}

	r := gin.Default()

	r.POST("/api/signup", api.signup)

	authenticated := r.Group("/", gin.BasicAuth(accounts(DBclient)))
	authenticated.POST("/api/interested", api.interested)
	authenticated.POST("/api/createEvent", api.createEventRoute)
	authenticated.GET("/api/events", api.events)
	authenticated.GET("/api/users", api.listUsers)

	// TODO: think more abt Api struct pattern
	// for now it feels nice enough
	// none of the below routes require db/gemini access
	// unsure what a nice way to handle this is
	// the api struct feels niceish but maybe it should be Env?
	// can I justify putting everything under api?
	// these all happen to also be frontend routes
	authenticated.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Main website",
			"isIndex":    true,
			"eventsList": getEvents(DBclient),
		})
	})

	authenticated.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"title": "Main website",
		})
	})

	authenticated.GET("/profile", func(c *gin.Context) {
		name, err := usernameFromAuthorization(c)
		if err != nil {
			log.Printf("failed to get username in profile: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// redirect to /profile/:name
		c.Redirect(http.StatusFound, "/profile/"+name)
	})

	authenticated.GET("/RSVP/:event", func(c *gin.Context) {
		c.HTML(http.StatusOK, "rsvp.html", gin.H{
			"title": "Main website",
		})
	})

	authenticated.GET("/profile/:name", func(c *gin.Context) {
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

	authenticated.GET("/events", func(c *gin.Context) {
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

	authenticated.GET("/filter", func(c *gin.Context) {
		c.HTML(http.StatusOK, "filter.html", gin.H{
			"title": "Main website",
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

	//front end routes
	r.Static("/js", "src/static/js")
	r.Static("/css", "src/static/css")
	r.Static("/imgs", "src/static/imgs")

	// run the server
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
