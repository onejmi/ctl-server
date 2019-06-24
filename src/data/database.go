package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DatabaseName - Name of main database which should be accessed by the application
const DatabaseName = "core"

//DatabaseClient - mongodb client which can be used to interface with the database
var DatabaseClient *mongo.Client

//Connect - connects to database and logs any errors that may have occured while doing so
func Connect() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	potentialClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = potentialClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	DatabaseClient = potentialClient

	fmt.Println("Connected to MongoDB!")
}
