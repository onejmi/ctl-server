package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var users *mongo.Collection

//Register - Register route (for creating a user)
func Register(w http.ResponseWriter, r *http.Request) {
	users = data.DatabaseClient.Database(data.DatabaseName).Collection("users")
	raw, _ := ioutil.ReadAll(r.Body)
	var newUser data.User
	json.Unmarshal(raw, &newUser)
	if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" {
		jsonError, _ := json.Marshal(data.Error{Message: "Please specify all user fields."})
		w.Write(jsonError)
	} else if exists(users, newUser.Username) {
		jsonError, _ := json.Marshal(data.Error{Message: "That user already exists"})
		w.Write(jsonError)
	} else {
		addUser(newUser)
		userJSON, _ := json.Marshal(newUser)
		w.Write(userJSON)
	}
}

func addUser(user data.User) {
	insertResponse, err := users.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted user: ", insertResponse.InsertedID)
}

func exists(users *mongo.Collection, username string) bool {
	count, err := users.CountDocuments(context.TODO(), bson.D{{Key: "username", Value: username}})
	if err != nil {
		log.Fatal(err)
	}
	return (count > 0)
}
