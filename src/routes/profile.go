package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"../data"
)

//GetProfile - Returns user JSON object
func GetProfile(w http.ResponseWriter, r *http.Request) {

	var user data.User

	var users *mongo.Collection = data.DatabaseClient.Database(data.DatabaseName).Collection("users")
	users.FindOne(context.TODO(), bson.D{{Key: "email", Value: r.Header.Get("email")}}).Decode(&user)

	jsonUser, _ := json.Marshal(user)
	w.Write(jsonUser)
}
