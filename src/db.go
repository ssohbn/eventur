package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TYPES

type User struct {
	Username string `form:"username" binding:"required"`
	Email string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`

	Bio      string `form:"bio"`
}

type Event struct {
	Title    string `form:"title" binding:"required"`
	Director string
	Blurb    string `form:"blurb"`

	// not required because tbd dates are allowed
	Date time.Time `form:"date" time_format:"2006-01-02"`

	// not required because tbd locations are allowed
	// maybe should be required for like city or something
	Location string `form:"location"`

	Img_url string `form:"img_url"`
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
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("Error inserting event:", err)
		return err
	} else {
		fmt.Println("Inserted event successfully")
	}
	return nil
}

// create event
func createEvent(client *mongo.Client, event Event) error {
	coll := client.Database("eventure").Collection("events")
	_, err := coll.InsertOne(context.Background(), event)
	if err != nil {
		fmt.Println("Error inserting event:", err)
		return err
	} else {
		fmt.Println("Inserted event successfully")
	}
	return nil
}

// get events
func getEvents(client *mongo.Client) []Event {
	coll := client.Database("eventure").Collection("events")

	filter := bson.M{} // empty filter to get all documents
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Event
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func getEventsByDirector(client *mongo.Client, directorName string) []Event {
	coll := client.Database("eventure").Collection("events")

	filter := bson.M{"director": directorName} // filter to get documents by director
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Event
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

// get user
func getUser(client *mongo.Client) []Event {
	coll := client.Database("eventure").Collection("users")

	filter := bson.M{} // empty filter to get all documents
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Event
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}
