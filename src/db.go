package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TYPES

type User struct {
	Username string `form:"username" binding:"required" bson:"username"`
	Password string `form:"password" binding:"required" bson:"password"`
	Token string `bson:"token"`

	// not checking the email binding is sketchy! 
	// recall: we are at a hackathon
	//
	// email binding fails on blank email but 
	// the field shouldnt be required.
	// I dont care to write my own function that binds it.
	Email string `form:"email" bson:"email"` // binding:"email"` 
	Bio      string `form:"bio" bson:"bio"`
}

type Event struct {
	Title string `form:"title" binding:"required" bson:"title"`
	Blurb string `form:"blurb" bson:"blurb"`

	// not required because tbd dates are allowed
	Date time.Time `form:"date" time_format:"2006-01-02" bson:"time"`

	// not required because tbd locations are allowed
	// maybe should be required for like city or something
	Location string `form:"location" bson:"location"`

	Img_url string `form:"img_url" bson:"img_url"`
}

func connectDB() (*mongo.Client, error) {
	godotenv.Load()
	uri := os.Getenv("MONGOURI")
	if uri == "" {
		log.Fatal("Set your 'MONGOURI' environment variable.")
	}
	// log.Println(uri)

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to db:, %s"))
	}

	return client, nil
}

func createUser(client *mongo.Client, user User) error {
	coll := client.Database("eventure").Collection("users")

	// we shold have unique names but  EH.  WHO CARE.  HACKATHON
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("Error inserting event:", err)
		return err
	} else {
		fmt.Println("Inserted event successfully")
	}
	return nil
}

func findUser(client *mongo.Client, username string) (User, error) {
	coll := client.Database("eventure").Collection("events")

	filter := bson.D{{"name", username}}

	result := User{}

	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	// Prints a message if no documents are matched or if any
	// other errors occur during the operation
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, errors.New("couldng get da document")
		}
		return result, errors.New(fmt.Sprintf("findone fucked up: %s\n", err))
	}

	return result, nil 
}

func allUsers(client *mongo.Client) ([]User, error) {
	var users []User
	
	return users, nil
}

// create event
func createEvent(client *mongo.Client, event Event) error {
	coll := client.Database("eventure").Collection("events")

	_, err := coll.InsertOne(context.TODO(), event)
	if err != nil {
		fmt.Println("Error inserting event:", err)
		return err
	} else {
		fmt.Println("Inserted event successfully")
	}
	return nil
}
