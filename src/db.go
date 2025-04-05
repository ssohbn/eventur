package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TYPES

type Event struct {
	Title string `form:"title" binding:"required"`
	Blurb string `form:"blurb"`

	// not required because tbd dates are allowed
	Date time.Time `form:"date" time_format:"2006-01-02"`

	// not required because tbd locations are allowed
	// maybe should be required for like city or something
	Location string `form:"location"`
}

func connectDB() *mongo.Client {
	godotenv.Load()
	uri := os.Getenv("MONGOURI")
	if uri == "" {
		log.Fatal("Set your 'MONGOURI' environment variable.")
	}
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		return nil
	}

	return client
}

// create event
func createEvent(client *mongo.Client, event Event) error {
	coll := client.Database("eventure").Collection("events")
	_, err := coll.InsertOne(context.TODO(), Event{
		Title:    "Test Event",
		Blurb:    "This is a test event.",
		Date:     time.Now(),
		Location: "Test Location",
	})
	if err != nil {
		fmt.Println("Error inserting event:", err)
		return err
	} else {
		fmt.Println("Inserted event successfully")
	}
	return nil
}
