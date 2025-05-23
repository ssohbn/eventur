package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Returns mongodb client thingy uhm remember to disconnect after
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
		fmt.Println("Error inserting user:", err)
		return err
	} else {
		fmt.Println("Inserted user successfully")
	}
	return nil
}

func findUser(client *mongo.Client, username string) (User, error) {
	coll := client.Database("eventure").Collection("users")

	filter := bson.D{{"username", username}}

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
	coll := client.Database("eventure").Collection("users")
	filter := bson.D{} // empty for all users

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find all users: %s\n", err))
	}

	// Unpacks the cursor into a slice
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unpack users into slice: %s\n"))
	}

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
func getUser(client *mongo.Client) []User {
	coll := client.Database("eventure").Collection("users")

	filter := bson.M{} // empty filter to get all documents
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []User
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func addInterest(client *mongo.Client, interest Interest) error {
	coll := client.Database("eventure").Collection("interest")
	_, err := coll.InsertOne(context.TODO(), interest)
	if err != nil {
		fmt.Println("Error inserting interest:", err)
		return err
	} else {
		fmt.Println("Inserted interest successfully")
	}
	return nil
}

// get events interested by user
func getInterestedEvents(client *mongo.Client, username string) []Event {
	coll := client.Database("eventure").Collection("interest")

	filter := bson.M{"username": username}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var eventNames []Interest
	if err = cursor.All(context.TODO(), &eventNames); err != nil {
		panic(err)
	}

	// get events from event names
	var events []Event
	for _, eventName := range eventNames {
		eventColl := client.Database("eventure").Collection("events")
		eventFilter := bson.M{"title": eventName.Event}
		eventCursor, err := eventColl.Find(context.TODO(), eventFilter)
		if err != nil {
			panic(err)
		}

		var result []Event
		if err = eventCursor.All(context.TODO(), &result); err != nil {
			panic(err)
		}
		events = append(events, result...)
	}

	return events
}
