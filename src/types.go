package main

import (
	"time"
)

// TYPES
type User struct {
	Username string `form:"username" binding:"required" bson:"username"`
	Password string `form:"password" binding:"required" bson:"password"`
	Token    string `bson:"token"`

	// not checking the email binding is sketchy!
	// recall: we are at a hackathon
	//
	// email binding fails on blank email but
	// the field shouldnt be required.
	// I dont care to write my own function that binds it.
	Email   string `form:"email" bson:"email"` // binding:"email"`
	Bio     string `form:"bio" bson:"bio"`
	Img_url string `form:"image" bson:"img_url"`
}

type Event struct {
	Title       string `form:"title" binding:"required" bson:"title"`
	Blurb       string `form:"blurb" bson:"blurb"`
	Description string `form:"description" bson:"description" binding:"required"`
	Tags        string `bson:"tags"`
	Director    string

	// not required because tbd dates are allowed
	Date time.Time `form:"date" time_format:"2006-01-02" bson:"time"`

	// not required because tbd locations are allowed
	// maybe should be required for like city or something
	Location string `form:"location" bson:"location"`

	Img_url string `form:"image" bson:"img_url"`
}

type Interest struct {
	Username string `form:"username" binding:"required" bson:"username"`
	Event    string `form:"event" binding:"required" bson:"event"`
}
